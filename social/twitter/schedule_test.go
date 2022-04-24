package twitter

import (
	"context"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp/cmpopts"

	"encore.dev/storage/sqldb"
)

func TestSchedule(t *testing.T) {
	c := qt.New(t)
	ctx := context.Background()

	// Schedule a tweet.
	p := &ScheduleParams{
		Tweet:  &TweetParams{Text: "hello, twitter world!"},
		SendAt: time.Now().Add(1 * time.Hour),
	}

	resp, err := Schedule(ctx, p)
	c.Assert(err, qt.IsNil)
	c.Assert(resp.ID, qt.Not(qt.Equals), 0)

	// Verify the database row exists and looks like we expect.
	var r struct {
		ID          int64
		TweetData   []byte
		TweetID     *string
		ScheduledAt time.Time
		SentAt      *time.Time
	}
	err = sqldb.QueryRow(ctx, `
		SELECT tweet_data, tweet_id, scheduled_at, sent_at
		FROM scheduled_tweet
		WHERE id = $1
	`, resp.ID).Scan(&r.TweetData, &r.TweetID, &r.ScheduledAt, &r.SentAt)
	c.Assert(err, qt.IsNil)
	c.Assert(r.TweetData, qt.JSONEquals, p.Tweet)
	c.Assert(r.TweetID, qt.IsNil)
	c.Assert(r.SentAt, qt.IsNil)
	c.Assert(r.ScheduledAt, qt.CmpEquals(cmpopts.EquateApproxTime(time.Millisecond)), p.SendAt)
}
