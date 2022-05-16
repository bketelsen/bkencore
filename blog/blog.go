package blog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"

	"encore.dev/rlog"
)

type BlogPost struct {
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
type BlogPostFull struct {
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

type GetBlogPostsParams struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetBlogPostsResponse struct {
	Count     int             `json:"count,omitempty"`
	BlogPosts []*BlogPostFull `json:"blog_posts"`
}

// GetBlogPost retrieves a blog post by slug.
//encore:api public method=GET path=/blog/:slug
func GetBlogPost(ctx context.Context, slug string) (*BlogPostFull, error) {
	var (
		b           BlogPostFull
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
		FROM "article"
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
			Message: "article not found",
		}
	}
	if primary_tag.Valid {
		t, err := GetTag(ctx, primary_tag.String)
		if err != nil {
			return nil, err
		}
		b.PrimaryTag = t
	}
	tresp, err := GetTagsByPost(ctx, b.Slug)
	if err != nil {
		return nil, err
	}
	b.Tags = tresp.Tags
	return &b, nil
}

/*
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
*/

// GetBlogPosts retrieves a list of blog posts with
// optional limit and offset.
//encore:api public method=GET path=/blog
func GetBlogPosts(ctx context.Context, params *GetBlogPostsParams) (*GetBlogPostsResponse, error) {
	rows, err := sqldb.Query(ctx, `
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
		FROM "article"
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return &GetBlogPostsResponse{
			Count:     0,
			BlogPosts: []*BlogPostFull{},
		}, err
	}
	defer rows.Close()

	var q []*BlogPostFull
	var i = 0
	for rows.Next() {
		var (
			b BlogPostFull
			t sql.NullString
		)
		err := rows.Scan(&b.Slug,
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
			&t,
			&b.URL)
		if err != nil {
			return &GetBlogPostsResponse{
				Count:     0,
				BlogPosts: []*BlogPostFull{},
			}, err
		}
		if t.Valid && t.String != "" {
			t, err := GetTag(ctx, t.String)
			if err != nil {
				return &GetBlogPostsResponse{
					Count:     0,
					BlogPosts: []*BlogPostFull{},
				}, err
			}
			b.PrimaryTag = t
		}
		tresp, err := GetTagsByPost(ctx, b.Slug)
		if err != nil {
			return nil, err
		}
		b.Tags = tresp.Tags

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

// Post receives incoming post CRUD webhooks from ghost.
//encore:api public raw
func PostHook(w http.ResponseWriter, req *http.Request) {
	// ... operate on the raw HTTP request ...
	ctx := req.Context()
	rlog.Info("received post hook")
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
	rlog.Info(string(b))
	var p PostHookPayload
	err = json.Unmarshal(b, &p)
	if err != nil {
		rlog.Error("error unmarshalling post hook payload", "err", err)
	}

	_, err = sqldb.Exec(ctx, `
		INSERT INTO "article" (
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
		p.Post.Current.Slug,
		p.Post.Current.ID,
		p.Post.Current.UUID,
		p.Post.Current.Title,
		p.Post.Current.HTML,
		p.Post.Current.Plaintext,
		p.Post.Current.FeatureImage,
		p.Post.Current.Featured,
		p.Post.Current.Status,
		p.Post.Current.Visibility,
		p.Post.Current.CreatedAt,
		p.Post.Current.UpdatedAt,
		p.Post.Current.PublishedAt,
		p.Post.Current.CustomExcerpt,
		p.Post.Current.CanonicalURL,
		p.Post.Current.Excerpt,
		p.Post.Current.ReadingTime,
		p.Post.Current.OgImage,
		p.Post.Current.OgTitle,
		p.Post.Current.OgDescription,
		p.Post.Current.TwitterImage,
		p.Post.Current.TwitterTitle,
		p.Post.Current.TwitterDescription,
		p.Post.Current.MetaTitle,
		p.Post.Current.MetaDescription,
		p.Post.Current.FeatureImageAlt,
		p.Post.Current.FeatureImageCaption,
		p.Post.Current.PrimaryTag.Slug,
		p.Post.Current.URL,
	)
	if err != nil {
		rlog.Error("error inserting data", "err", err)
	}

	for _, tag := range p.Post.Current.Tags {
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

		err := CreateTag(ctx, &t)
		if err != nil {
			rlog.Error("error inserting tag", "err", err)
		}
		err = CreatePostTag(ctx, &t, &BlogPost{Slug: p.Post.Current.Slug})
		if err != nil {
			rlog.Error("error inserting article_tag", "err", err)
		}
	}
}

type PostHookPayload struct {
	Post struct {
		Current struct {
			ID                   string    `json:"id"`
			UUID                 string    `json:"uuid"`
			Title                string    `json:"title"`
			Slug                 string    `json:"slug"`
			Mobiledoc            string    `json:"mobiledoc"`
			HTML                 string    `json:"html"`
			CommentID            string    `json:"comment_id"`
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
			CodeinjectionHead    string    `json:"codeinjection_head"`
			CodeinjectionFoot    string    `json:"codeinjection_foot"`
			CustomTemplate       string    `json:"custom_template"`
			CanonicalURL         string    `json:"canonical_url"`
			NewsletterID         string    `json:"newsletter_id"`
			Tags                 []struct {
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
			PrimaryTag          Tag    `json:"primary_tag"`
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
			EmailSubject        string `json:"email_subject"`
			Frontmatter         string `json:"frontmatter"`
			FeatureImageAlt     string `json:"feature_image_alt"`
			FeatureImageCaption string `json:"feature_image_caption"`
			EmailOnly           bool   `json:"email_only"`
		} `json:"current"`
		Previous struct {
			Mobiledoc string    `json:"mobiledoc"`
			UpdatedAt time.Time `json:"updated_at"`
			HTML      string    `json:"html"`
			Plaintext string    `json:"plaintext"`
		} `json:"previous"`
	} `json:"post"`
}
