package blog

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/russross/blackfriday/v2"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type Page struct {
	Slug          string
	CreatedAt     time.Time `json:"created_at" yaml:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
	Published     bool
	Title         string
	Subtitle      string
	HeroText      string
	Summary       string
	Body          string
	BodyRendered  string
	FeaturedImage string `json:"featured_image"`
}

type CreatePageParams struct {
	Published     bool
	Title         string
	Subtitle      string
	HeroText      string
	Summary       string
	Body          string
	FeaturedImage string // empty string means no image
}

// GetPage retrieves a page by slug.
//encore:api public method=GET path=/page/:slug
func GetPage(ctx context.Context, slug string) (*Page, error) {
	var (
		p   Page
		img sql.NullString
	)
	err := sqldb.QueryRow(ctx, `
		SELECT slug, created_at, modified_at, published, title, subtitle, hero_text, summary, body, body_rendered, featured_image
		FROM "page"
		WHERE slug = $1
	`, slug).Scan(&p.Slug, &p.CreatedAt, &p.ModifiedAt, &p.Published, &p.Title, &p.Subtitle, &p.HeroText, &p.Summary, &p.Body, &p.BodyRendered, &img)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "page not found",
		}
	}
	p.FeaturedImage = img.String
	return &p, nil
}

// CreatePage creates a new page, or updates it if it already exists.
//encore:api auth method=PUT path=/page/:slug
func CreatePage(ctx context.Context, slug string, p *CreatePageParams) error {
	img := sql.NullString{String: p.FeaturedImage, Valid: p.FeaturedImage != ""}
	rendered := blackfriday.Run([]byte(p.Body))
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "page" (
			slug, created_at, modified_at, published, title, subtitle,
			hero_text, summary, body, body_rendered, featured_image
		)
		VALUES ($1, NOW(), NOW(), $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (slug) DO UPDATE
		SET modified_at = NOW(), published = $2, title = $3, subtitle = $4,
			hero_text = $5, summary = $6, body = $7, body_rendered = $8, featured_image = $9
	`, slug, p.Published, p.Title, p.Subtitle, p.HeroText, p.Summary, p.Body, string(rendered), img)
	if err != nil {
		return fmt.Errorf("insert page: %v", err)
	}
	return nil
}
