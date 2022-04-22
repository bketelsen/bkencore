// Package twitter integrates with the Twitter API to send tweets.
package twitter

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type TweetParams struct {
	// Text is the text to tweet.
	Text string
}

type TweetResponse struct {
	// ID is the tweet id.
	ID string
}

// Tweet sends a tweet using the Twitter API.
//encore:api public method=POST path=/twitter/tweet
func Tweet(ctx context.Context, p *TweetParams) (*TweetResponse, error) {
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
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, eb.Cause(err).Msg("unable to decode response").Err()
	}
	return &TweetResponse{ID: respData.ID}, nil
}

// httpClient returns an HTTP client to make authenticated requests to the Twitter API.
func httpClient(ctx context.Context) *http.Client {
	return cfg.Client(ctx, &oauth2.Token{
		TokenType:    "bearer",
		RefreshToken: secrets.TwitterRefreshToken,
		Expiry:       time.Now().Add(-1 * time.Hour),
	})
}

// OAuthBegin begins an OAuth handshake.
//encore:api public raw method=GET path=/twitter/oauth/begin
func OAuthBegin(w http.ResponseWriter, req *http.Request) {
	state, challenge, err := generateStateAndChallenge()
	if err != nil {
		rlog.Error("unable to generate random state", "err", err)
		http.Error(w, "unable to generate state", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "oauth2_challenge",
		Value: challenge,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "oauth2_state",
		Value: state,
	})

	url := cfg.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)
	http.Redirect(w, req, url, http.StatusFound)
}

// OAuthToken retrieves an OAuth token.
//encore:api public raw method=GET path=/twitter/oauth/token
func OAuthToken(w http.ResponseWriter, req *http.Request) {
	wantState := cookieValue(req, "state")
	challenge := cookieValue(req, "challenge")
	if got := req.FormValue("state"); got != wantState {
		rlog.Error("state mismatch", "got", got, "want", wantState)
		http.Error(w, "bad state", http.StatusBadRequest)
		return
	}

	tok, err := cfg.Exchange(req.Context(), req.FormValue("code"),
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_verifier", challenge),
	)
	if err != nil {
		rlog.Error("unable to exchange code", "err", err)
		http.Error(w, "unable to exchange code", http.StatusInternalServerError)
		return
	}
	rlog.Info("got token", "tok", tok)
}

func generateStateAndChallenge() (state, challenge string, err error) {
	var data [32]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", "", err
	}
	state = base64.URLEncoding.EncodeToString(data[:16])
	challenge = base64.URLEncoding.EncodeToString(data[16:])
	return state, challenge, nil
}

func cookieValue(req *http.Request, name string) string {
	if cookie, err := req.Cookie(name); err == nil {
		return cookie.Value
	}
	return ""
}

var secrets struct {
	TwitterClientID     string
	TwitterClientSecret string
	TwitterRefreshToken string
}

var cfg = oauth2.Config{
	ClientID:     secrets.TwitterClientID,
	ClientSecret: secrets.TwitterClientSecret,
	Endpoint: oauth2.Endpoint{
		AuthURL:   "https://twitter.com/i/oauth2/authorize",
		TokenURL:  "https://api.twitter.com/2/oauth2/token",
		AuthStyle: oauth2.AuthStyleInHeader,
	},
	RedirectURL: "http://localhost:4000/twitter/oauth/token", // TODO
	Scopes:      []string{"tweet.read", "tweet.write", "users.read", "offline.access"},
}
