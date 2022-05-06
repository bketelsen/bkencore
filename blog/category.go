package blog

import (
	"context"
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type GetCategoriesResponse struct {
	Count      int         `json:"count,omitempty"`
	Categories []*Category `json:"categories,omitempty"`
}
type Category struct {
	Category string `json:"category,omitempty"`
	Summary  string `json:"summary,omitempty"`
}

// GetCategory retrieves a category by slug.
//encore:api public method=GET path=/category/:category
func GetCategory(ctx context.Context, category string) (*Category, error) {
	var (
		c Category
	)
	err := sqldb.QueryRow(ctx, `
		SELECT category , summary
		FROM "category"
		WHERE category = $1
	`, category).Scan(&c.Category, &c.Summary)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "tag not found",
		}
	}
	return &c, nil
}

// CreateTag creates a new blog post.
//encore:api auth
func CreateCategory(ctx context.Context, category *Category) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "category" (category, summary)
		VALUES ($1,  $2)
		ON CONFLICT (category) DO UPDATE
		SET summary = $2
	`, category.Category, category.Summary)

	if err != nil {
		return fmt.Errorf("insert category: %v", err)
	}

	return nil

}

// GetCategories retrieves a list of categories
//encore:api public method=GET path=/category
func GetCategories(ctx context.Context) (*GetCategoriesResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT category, summary
		FROM "category"
	`)
	if err != nil {
		return &GetCategoriesResponse{
			Count:      0,
			Categories: []*Category{},
		}, err
	}
	defer rows.Close()

	var q []*Category
	var i = 0
	for rows.Next() {
		var (
			c Category
		)
		err := rows.Scan(&c.Category, &c.Summary)
		if err != nil {
			return &GetCategoriesResponse{
				Count:      0,
				Categories: []*Category{},
			}, err
		}
		q = append(q, &c)
		i = i + 1
	}
	return &GetCategoriesResponse{
		Count:      i,
		Categories: q,
	}, rows.Err()
}
