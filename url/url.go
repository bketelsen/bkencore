package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/url"

	"encore.dev/storage/sqldb"
)

type URL struct {
	ID       string `json:"id,omitempty"`        // short-form URL id
	URL      string `json:"url,omitempty"`       // original URL, in long form
	ShortURL string `json:"short_url,omitempty"` // short URL
}
type GetListResponse struct {
	Count int    `json:"count,omitempty"`
	URLS  []*URL `json:"urls,omitempty"`
}

type ShortenParams struct {
	URL string `json:"url,omitempty"` // the URL to shorten
}

// Shorten shortens a URL.
//encore:api public method=POST path=/url
func Shorten(ctx context.Context, p *ShortenParams) (*URL, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	} else if err := insert(ctx, id, p.URL); err != nil {
		return nil, err
	}
	return &URL{
		ID:       id,
		URL:      p.URL,
		ShortURL: FormatShortURL(id),
	}, nil
}

// Get retrieves the original URL for the id.
//encore:api public method=GET path=/url/:id
func Get(ctx context.Context, id string) (*URL, error) {
	u := &URL{ID: id}
	err := sqldb.QueryRow(ctx, `
        SELECT original_url FROM url
        WHERE id = $1
    `, id).Scan(&u.URL)
	return u, err
}

// List retrieves all shortened URLs
//encore:api public method=GET path=/url
func List(ctx context.Context) (*GetListResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT id, original_url
		FROM "url"
	`)
	if err != nil {
		return &GetListResponse{
			Count: 0,
			URLS:  []*URL{},
		}, err
	}
	defer rows.Close()

	var q []*URL
	var i = 0
	for rows.Next() {
		var (
			b URL
		)
		err := rows.Scan(&b.ID, &b.URL)
		if err != nil {
			return &GetListResponse{
				Count: 0,
				URLS:  []*URL{},
			}, err
		}

		q = append(q, &b)
		i = i + 1
	}
	return &GetListResponse{
		Count: i,
		URLS:  q,
	}, rows.Err()
}

// generateID generates a random short ID.
func generateID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// insert inserts a URL into the database.
func insert(ctx context.Context, id, url string) error {
	_, err := sqldb.Exec(ctx, `
        INSERT INTO url (id, original_url)
        VALUES ($1, $2)
    `, id, url)
	return err
}

// FormatShortURL formats a short id into the short URL.
func FormatShortURL(id string) string {
	return "https://url.bjk.fyi/" + url.PathEscape(id)
}
