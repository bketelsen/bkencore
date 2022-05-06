package blog

import (
	"context"
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type Tag struct {
	Tag     string `json:"tag,omitempty"`
	Summary string `json:"summary,omitempty"`
}
type GetTagsResponse struct {
	Count int    `json:"count,omitempty"`
	Tags  []*Tag `json:"tags,omitempty"`
}

// GetTag retrieves a tag by slug.
//encore:api public method=GET path=/tag/:tag
func GetTag(ctx context.Context, tag string) (*Tag, error) {
	var (
		t Tag
	)
	err := sqldb.QueryRow(ctx, `
		SELECT tag , summary
		FROM "tag"
		WHERE tag = $1
	`, tag).Scan(&t.Tag, &t.Summary)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "tag not found",
		}
	}
	return &t, nil
}

// CreateTag creates a new blog post.
//encore:api auth
func CreateTag(ctx context.Context, tag *Tag) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "tag" (tag, summary)
		VALUES ($1,  $2)
		ON CONFLICT (tag) DO UPDATE
		SET summary = $2
	`, tag.Tag, tag.Summary)

	if err != nil {
		return fmt.Errorf("insert tag: %v", err)
	}

	return nil

}

// GetTags retrieves a list of tags
//encore:api public method=GET path=/tag
func GetTags(ctx context.Context) (*GetTagsResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT tag, summary
		FROM "tag"
	`)
	if err != nil {
		return &GetTagsResponse{
			Count: 0,
			Tags:  []*Tag{},
		}, err
	}
	defer rows.Close()

	var q []*Tag
	var i = 0
	for rows.Next() {
		var (
			t Tag
		)
		err := rows.Scan(&t.Tag, &t.Summary)
		if err != nil {
			return &GetTagsResponse{
				Count: 0,
				Tags:  []*Tag{},
			}, err
		}
		q = append(q, &t)
		i = i + 1
	}
	return &GetTagsResponse{
		Count: i,
		Tags:  q,
	}, rows.Err()
}

// GetTagsBySlug retrieves a list of tags for a post
//encore:api public method=GET path=/tagbyslug/:slug
func GetTagsBySlug(ctx context.Context, slug string) (*GetTagsResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT tag, summary
		FROM "tag"
		WHERE tag IN (select tag from article_tag where slug = $1)
	`, slug)
	if err != nil {
		return &GetTagsResponse{
			Count: 0,
			Tags:  []*Tag{},
		}, err
	}
	defer rows.Close()

	var q []*Tag
	var i = 0
	for rows.Next() {
		var (
			t Tag
		)
		err := rows.Scan(&t.Tag, &t.Summary)
		if err != nil {
			return &GetTagsResponse{
				Count: 0,
				Tags:  []*Tag{},
			}, err
		}
		q = append(q, &t)
		i = i + 1
	}
	return &GetTagsResponse{
		Count: i,
		Tags:  q,
	}, rows.Err()
}
