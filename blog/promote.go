package blog

import (
	"context"
	"fmt"
	"time"

	"encore.app/email"
	"encore.app/social/twitter"
	"encore.app/url"
	"encore.dev/beta/errs"
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
//encore:api auth method=POST path=/blog/:slug/promote
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

	// Create a template for the email
	err = email.CreateTemplate(ctx, slug, &email.CreateTemplateParams{
		Sender:   "Brian Ketelsen <me@brian.dev>", // TODO
		Subject:  post.Title,
		BodyText: post.Body,
		BodyHTML: post.BodyRendered,
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to create email template").Err()
	}

	// Schedule emails
	sendAt := time.Now() // TODO factor in p.Schedule
	_, err = email.ScheduleAll(ctx, &email.ScheduleAllParams{
		TemplateID: slug,
		SendAt:     &sendAt,
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to get email subscribers").Err()
	}

	// Schedule twitter
	_, err = twitter.Schedule(ctx, &twitter.ScheduleParams{
		SendAt: sendAt,
		Tweet: &twitter.TweetParams{
			// TODO very placeholder tweet text
			Text: fmt.Sprintf("%s\n\n%s", post.Summary, short.ShortURL),
		},
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to schedule twitter post").Err()
	}

	return nil
}
