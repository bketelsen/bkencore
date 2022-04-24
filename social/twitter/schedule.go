// Package twitter integrates with the Twitter API to send tweets.
package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/cron"
	"encore.dev/storage/sqldb"
)

type TweetParams struct {
	// Text is the text to tweet.
	Text string
}

type TweetResponse struct {
	// ID is the tweet id.
	ID string
}

// Tweet writes a mock tweet to the database.
//encore:api auth method=POST path=/twitter/tweet
func Tweet(ctx context.Context, p *TweetParams) (*TweetResponse, error) {
	var id int64
	err := sqldb.QueryRow(ctx, `
		INSERT INTO mock_tweet (body)
		VALUES ($1)
		RETURNING id
	`, p.Text).Scan(&id)
	return &TweetResponse{ID: fmt.Sprintf("mock-%d", id)}, err
}

// Tweet sends a tweet using the Twitter API.
//encore:api auth method=POST path=/twitter/tweet/for-real
func TweetForReal(ctx context.Context, p *TweetParams) (*TweetResponse, error) {
	eb := errs.B()
	client := httpClient(ctx)
	data, _ := json.Marshal(map[string]any{
		"text": p.Text,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.twitter.com/2/tweets", bytes.NewReader(data))
	if err != nil {
		return nil, eb.Cause(err).Msg("unable to create request").Err()
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, eb.Cause(err).Msg("unable to make request").Err()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, eb.Msgf("got non-201 response from Twitter API: %s (%s)", resp.Status, body).Err()
	}

	var respData struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, eb.Cause(err).Msg("unable to decode response").Err()
	}
	return &TweetResponse{ID: respData.Data.ID}, nil
}

type ScheduleParams struct {
	// Tweet is the tweet to schedule.
	Tweet *TweetParams

	// SendAt is the time to send it at.
	SendAt time.Time
}

type ScheduleResponse struct {
	// ID is the unique id of the scheduled tweet.
	ID int64
}

// Schedule schedules a tweet to be posted at a certain time.
//encore:api private method=POST path=/twitter/schedule
func Schedule(ctx context.Context, p *ScheduleParams) (*ScheduleResponse, error) {
	eb := errs.B()
	data, err := json.Marshal(p.Tweet)
	if err != nil {
		return nil, eb.Cause(err).Msg("unable to marshal tweet").Err()
	}

	var id int64
	err = sqldb.QueryRow(ctx, `
		INSERT INTO scheduled_tweet (tweet_data, scheduled_at)
		VALUES ($1, $2)
		RETURNING id
	`, data, p.SendAt).Scan(&id)
	if err != nil {
		return nil, eb.Cause(err).Msg("unable to insert row").Err()
	}
	return &ScheduleResponse{ID: id}, nil
}

// SendDue posts tweets that are due.
//encore:api auth method=POST path=/twitter/send-due
func SendDue(ctx context.Context) error {
	tweet, err := queryDueTweet(ctx, time.Now())
	if errors.Is(err, sqldb.ErrNoRows) {
		return nil
	} else if err != nil {
		return err
	}

	resp, err := Tweet(ctx, tweet.Tweet)
	if err != nil {
		return err
	}
	// Use a separate context in case the parent ctx has been canceled.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err = sqldb.Exec(ctx, `
		UPDATE scheduled_tweet
		SET sent_at = NOW(), tweet_id = $2
		WHERE id = $1
	`, tweet.ID, resp.ID)
	return err
}

// scheduledTweet represents a database row for scheduled tweets.
type scheduledTweet struct {
	ID          int64
	Tweet       *TweetParams
	TweetID     *string
	ScheduledAt time.Time
	SentAt      *time.Time
}

// queryDueTweet reports a tweet that is due to be sent at time t.
// It returns a maximum of one tweet to avoid spamming the Twitter timeline
// with multiple tweets at the same time.
// If no tweet is due it reports sqldb.ErrNoRows.
func queryDueTweet(ctx context.Context, t time.Time) (*scheduledTweet, error) {
	var (
		tweet  scheduledTweet
		data   []byte
		params TweetParams
	)
	err := sqldb.QueryRow(ctx, `
		SELECT id, tweet_data, tweet_id, scheduled_at, sent_at
		FROM scheduled_tweet
		WHERE sent_at IS NULL AND scheduled_at <= $1
		LIMIT 1
	`, t).Scan(&tweet.ID, &data, &tweet.TweetID, &tweet.ScheduledAt, &tweet.SentAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &params); err != nil {
		return nil, err
	}
	tweet.Tweet = &params
	return &tweet, nil
}

// Send tweets due to be delivered every minute.
var _ = cron.NewJob("send-due-tweets", cron.JobConfig{
	Title:    "Send due tweets",
	Every:    1 * cron.Minute,
	Endpoint: SendDue,
})
