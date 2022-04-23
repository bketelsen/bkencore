package blog

import (
	"context"
	"fmt"
	"time"

	"encore.app/email"
	"encore.app/social/twitter"
	"encore.app/url"
	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
)

type ScheduleType string

const (
	ScheduleAuto ScheduleType = "auto"
	ScheduleNow  ScheduleType = "now"
)

type PromoteParams struct {
	// Schedule decides how the promotion should be scheduled.
	// Valid values are "auto" for scheduling it at a suitable time
	// based on the current posting schedule, and "now" to schedule it immediately.
	Schedule ScheduleType
}

// Promote schedules the promotion a blog post.
//encore:api public method=POST path=/blog/:slug/promote
func Promote(ctx context.Context, slug string, p *PromoteParams) error {
	eb := errs.B().Meta("slug", slug)
	post, err := GetBlogPost(ctx, slug)
	if err != nil {
		return eb.Cause(err).Msg("unable to get blog post").Err()
	}

	// Generate a short URL for the blog post
	short, err := url.Shorten(ctx, &url.ShortenParams{
		URL: "/blog/" + slug, // TODO
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to generate short url").Err()
	}
	shortURL := fmt.Sprintf("/s/%s", short.ID) // TODO make proper

	// Create a template for the email
	err = email.CreateTemplate(ctx, slug, &email.CreateTemplateParams{
		Sender:   "todo@example.org", // TODO
		Subject:  post.Title,
		BodyText: post.Body,
		BodyHTML: post.BodyRendered,
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to create email template").Err()
	}

	// Schedule emails
	sendAt := time.Now() // TODO factor in p.Schedule
	subscribers, err := getSubscribers(ctx)
	if err != nil {
		return eb.Cause(err).Msg("unable to get email subscribers").Err()
	}
	_, err = email.Schedule(ctx, &email.ScheduleParams{
		TemplateID:     slug,
		EmailAddresses: subscribers,
		SendAt:         &sendAt,
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to get email subscribers").Err()
	}

	// Schedule twitter
	_, err = twitter.Schedule(ctx, &twitter.ScheduleParams{
		SendAt: sendAt,
		Tweet: &twitter.TweetParams{
			// TODO very placeholder tweet text
			Text: fmt.Sprintf("New blog post: %s\n\n%s",
				post.Title, shortURL),
		},
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to schedule twitter post").Err()
	}

	return nil
}

type SubscribeParams struct {
	Email string
}

// Subscribe subscribes to the email newsletter for a given email.
//encore:api public method=POST path=/subscribe
func Subscribe(ctx context.Context, p *SubscribeParams) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO email_subscriber (email)
		VALUES ($1)
		ON CONFLICT (email) DO NOTHING
	`, p.Email)
	return err
}

// getSubscribers returns all the subscribers to the email newsletter.
func getSubscribers(ctx context.Context) (emails []string, err error) {
	rows, err := sqldb.Query(ctx, `
		SELECT email
		FROM email_subscriber
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, rows.Err()
}
