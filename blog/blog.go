package blog

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/russross/blackfriday/v2"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type BlogPost struct {
	Slug          string    `json:"slug,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty" yaml:"created_at"`
	ModifiedAt    time.Time `json:"modified_at,omitempty"`
	Published     bool      `json:"published,omitempty"`
	Title         string    `json:"title,omitempty"`
	Summary       string    `json:"summary,omitempty"`
	Body          string    `json:"body,omitempty"`
	BodyRendered  string    `json:"body_rendered,omitempty"`
	FeaturedImage string    `json:"featured_image,omitempty"`
	Category      *Category `json:"category"`
	Tags          []*Tag    `json:"tags"`
}

type CreateBlogPostParams struct {
	Slug          string   `json:"slug,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	ModifiedAt    string   `json:"modified_at,omitempty"`
	Published     bool     `json:"published,omitempty"`
	Title         string   `json:"title,omitempty"`
	Summary       string   `json:"summary,omitempty"`
	Body          string   `json:"body,omitempty"`
	FeaturedImage string   `json:"featured_image,omitempty"`
	Category      string   `json:"category,omitempty"`
	Tags          []string `json:"tags"`
}

type GetBlogPostsParams struct {
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
}

type GetBlogPostsResponse struct {
	Count     int         `json:"count,omitempty"`
	BlogPosts []*BlogPost `json:"blog_posts"`
}

// GetBlogPost retrieves a blog post by slug.
//encore:api public method=GET path=/blog/:slug
func GetBlogPost(ctx context.Context, slug string) (*BlogPost, error) {
	var (
		b        BlogPost
		img      sql.NullString
		category sql.NullString
	)
	err := sqldb.QueryRow(ctx, `
		SELECT slug, created_at, published, modified_at, title, summary, body, body_rendered, featured_image, category
		FROM "article"
		WHERE slug = $1
	`, slug).Scan(&b.Slug, &b.CreatedAt, &b.Published, &b.ModifiedAt, &b.Title, &b.Summary, &b.Body, &b.BodyRendered, &img, &category)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "article not found",
		}
	}
	b.FeaturedImage = img.String
	tbsr, err := GetTagsBySlug(ctx, b.Slug)
	b.Tags = tbsr.Tags
	if category.Valid {
		cat, err := GetCategory(ctx, category.String)
		if err != nil {
			return nil, &errs.Error{
				Code:    errs.NotFound,
				Message: "error finding category",
			}
		}
		b.Category = cat
	}
	return &b, nil
}

// CreateBlogPost creates a new blog post.
//encore:api auth
func CreateBlogPost(ctx context.Context, params *CreateBlogPostParams) error {
	img := sql.NullString{String: params.FeaturedImage, Valid: params.FeaturedImage != ""}
	if params.Category == "" {
		params.Category = "miscellaneous"
	}
	rendered := blackfriday.Run([]byte(params.Body))
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "article" (slug, created_at, published, modified_at, title, summary,body, body_rendered, featured_image, category)
		VALUES ($1,  $2, $3,  $4, $5,  $6, $7, $8, $9, $10)
		ON CONFLICT (slug) DO UPDATE
		SET created_at = $2, published = $3, modified_at = $4, title = $5, summary = $6, body = $7, body_rendered = $8, featured_image = $9, category = $10
	`, params.Slug, params.CreatedAt, params.Published, params.ModifiedAt, params.Title, params.Summary, params.Body, string(rendered), img, params.Category)

	if err != nil {
		return fmt.Errorf("insert article: %v", err)
	}
	// now insert tags

	for _, t := range params.Tags {
		_, err = sqldb.Exec(ctx, `
		INSERT INTO "article_tag" (slug, tag)
		VALUES ($1,  $2)
		ON CONFLICT DO NOTHING
	`, params.Slug, t)
		if err != nil {
			return fmt.Errorf("insert article_tag: %v", err)
		}
	}
	return nil

}

// GetBlogPosts retrieves a list of blog posts with
// optional limit and offset.
//encore:api public method=GET path=/blog
func GetBlogPosts(ctx context.Context, params *GetBlogPostsParams) (*GetBlogPostsResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT slug, created_at, published, modified_at, title, summary, body, body_rendered, featured_image, category
		FROM "article"
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return &GetBlogPostsResponse{
			Count:     0,
			BlogPosts: []*BlogPost{},
		}, err
	}
	defer rows.Close()

	var q []*BlogPost
	var i = 0
	for rows.Next() {
		var (
			b        BlogPost
			img      sql.NullString
			category sql.NullString
		)
		err := rows.Scan(&b.Slug, &b.CreatedAt, &b.Published, &b.ModifiedAt, &b.Title, &b.Summary, &b.Body, &b.BodyRendered, &img, &category)
		if err != nil {
			return &GetBlogPostsResponse{
				Count:     0,
				BlogPosts: []*BlogPost{},
			}, err
		}
		b.FeaturedImage = img.String
		tbsr, err := GetTagsBySlug(ctx, b.Slug)
		b.Tags = tbsr.Tags
		if category.Valid {
			cat, err := GetCategory(ctx, category.String)
			if err != nil {
				return &GetBlogPostsResponse{
					Count:     0,
					BlogPosts: []*BlogPost{},
				}, err
			}
			b.Category = cat
		}
		q = append(q, &b)
		i = i + 1
	}
	return &GetBlogPostsResponse{
		Count:     i,
		BlogPosts: q,
	}, rows.Err()
}

//encore:authhandler
func AuthHandler(ctx context.Context, token string) (auth.UID, error) {
	eb := errs.B().Meta("auth", token)

	if token != secrets.AuthPassword {
		return "", eb.Code(errs.Unauthenticated).Msg("authentication failure").Err()
	}
	return "", nil
}
