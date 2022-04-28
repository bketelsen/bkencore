// Package twitter integrates with the Twitter API to send tweets.
package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func twitterClient() *twitter.Client {
	config := oauth1.NewConfig(secrets.TwitterAPIKey, secrets.TwitterAPISecret)
	token := oauth1.NewToken(secrets.TwitterAccessToken, secrets.TwitterAccessSecret)
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	client := twitter.NewClient(httpClient)
	return client
}

var secrets struct {
	TwitterAccessToken  string
	TwitterAccessSecret string
	TwitterClientID     string
	TwitterClientSecret string
	TwitterRefreshToken string
	TwitterAPIKey       string
	TwitterAPISecret    string
}
