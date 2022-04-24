CREATE TABLE "scheduled_tweet" (
    -- id is a unique id for this scheduled tweet.
    id BIGSERIAL PRIMARY KEY,

	-- tweet_data is the tweet request payload to tweet.
    -- It is in the format of *twitter.TweetParams.
    tweet_data JSONB NOT NULL,

	-- tweet_id is the unique id of the tweet.
    -- It is non-null when the tweet has been successfully sent.
    tweet_id TEXT NULL,

    -- scheduled_at is the time when the tweet is scheduled for posting.
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- tweeted_at is the timestamp when the tweet was posted.
    sent_at TIMESTAMP WITH TIME ZONE NULL
);
