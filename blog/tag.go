package blog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

type Tag struct {
	Name               string    `json:"slug_name"`
	Slug               string    `json:"slug"`
	Description        string    `json:"slug_description"`
	FeatureImage       string    `json:"feature_image"`
	Visibility         string    `json:"visibility"`
	OgImage            string    `json:"og_image"`
	OgTitle            string    `json:"og_title"`
	OgDescription      string    `json:"og_description"`
	TwitterImage       string    `json:"twitter_image"`
	TwitterTitle       string    `json:"twitter_title"`
	TwitterDescription string    `json:"twitter_description"`
	MetaTitle          string    `json:"meta_title"`
	MetaDescription    string    `json:"meta_description"`
	AccentColor        string    `json:"accent_color"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	URL                string    `json:"slug_url"`
}

type GetTagsResponse struct {
	Count int    `json:"count,omitempty"`
	Tags  []*Tag `json:"tags,omitempty"`
}

// GetTag retrieves a tag by slug.
//encore:api public method=GET path=/tag/:slug
func GetTag(ctx context.Context, slug string) (*Tag, error) {
	var (
		t                   Tag
		slug_description    sql.NullString
		feature_image       sql.NullString
		visibility          sql.NullString
		og_image            sql.NullString
		og_title            sql.NullString
		og_description      sql.NullString
		twitter_image       sql.NullString
		twitter_title       sql.NullString
		twitter_description sql.NullString
		meta_title          sql.NullString
		meta_description    sql.NullString
		accent_color        sql.NullString
	)
	err := sqldb.QueryRow(ctx, `
		SELECT
		slug,
		slug_name,
		slug_description,
		feature_image,
		visibility,
		og_image,
		og_title,
		og_description,
		twitter_image,
		twitter_title,
		twitter_description,
		meta_title,
		meta_description,
		accent_color,
		created_at,
		updated_at,
		slug_url
		FROM "tag"
		WHERE slug = $1
	`, slug).Scan(
		&t.Slug,
		&t.Name,
		&slug_description,
		&feature_image,
		&visibility,
		&og_image,
		&og_title,
		&og_description,
		&twitter_image,
		&twitter_title,
		&twitter_description,
		&meta_title,
		&meta_description,
		&accent_color,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.URL,
	)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "tag not found",
		}
	}
	t.Description = slug_description.String
	t.FeatureImage = feature_image.String
	t.Visibility = visibility.String
	t.OgImage = og_image.String
	t.OgTitle = og_title.String
	t.OgDescription = og_description.String
	t.TwitterImage = twitter_image.String
	t.TwitterTitle = twitter_title.String
	t.TwitterDescription = twitter_description.String
	t.MetaTitle = meta_title.String
	t.MetaDescription = meta_description.String
	t.AccentColor = accent_color.String
	return &t, nil
}

// CreateTag creates a new blog post.
//encore:api private
func CreateTag(ctx context.Context, t *Tag) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "tag" (
			slug,
			slug_name,
			slug_description,
			feature_image,
			visibility,
			og_image,
			og_title,
			og_description,
			twitter_image,
			twitter_title,
			twitter_description,
			meta_title,
			meta_description,
			accent_color,
			created_at,
			updated_at,
			slug_url
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17
		)
		ON CONFLICT (slug) DO UPDATE
		SET
			slug_name = $2,
			slug_description = $3,
			feature_image = $4,
			visibility = $5,
			og_image = $6,
			og_title = $7,
			og_description = $8,
			twitter_image = $9,
			twitter_title = $10,
			twitter_description = $11,
			meta_title = $12,
			meta_description = $13,
			accent_color = $14,
			created_at = $15,
			updated_at = $16,
			slug_url = $17

	`,
		t.Slug,
		t.Name,
		t.Description,
		t.FeatureImage,
		t.Visibility,
		t.OgImage,
		t.OgTitle,
		t.OgDescription,
		t.TwitterImage,
		t.TwitterTitle,
		t.TwitterDescription,
		t.MetaTitle,
		t.MetaDescription,
		t.AccentColor,
		t.CreatedAt,
		t.UpdatedAt,
		t.URL,
	)

	if err != nil {
		return fmt.Errorf("insert tag: %v", err)
	}

	return nil
}

// CreatePostTag creates a new association between a post and a tag.
func CreatePostTag(ctx context.Context, t *Tag, b *BlogPost) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "article_tag" (
			slug,
			tag
		)
		VALUES (
			$1,
			$2
		)
		ON CONFLICT (slug,tag) DO NOTHING
	`,
		b.Slug,
		t.Slug,
	)

	if err != nil {
		return fmt.Errorf("insert article_tag: %v", err)
	}

	return nil
}

// CreatePageTag creates a association record for Page tags.
func CreatePageTag(ctx context.Context, t *Tag, p *Page) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO "page_tag" (
			slug,
			tag
		)
		VALUES (
			$1,
			$2
		)
		ON CONFLICT (slug,tag) DO NOTHING
	`,
		p.Slug,
		t.Slug,
	)

	if err != nil {
		return fmt.Errorf("insert page_tag: %v", err)
	}

	return nil
}

// GetTags retrieves a list of tags
//encore:api public method=GET path=/tag
func GetTags(ctx context.Context) (*GetTagsResponse, error) {
	var (
		slug_description    sql.NullString
		feature_image       sql.NullString
		visibility          sql.NullString
		og_image            sql.NullString
		og_title            sql.NullString
		og_description      sql.NullString
		twitter_image       sql.NullString
		twitter_title       sql.NullString
		twitter_description sql.NullString
		meta_title          sql.NullString
		meta_description    sql.NullString
		accent_color        sql.NullString
	)
	rows, err := sqldb.Query(ctx, `
		SELECT
		slug,
		slug_name,
		slug_description,
		feature_image,
		visibility,
		og_image,
		og_title,
		og_description,
		twitter_image,
		twitter_title,
		twitter_description,
		meta_title,
		meta_description,
		accent_color,
		created_at,
		updated_at,
		slug_url
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
		err := rows.Scan(
			&t.Slug,
			&t.Name,
			&slug_description,
			&feature_image,
			&visibility,
			&og_image,
			&og_title,
			&og_description,
			&twitter_image,
			&twitter_title,
			&twitter_description,
			&meta_title,
			&meta_description,
			&accent_color,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.URL,
		)
		if err != nil {
			return &GetTagsResponse{
				Count: 0,
				Tags:  []*Tag{},
			}, err
		}
		t.Description = slug_description.String
		t.FeatureImage = feature_image.String
		t.Visibility = visibility.String
		t.OgImage = og_image.String
		t.OgTitle = og_title.String
		t.OgDescription = og_description.String
		t.TwitterImage = twitter_image.String
		t.TwitterTitle = twitter_title.String
		t.TwitterDescription = twitter_description.String
		t.MetaTitle = meta_title.String
		t.MetaDescription = meta_description.String
		t.AccentColor = accent_color.String
		q = append(q, &t)
		i = i + 1
	}
	return &GetTagsResponse{
		Count: i,
		Tags:  q,
	}, rows.Err()
}

// GetTagsBySlug retrieves a list of tags for a post
//encore:api public method=GET path=/tagsbypost/:slug
func GetTagsByPost(ctx context.Context, slug string) (*GetTagsResponse, error) {
	var (
		slug_description    sql.NullString
		feature_image       sql.NullString
		visibility          sql.NullString
		og_image            sql.NullString
		og_title            sql.NullString
		og_description      sql.NullString
		twitter_image       sql.NullString
		twitter_title       sql.NullString
		twitter_description sql.NullString
		meta_title          sql.NullString
		meta_description    sql.NullString
		accent_color        sql.NullString
	)
	rows, err := sqldb.Query(ctx, `
		SELECT
		slug,
		slug_name,
		slug_description,
		feature_image,
		visibility,
		og_image,
		og_title,
		og_description,
		twitter_image,
		twitter_title,
		twitter_description,
		meta_title,
		meta_description,
		accent_color,
		created_at,
		updated_at,
		slug_url
		FROM "tag"
		WHERE slug IN (select tag from article_tag where slug = $1)
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
		err := rows.Scan(
			&t.Slug,
			&t.Name,
			&slug_description,
			&feature_image,
			&visibility,
			&og_image,
			&og_title,
			&og_description,
			&twitter_image,
			&twitter_title,
			&twitter_description,
			&meta_title,
			&meta_description,
			&accent_color,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.URL,
		)
		if err != nil {
			return &GetTagsResponse{
				Count: 0,
				Tags:  []*Tag{},
			}, err
		}
		t.Description = slug_description.String
		t.FeatureImage = feature_image.String
		t.Visibility = visibility.String
		t.OgImage = og_image.String
		t.OgTitle = og_title.String
		t.OgDescription = og_description.String
		t.TwitterImage = twitter_image.String
		t.TwitterTitle = twitter_title.String
		t.TwitterDescription = twitter_description.String
		t.MetaTitle = meta_title.String
		t.MetaDescription = meta_description.String
		t.AccentColor = accent_color.String
		q = append(q, &t)
		i = i + 1
	}
	return &GetTagsResponse{
		Count: i,
		Tags:  q,
	}, rows.Err()
}

// GetTagsBySlug retrieves a list of tags for a post
//encore:api public method=GET path=/tagsbypage/:slug
func GetTagsByPage(ctx context.Context, slug string) (*GetTagsResponse, error) {
	var (
		slug_description    sql.NullString
		feature_image       sql.NullString
		visibility          sql.NullString
		og_image            sql.NullString
		og_title            sql.NullString
		og_description      sql.NullString
		twitter_image       sql.NullString
		twitter_title       sql.NullString
		twitter_description sql.NullString
		meta_title          sql.NullString
		meta_description    sql.NullString
		accent_color        sql.NullString
	)
	rows, err := sqldb.Query(ctx, `
		SELECT
		slug,
		slug_name,
		slug_description,
		feature_image,
		visibility,
		og_image,
		og_title,
		og_description,
		twitter_image,
		twitter_title,
		twitter_description,
		meta_title,
		meta_description,
		accent_color,
		created_at,
		updated_at,
		slug_url
		FROM "tag"
		WHERE slug IN (select tag from page_tag where slug = $1)
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
		err := rows.Scan(
			&t.Slug,
			&t.Name,
			&slug_description,
			&feature_image,
			&visibility,
			&og_image,
			&og_title,
			&og_description,
			&twitter_image,
			&twitter_title,
			&twitter_description,
			&meta_title,
			&meta_description,
			&accent_color,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.URL,
		)
		if err != nil {
			return &GetTagsResponse{
				Count: 0,
				Tags:  []*Tag{},
			}, err
		}
		t.Description = slug_description.String
		t.FeatureImage = feature_image.String
		t.Visibility = visibility.String
		t.OgImage = og_image.String
		t.OgTitle = og_title.String
		t.OgDescription = og_description.String
		t.TwitterImage = twitter_image.String
		t.TwitterTitle = twitter_title.String
		t.TwitterDescription = twitter_description.String
		t.MetaTitle = meta_title.String
		t.MetaDescription = meta_description.String
		t.AccentColor = accent_color.String
		q = append(q, &t)
		i = i + 1
	}
	return &GetTagsResponse{
		Count: i,
		Tags:  q,
	}, rows.Err()
}

// Post receives incoming post CRUD webhooks from ghost.
//encore:api public raw
func TagHook(w http.ResponseWriter, req *http.Request) {
	// ... operate on the raw HTTP request ...
	rlog.Info("received post hook")
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	rlog.Info(string(b))
	var t TagHookPayload
	err = json.Unmarshal(b, &t)
	if err != nil {
		rlog.Error("error unmarshalling post hook payload", "err", err)
	}

	rlog.Info(t.Tag.Current.Name)
	_, err = sqldb.Exec(context.Background(), `
		INSERT INTO "tag" (
			slug,
			slug_name,
			slug_description,
			feature_image,
			visibility,
			og_image,
			og_title,
			og_description,
			twitter_image,
			twitter_title,
			twitter_description,
			meta_title,
			meta_description,
			accent_color,
			created_at,
			updated_at,
			slug_url
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17
		)
		ON CONFLICT (slug) DO UPDATE
		SET
			slug_name = $2,
			slug_description = $3,
			feature_image = $4,
			visibility = $5,
			og_image = $6,
			og_title = $7,
			og_description = $8,
			twitter_image = $9,
			twitter_title = $10,
			twitter_description = $11,
			meta_title = $12,
			meta_description = $13,
			accent_color = $14,
			created_at = $15,
			updated_at = $16,
			slug_url = $17

	`,
		t.Tag.Current.Slug,
		t.Tag.Current.Name,
		t.Tag.Current.Description,
		t.Tag.Current.FeatureImage,
		t.Tag.Current.Visibility,
		t.Tag.Current.OgImage,
		t.Tag.Current.OgTitle,
		t.Tag.Current.OgDescription,
		t.Tag.Current.TwitterImage,
		t.Tag.Current.TwitterTitle,
		t.Tag.Current.TwitterDescription,
		t.Tag.Current.MetaTitle,
		t.Tag.Current.MetaDescription,
		t.Tag.Current.AccentColor,
		t.Tag.Current.CreatedAt,
		t.Tag.Current.UpdatedAt,
		t.Tag.Current.URL,
	)
	if err != nil {
		rlog.Error("error inserting data", "err", err)
	}
}

type TagHookPayload struct {
	Tag struct {
		Current struct {
			ID                 string    `json:"id"`
			Name               string    `json:"name"`
			Slug               string    `json:"slug"`
			Description        string    `json:"description"`
			FeatureImage       string    `json:"feature_image"`
			Visibility         string    `json:"visibility"`
			OgImage            string    `json:"og_image"`
			OgTitle            string    `json:"og_title"`
			OgDescription      string    `json:"og_description"`
			TwitterImage       string    `json:"twitter_image"`
			TwitterTitle       string    `json:"twitter_title"`
			TwitterDescription string    `json:"twitter_description"`
			MetaTitle          string    `json:"meta_title"`
			MetaDescription    string    `json:"meta_description"`
			CodeinjectionHead  string    `json:"codeinjection_head"`
			CodeinjectionFoot  string    `json:"codeinjection_foot"`
			CanonicalURL       string    `json:"canonical_url"`
			AccentColor        string    `json:"accent_color"`
			CreatedAt          time.Time `json:"created_at"`
			UpdatedAt          time.Time `json:"updated_at"`
			URL                string    `json:"url"`
		} `json:"current"`
		Previous struct {
			Description interface{} `json:"description"`
			UpdatedAt   time.Time   `json:"updated_at"`
		} `json:"previous"`
	} `json:"tag"`
}
