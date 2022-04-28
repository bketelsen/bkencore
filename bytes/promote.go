package bytes

import (
	"context"
	"fmt"
	"time"

	"encore.app/social/twitter"
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
//encore:api auth method=POST path=/bytes/:id/promote
func Promote(ctx context.Context, id int64, p *PromoteParams) error {
	eb := errs.B().Meta("id", id)
	byte, err := Get(ctx, id)
	if err != nil {
		return eb.Cause(err).Msg("unable to get byte ").Err()
	}
	sendAt := time.Now() // TODO factor in p.Schedule

	// Schedule twitter
	_, err = twitter.Schedule(ctx, &twitter.ScheduleParams{
		SendAt: sendAt,
		Tweet: &twitter.TweetParams{
			// TODO very placeholder tweet text
			Text: fmt.Sprintf("Quick Byte: %s - %s \n\n%s", byte.Title, byte.Summary, byte.URL),
		},
	})
	if err != nil {
		return eb.Cause(err).Msg("unable to schedule twitter post").Err()
	}

	return nil
}
