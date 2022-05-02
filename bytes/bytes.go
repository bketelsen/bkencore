// Service bytes stores interesting links.
package bytes

import (
	"context"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type PublishParams struct {
	Title   string
	Summary string
	URL     string
}

type PublishResponse struct {
	ID int64
}

// Publish publishes a byte.
//encore:api auth method=POST path=/bytes
func Publish(ctx context.Context, p *PublishParams) (*PublishResponse, error) {
	var id int64
	err := sqldb.QueryRow(ctx, `
		INSERT INTO byte (title, summary, url)
		VALUES ($1, $2, $3)
		ON CONFLICT (url) DO UPDATE
		SET id = byte.id
		RETURNING id
	`, p.Title, p.Summary, p.URL).Scan(&id)
	return &PublishResponse{ID: id}, err
}

type ListParams struct {
	Limit  int
	Offset int
}

type Byte struct {
	ID      int64
	Title   string
	Summary string
	URL     string
	Created time.Time
}

type ListResponse struct {
	Bytes []Byte
}

// Get retrieves a byte.
//encore:api public method=GET path=/bytes/:id
func Get(ctx context.Context, id int64) (*Byte, error) {
	var b Byte
	err := sqldb.QueryRow(ctx, `
		SELECT id, title, summary, url, created_at	
		FROM "byte"
		WHERE id = $1
	`, id).Scan(&b.ID, &b.Title, &b.Summary, &b.URL, &b.Created)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "byte not found",
		}
	}
	return &b, nil

}

// List lists published bytes.
//encore:api public method=GET path=/bytes
func List(ctx context.Context, p *ListParams) (*ListResponse, error) {
	offset := getOrDefault(p.Offset, 0)
	limit := getOrDefault(p.Limit, 100)
	rows, err := sqldb.Query(ctx, `
		SELECT id, title, summary, url, created_at
		FROM byte 
		ORDER BY id desc
		OFFSET $1
		LIMIT $2
	`, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bytes []Byte
	for rows.Next() {
		var b Byte
		if err := rows.Scan(&b.ID, &b.Title, &b.Summary, &b.URL, &b.Created); err != nil {
			return nil, err
		}
		bytes = append(bytes, b)
	}
	return &ListResponse{Bytes: bytes}, rows.Err()
}

func getOrDefault(n, def int) int {
	if n == 0 {
		return def
	}
	return n
}
