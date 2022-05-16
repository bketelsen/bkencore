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

type Page struct {
	ID                   string    `json:"id"`
	UUID                 string    `json:"uuid"`
	Title                string    `json:"title"`
	Slug                 string    `json:"slug"`
	HTML                 string    `json:"html"`
	Plaintext            string    `json:"plaintext"`
	FeatureImage         string    `json:"feature_image"`
	Featured             bool      `json:"featured"`
	Status               string    `json:"status"`
	Visibility           string    `json:"visibility"`
	EmailRecipientFilter string    `json:"email_recipient_filter"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	PublishedAt          time.Time `json:"published_at"`
	CustomExcerpt        string    `json:"custom_excerpt"`
	CanonicalURL         string    `json:"canonical_url"`
	PrimaryTag           string    `json:"primary_tag"`
	URL                  string    `json:"url"`
	Excerpt              string    `json:"excerpt"`
	ReadingTime          int       `json:"reading_time"`
	OgImage              string    `json:"og_image"`
	OgTitle              string    `json:"og_title"`
	OgDescription        string    `json:"og_description"`
	TwitterImage         string    `json:"twitter_image"`
	TwitterTitle         string    `json:"twitter_title"`
	TwitterDescription   string    `json:"twitter_description"`
	MetaTitle            string    `json:"meta_title"`
	MetaDescription      string    `json:"meta_description"`
	FeatureImageAlt      string    `json:"feature_image_alt"`
	FeatureImageCaption  string    `json:"feature_image_caption"`
}

type PageFull struct {
	ID                   string    `json:"id"`
	UUID                 string    `json:"uuid"`
	Title                string    `json:"title"`
	Slug                 string    `json:"slug"`
	HTML                 string    `json:"html"`
	Plaintext            string    `json:"plaintext"`
	FeatureImage         string    `json:"feature_image"`
	Featured             bool      `json:"featured"`
	Status               string    `json:"status"`
	Visibility           string    `json:"visibility"`
	EmailRecipientFilter string    `json:"email_recipient_filter"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	PublishedAt          time.Time `json:"published_at"`
	CustomExcerpt        string    `json:"custom_excerpt"`
	CanonicalURL         string    `json:"canonical_url"`
	URL                  string    `json:"url"`
	Excerpt              string    `json:"excerpt"`
	ReadingTime          int       `json:"reading_time"`
	OgImage              string    `json:"og_image"`
	OgTitle              string    `json:"og_title"`
	OgDescription        string    `json:"og_description"`
	TwitterImage         string    `json:"twitter_image"`
	TwitterTitle         string    `json:"twitter_title"`
	TwitterDescription   string    `json:"twitter_description"`
	MetaTitle            string    `json:"meta_title"`
	MetaDescription      string    `json:"meta_description"`
	FeatureImageAlt      string    `json:"feature_image_alt"`
	FeatureImageCaption  string    `json:"feature_image_caption"`
	PrimaryTag           *Tag      `json:"primary_tag"`
	Tags                 []*Tag    `json:"tags"`
}
type CreatePageParams struct {
	Published     bool   `json:"published,omitempty"`
	Title         string `json:"title,omitempty"`
	Subtitle      string `json:"subtitle,omitempty"`
	HeroText      string `json:"hero_text,omitempty"`
	Summary       string `json:"summary,omitempty"`
	Body          string `json:"body,omitempty"`
	FeaturedImage string `json:"featured_image,omitempty"` // empty string means no image
}

// GetPage retrieves a page by slug.
//encore:api public method=GET path=/page/:slug
func GetPage(ctx context.Context, slug string) (*PageFull, error) {
	var (
		b           PageFull
		primary_tag sql.NullString
	)
	err := sqldb.QueryRow(ctx, `
		SELECT
		slug,
		id,
		uuid,
		title,
		html,
		plaintext,
		feature_image,
		featured,
		status,
		visibility,
		created_at,
		updated_at,
		published_at,
		custom_excerpt,
		canonical_url,
		excerpt,
		reading_time,
		og_image,
		og_title,
		og_description,
		twitter_image,
		twitter_title,
		twitter_description,
		meta_title,
		meta_description,
		feature_image_alt,
		feature_image_caption,
		primary_tag,
		url
		FROM "page"
		WHERE slug = $1
	`, slug).Scan(&b.Slug,
		&b.ID,
		&b.UUID,
		&b.Title,
		&b.HTML,
		&b.Plaintext,
		&b.FeatureImage,
		&b.Featured,
		&b.Status,
		&b.Visibility,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.PublishedAt,
		&b.CustomExcerpt,
		&b.CanonicalURL,
		&b.Excerpt,
		&b.ReadingTime,
		&b.OgImage,
		&b.OgTitle,
		&b.OgDescription,
		&b.TwitterImage,
		&b.TwitterTitle,
		&b.TwitterDescription,
		&b.MetaTitle,
		&b.MetaDescription,
		&b.FeatureImageAlt,
		&b.FeatureImageCaption,
		&primary_tag,
		&b.URL)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.NotFound,
			Message: "page not found",
		}
	}
	if primary_tag.Valid && primary_tag.String != "" {
		t, err := GetTag(ctx, primary_tag.String)
		if err != nil {
			return nil, err
		}
		b.PrimaryTag = t
	}
	tresp, err := GetTagsByPage(ctx, b.Slug)
	if err != nil {
		return nil, err
	}
	b.Tags = tresp.Tags
	return &b, nil
}

/*
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
*/

// Post receives incoming post CRUD webhooks from ghost.
//encore:api public raw
func PageHook(w http.ResponseWriter, req *http.Request) {
	rlog.Info("received page hook")
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	rlog.Info(string(b))
	var p PageHookPayload
	err = json.Unmarshal(b, &p)
	if err != nil {
		rlog.Error("error unmarshalling page hook payload", "err", err)
	}

	_, err = sqldb.Exec(req.Context(), `
		INSERT INTO "page" (
			slug,
			id,
			uuid,
			title,
			html,
			plaintext,
			feature_image,
			featured,
			status,
			visibility,
			created_at,
			updated_at,
			published_at,
			custom_excerpt,
			canonical_url,
			excerpt,
			reading_time,
			og_image,
			og_title,
			og_description,
			twitter_image,
			twitter_title,
			twitter_description,
			meta_title,
			meta_description,
			feature_image_alt,
			feature_image_caption,
			primary_tag,
			url
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
			$17,
			$18,
			$19,
			$20,
			$21,
			$22,
			$23,
			$24,
			$25,
			$26,
			$27,
			$28,
			$29
		)
		ON CONFLICT (slug) DO UPDATE
		SET
			id = $2,
			uuid = $3,
			title = $4,
			html = $5,
			plaintext = $6,
			feature_image = $7,
			featured = $8,
			status = $9,
			visibility = $10,
			created_at = $11,
			updated_at = $12,
			published_at = $13,
			custom_excerpt = $14,
			canonical_url = $15,
			excerpt = $16,
			reading_time = $17,
			og_image = $18,
			og_title = $19,
			og_description = $20,
			twitter_image = $21,
			twitter_title = $22,
			twitter_description = $23,
			meta_title 	= $24,
			meta_description = $25,
			feature_image_alt = $26,
			feature_image_caption = $27,
			primary_tag = $28,
			url = $29
	`,
		p.Page.Current.Slug,
		p.Page.Current.ID,
		p.Page.Current.UUID,
		p.Page.Current.Title,
		p.Page.Current.HTML,
		p.Page.Current.Plaintext,
		p.Page.Current.FeatureImage,
		p.Page.Current.Featured,
		p.Page.Current.Status,
		p.Page.Current.Visibility,
		p.Page.Current.CreatedAt,
		p.Page.Current.UpdatedAt,
		p.Page.Current.PublishedAt,
		p.Page.Current.CustomExcerpt,
		p.Page.Current.CanonicalURL,
		p.Page.Current.Excerpt,
		p.Page.Current.ReadingTime,
		p.Page.Current.OgImage,
		p.Page.Current.OgTitle,
		p.Page.Current.OgDescription,
		p.Page.Current.TwitterImage,
		p.Page.Current.TwitterTitle,
		p.Page.Current.TwitterDescription,
		p.Page.Current.MetaTitle,
		p.Page.Current.MetaDescription,
		p.Page.Current.FeatureImageAlt,
		p.Page.Current.FeatureImageCaption,
		p.Page.Current.PrimaryTag.Slug,
		p.Page.Current.URL,
	)
	if err != nil {
		rlog.Error("error inserting data", "err", err)
	}

	for _, tag := range p.Page.Current.Tags {
		var t Tag
		t.AccentColor = tag.AccentColor
		t.CreatedAt = tag.CreatedAt
		t.Description = tag.Description
		t.FeatureImage = tag.FeatureImage
		t.MetaDescription = tag.MetaDescription
		t.MetaTitle = tag.MetaTitle
		t.Name = tag.Name
		t.OgDescription = tag.OgDescription
		t.OgImage = tag.OgImage
		t.OgTitle = tag.OgTitle
		t.Slug = tag.Slug
		t.TwitterDescription = tag.TwitterDescription
		t.TwitterImage = tag.TwitterImage
		t.TwitterTitle = tag.TwitterTitle
		t.TwitterDescription = tag.TwitterDescription
		t.URL = tag.URL
		t.UpdatedAt = tag.UpdatedAt

		err := CreateTag(req.Context(), &t)
		if err != nil {
			rlog.Error("error inserting tag", "err", err)
		}
		err = CreatePageTag(req.Context(), &t, &Page{Slug: p.Page.Current.Slug})
		if err != nil {
			rlog.Error("error inserting page_tag", "err", err)
		}
	}

}

type PageHookPayload struct {
	Page struct {
		Current struct {
			ID                string    `json:"id"`
			UUID              string    `json:"uuid"`
			Title             string    `json:"title"`
			Slug              string    `json:"slug"`
			Mobiledoc         string    `json:"mobiledoc"`
			HTML              string    `json:"html"`
			CommentID         string    `json:"comment_id"`
			Plaintext         string    `json:"plaintext"`
			FeatureImage      string    `json:"feature_image"`
			Featured          bool      `json:"featured"`
			Status            string    `json:"status"`
			Visibility        string    `json:"visibility"`
			CreatedAt         time.Time `json:"created_at"`
			UpdatedAt         time.Time `json:"updated_at"`
			PublishedAt       time.Time `json:"published_at"`
			CustomExcerpt     string    `json:"custom_excerpt"`
			CodeinjectionHead string    `json:"codeinjection_head"`
			CodeinjectionFoot string    `json:"codeinjection_foot"`
			CustomTemplate    string    `json:"custom_template"`
			CanonicalURL      string    `json:"canonical_url"`
			Tags              []struct {
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
			} `json:"tags"`

			PrimaryTag struct {
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
			} `json:"primary_tag"`
			URL                 string `json:"url"`
			Excerpt             string `json:"excerpt"`
			ReadingTime         int    `json:"reading_time"`
			OgImage             string `json:"og_image"`
			OgTitle             string `json:"og_title"`
			OgDescription       string `json:"og_description"`
			TwitterImage        string `json:"twitter_image"`
			TwitterTitle        string `json:"twitter_title"`
			TwitterDescription  string `json:"twitter_description"`
			MetaTitle           string `json:"meta_title"`
			MetaDescription     string `json:"meta_description"`
			Frontmatter         string `json:"frontmatter"`
			FeatureImageAlt     string `json:"feature_image_alt"`
			FeatureImageCaption string `json:"feature_image_caption"`
		} `json:"current"`
	} `json:"page"`
}
