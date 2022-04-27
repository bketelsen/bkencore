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
	Slug          string
	CreatedAt     time.Time
	Published     bool
	ModifiedAt    time.Time
	Title         string
	Summary       string
	Body          string
	BodyRendered  string
	FeaturedImage string
}

type CreateBlogPostParams struct {
	Slug          string
	Published     bool
	Title         string
	Summary       string
	Body          string
	FeaturedImage string
}

type GetBlogPostsParams struct {
	Limit  int
	Offset int
}

type GetBlogPostsResponse struct {
	Count     int
	BlogPosts []*BlogPost
}

// GetBlogPost retrieves a blog post by slug.
//encore:api public method=GET path=/blog/:slug
func GetBlogPost(ctx context.Context, slug string) (*BlogPost, error) {
	var (
		b   BlogPost
		img sql.NullString
	)
	err := sqldb.QueryRow(ctx, `
		SELECT slug, created_at, published, modified_at, title, summary, body, body_rendered, featured_image
		FROM "article"
		WHERE slug = $1
	`, slug).Scan(&b.Slug, &b.CreatedAt, &b.Published, &b.ModifiedAt, &b.Title, &b.Summary, &b.Body, &b.BodyRendered, &img)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "article not found",
		}
	}
	b.FeaturedImage = img.String
	return &b, nil
}

// CreateBlogPost creates a new blog post.
//encore:api auth
func CreateBlogPost(ctx context.Context, params *CreateBlogPostParams) error {
	img := sql.NullString{String: params.FeaturedImage, Valid: params.FeaturedImage != ""}
	rendered := blackfriday.Run([]byte(params.Body))
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "article" (slug, created_at, published, modified_at, title, summary,body, body_rendered, featured_image)
		VALUES ($1,  NOW(), $2, NOW(), $3, $4, $5, $6, $7)
		ON CONFLICT (slug) DO UPDATE
		SET published = $2, modified_at = NOW(), title = $3, summary = $4, body = $5, body_rendered = $6, featured_image = $7
	`, params.Slug, params.Published, params.Title, params.Summary, params.Body, string(rendered), img)

	if err != nil {
		return fmt.Errorf("insert article: %v", err)
	}

	return nil

}

// GetBlogPosts retrieves a list of blog posts with
// optional limit and offset.
//encore:api public method=GET path=/blog
func GetBlogPosts(ctx context.Context, params *GetBlogPostsParams) (*GetBlogPostsResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT slug, created_at, published, modified_at, title, summary, body, body_rendered, featured_image
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
			b   BlogPost
			img sql.NullString
		)
		err := rows.Scan(&b.Slug, &b.CreatedAt, &b.Published, &b.ModifiedAt, &b.Title, &b.Summary, &b.Body, &b.BodyRendered, &img)
		if err != nil {
			return &GetBlogPostsResponse{
				Count:     0,
				BlogPosts: []*BlogPost{},
			}, err
		}
		b.FeaturedImage = img.String
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
