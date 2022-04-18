package blog

import (
	"context"
	"fmt"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"github.com/russross/blackfriday/v2"
)

type BlogPost struct {
	Slug         string
	CreatedAt    time.Time
	Published    bool
	ModifiedAt   time.Time
	Title        string
	Summary      string
	Body         string
	BodyRendered string
}

type CreateBlogPostParams struct {
	Slug      string
	Published bool
	Title     string
	Summary   string
	Body      string
}

type GetBlogPostParams struct {
	Slug string
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
	// Use slug to query database...
	var b BlogPost
	err := sqldb.QueryRow(ctx, `
		SELECT slug, created_at, published, modified_at, title, summary, body, body_rendered
		FROM "article"
		WHERE slug = $1
	`, slug).Scan(&b.Slug, &b.CreatedAt, &b.Published, &b.ModifiedAt, &b.Title, &b.Summary, &b.Body, &b.BodyRendered)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "article not found",
		}
	}
	return &b, nil
}

// CreateBlogPost creates a new blog post.
//encore:api auth
func CreateBlogPost(ctx context.Context, params *CreateBlogPostParams) error {

	rendered := blackfriday.Run([]byte(params.Body))
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "article" (slug, created_at, published, modified_at, title, summary,body, body_rendered)
		VALUES ($1,  NOW(), $2, NOW(), $3, $4, $5, $6)
		ON CONFLICT (slug) DO UPDATE
		SET published = $2, modified_at = NOW(), title = $3, summary = $4, body = $5, body_rendered = $6

	`, params.Slug, params.Published, params.Title, params.Summary, params.Body, string(rendered))
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
	SELECT slug, created_at, published, modified_at, title, summary, body, body_rendered
	FROM "article"
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
		var b BlogPost
		err := rows.Scan(&b.Slug, &b.CreatedAt, &b.Published, &b.ModifiedAt, &b.Title, &b.Summary, &b.Body, &b.BodyRendered)
		if err != nil {
			return &GetBlogPostsResponse{
				Count:     0,
				BlogPosts: []*BlogPost{},
			}, err
		}
		q = append(q, &b)
		i = i + 1
	}
	return &GetBlogPostsResponse{
		Count:     i,
		BlogPosts: q,
	}, rows.Err()
}
