// Package twitter integrates with the Twitter API to send tweets.
package twitter

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"encore.dev/rlog"
)

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
	wantState := cookieValue(req, "oauth2_state")
	challenge := cookieValue(req, "oauth2_challenge")
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

// httpClient returns an HTTP client to make authenticated requests to the Twitter API.
func httpClient(ctx context.Context) *http.Client {
	return cfg.Client(ctx, &oauth2.Token{
		TokenType:    "bearer",
		RefreshToken: secrets.TwitterRefreshToken,
		Expiry:       time.Now().Add(-1 * time.Hour),
	})
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
