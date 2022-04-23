package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an API client for the devweek-k65i Encore application.
type Client struct {
	Blog    BlogClient
	Email   EmailClient
	Hello   HelloClient
	Twitter TwitterClient
	Url     UrlClient
}

// BaseURL is the base URL for calling the Encore application's API.
type BaseURL string

const Local BaseURL = "http://localhost:4000"

// Environment returns a BaseURL for calling the cloud environment with the given name.
func Environment(name string) BaseURL {
	return BaseURL(fmt.Sprintf("https://%s-devweek-k65i.encr.app", name))
}

// Option allows you to customise the baseClient used by the Client
type Option = func(client *baseClient) error

// New returns a Client for calling the public and authenticated APIs of your Encore application.
// You can customize the behaviour of the client using the given Option functions, such as WithHTTPClient or WithAuthToken.
func New(target BaseURL, options ...Option) (*Client, error) {
	// Parse the base URL where the Encore application is being hosted
	baseURL, err := url.Parse(string(target))
	if err != nil {
		return nil, fmt.Errorf("unable to parse base url: %w", err)
	}

	// Create a client with sensible defaults
	base := &baseClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		userAgent:  "devweek-k65i-Generated-Client (Encore/v1.0.0)",
	}

	// Apply any given options
	for _, option := range options {
		if err := option(base); err != nil {
			return nil, fmt.Errorf("unable to apply client option: %w", err)
		}
	}

	return &Client{
		Blog:    &blogClient{base},
		Email:   &emailClient{base},
		Hello:   &helloClient{base},
		Twitter: &twitterClient{base},
		Url:     &urlClient{base},
	}, nil
}

// WithHTTPClient can be used to configure the underlying HTTP client used when making API calls.
//
// Defaults to http.DefaultClient
func WithHTTPClient(client HTTPDoer) Option {
	return func(base *baseClient) error {
		base.httpClient = client
		return nil
	}
}

// WithAuthToken allows you to set the auth token to be used for each request
func WithAuthToken(token string) Option {
	return func(base *baseClient) error {
		base.tokenGenerator = func(_ context.Context) (string, error) {
			return token, nil
		}
		return nil
	}
}

// WithAuthFunc allows you to pass a function which is called for each request to return an access token.
func WithAuthFunc(tokenGenerator func(ctx context.Context) (string, error)) Option {
	return func(base *baseClient) error {
		base.tokenGenerator = tokenGenerator
		return nil
	}
}

type BlogBlogPost struct {
	Slug         string
	CreatedAt    time.Time `qs:"created_at"`
	Published    bool      `yaml:"published"`
	ModifiedAt   time.Time `qs:"modified_at"`
	Title        string    `json:"title" yaml:"title"`
	Summary      string    `yaml:"summary"`
	Body         string
	BodyRendered string `qs:"body_rendered"`
}

type BlogCreateBlogPostParams struct {
	Slug      string
	Published bool
	Title     string
	Summary   string
	Body      string
}

type BlogGetBlogPostsParams struct {
	Limit  int
	Offset int
}

type BlogGetBlogPostsResponse struct {
	Count     int
	BlogPosts []BlogBlogPost `qs:"blog_posts"`
}

// BlogClient Provides you access to call public and authenticated APIs on blog. The concrete implementation is blogClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type BlogClient interface {

	// CreateBlogPost creates a new blog post.
	CreateBlogPost(ctx context.Context, params BlogCreateBlogPostParams) error

	// GetBlogPost retrieves a blog post by slug.
	GetBlogPost(ctx context.Context, slug string) (BlogBlogPost, error)

	// GetBlogPosts retrieves a list of blog posts with
	// optional limit and offset.
	GetBlogPosts(ctx context.Context, params BlogGetBlogPostsParams) (BlogGetBlogPostsResponse, error)
}

type blogClient struct {
	base *baseClient
}

var _ BlogClient = (*blogClient)(nil)

// CreateBlogPost creates a new blog post.
func (c *blogClient) CreateBlogPost(ctx context.Context, params BlogCreateBlogPostParams) error {
	_, err := callAPI[struct{}](ctx, c.base, "POST", "/blog.CreateBlogPost", params)
	return err
}

// GetBlogPost retrieves a blog post by slug.
func (c *blogClient) GetBlogPost(ctx context.Context, slug string) (BlogBlogPost, error) {
	return callAPI[BlogBlogPost](ctx, c.base, "GET", fmt.Sprintf("/blog/%s", slug), nil)
}

// GetBlogPosts retrieves a list of blog posts with
// optional limit and offset.
func (c *blogClient) GetBlogPosts(ctx context.Context, params BlogGetBlogPostsParams) (BlogGetBlogPostsResponse, error) {
	queryString := url.Values{
		"limit":  []string{fmt.Sprint(params.Limit)},
		"offset": []string{fmt.Sprint(params.Offset)},
	}
	return callAPI[BlogGetBlogPostsResponse](ctx, c.base, "GET", fmt.Sprintf("/blog?%s", queryString.Encode()), nil)
}

type EmailCreateTemplateParams struct {
	Sender   string // sender email
	Subject  string // subject line to use
	BodyText string `qs:"body_text"` // plaintext body
	BodyHTML string `qs:"body_html"` // html body
}

type EmailUnsubscribeParams struct {
	Token string // Token is the unsubscribe token in to the email.
}

// EmailClient Provides you access to call public and authenticated APIs on email. The concrete implementation is emailClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type EmailClient interface {

	// CreateTemplate creates an email template.
	// If the template with that id already exists it is updated.
	CreateTemplate(ctx context.Context, id string, params EmailCreateTemplateParams) error

	// Unsubscribe unsubscribes the user from the email list.
	Unsubscribe(ctx context.Context, params EmailUnsubscribeParams) error
}

type emailClient struct {
	base *baseClient
}

var _ EmailClient = (*emailClient)(nil)

// CreateTemplate creates an email template.
// If the template with that id already exists it is updated.
func (c *emailClient) CreateTemplate(ctx context.Context, id string, params EmailCreateTemplateParams) error {
	_, err := callAPI[struct{}](ctx, c.base, "PUT", fmt.Sprintf("/email/templates/%s", id), params)
	return err
}

// Unsubscribe unsubscribes the user from the email list.
func (c *emailClient) Unsubscribe(ctx context.Context, params EmailUnsubscribeParams) error {
	_, err := callAPI[struct{}](ctx, c.base, "POST", "/email/unsubscribe", params)
	return err
}

type HelloResponse struct {
	Message string
}

// HelloClient Provides you access to call public and authenticated APIs on hello. The concrete implementation is helloClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type HelloClient interface {

	// This is a simple REST API that responds with a personalized greeting.
	// To call it, run in your terminal:
	//
	//     curl http://localhost:4000/hello/World
	World(ctx context.Context, name string) (HelloResponse, error)
}

type helloClient struct {
	base *baseClient
}

var _ HelloClient = (*helloClient)(nil)

// This is a simple REST API that responds with a personalized greeting.
// To call it, run in your terminal:
//     curl http://localhost:4000/hello/World
func (c *helloClient) World(ctx context.Context, name string) (HelloResponse, error) {
	return callAPI[HelloResponse](ctx, c.base, "GET", fmt.Sprintf("/hello/%s", name), nil)
}

type TwitterTweetParams struct {
	Text string // Text is the text to tweet.
}

type TwitterTweetResponse struct {
	ID string // ID is the tweet id.
}

// TwitterClient Provides you access to call public and authenticated APIs on twitter. The concrete implementation is twitterClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type TwitterClient interface {

	// OAuthBegin begins an OAuth handshake.
	OAuthBegin(ctx context.Context, request *http.Request) (*http.Response, error)

	// OAuthToken retrieves an OAuth token.
	OAuthToken(ctx context.Context, request *http.Request) (*http.Response, error)

	// Tweet sends a tweet using the Twitter API.
	Tweet(ctx context.Context, params TwitterTweetParams) (TwitterTweetResponse, error)
}

type twitterClient struct {
	base *baseClient
}

var _ TwitterClient = (*twitterClient)(nil)

// OAuthBegin begins an OAuth handshake.
func (c *twitterClient) OAuthBegin(ctx context.Context, request *http.Request) (*http.Response, error) {
	path, err := url.Parse("/twitter/oauth/begin")
	if err != nil {
		return nil, fmt.Errorf("unable to parse api url: %w", err)
	}
	request = request.WithContext(ctx)
	request.URL = path

	return c.base.Do(request)
}

// OAuthToken retrieves an OAuth token.
func (c *twitterClient) OAuthToken(ctx context.Context, request *http.Request) (*http.Response, error) {
	path, err := url.Parse("/twitter/oauth/token")
	if err != nil {
		return nil, fmt.Errorf("unable to parse api url: %w", err)
	}
	request = request.WithContext(ctx)
	request.URL = path

	return c.base.Do(request)
}

// Tweet sends a tweet using the Twitter API.
func (c *twitterClient) Tweet(ctx context.Context, params TwitterTweetParams) (TwitterTweetResponse, error) {
	return callAPI[TwitterTweetResponse](ctx, c.base, "POST", "/twitter/tweet", params)
}

type UrlShortenParams struct {
	URL string // the URL to shorten
}

type UrlURL struct {
	ID  string // short-form URL id
	URL string // complete URL, in long form
}

// UrlClient Provides you access to call public and authenticated APIs on url. The concrete implementation is urlClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type UrlClient interface {

	// Get retrieves the original URL for the id.
	Get(ctx context.Context, id string) (UrlURL, error)

	// Shorten shortens a URL.
	Shorten(ctx context.Context, params UrlShortenParams) (UrlURL, error)
}

type urlClient struct {
	base *baseClient
}

var _ UrlClient = (*urlClient)(nil)

// Get retrieves the original URL for the id.
func (c *urlClient) Get(ctx context.Context, id string) (UrlURL, error) {
	return callAPI[UrlURL](ctx, c.base, "GET", fmt.Sprintf("/url/%s", id), nil)
}

// Shorten shortens a URL.
func (c *urlClient) Shorten(ctx context.Context, params UrlShortenParams) (UrlURL, error) {
	return callAPI[UrlURL](ctx, c.base, "POST", "/url", params)
}

// HTTPDoer is an interface which can be used to swap out the default
// HTTP client (http.DefaultClient) with your own custom implementation.
// This can be used to inject middleware or mock responses during unit tests.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// baseClient holds all the information we need to make requests to an Encore application
type baseClient struct {
	tokenGenerator func(ctx context.Context) (string, error) // The function which will add the bearer token to the requests
	httpClient     HTTPDoer                                  // The HTTP client which will be used for all API requests
	baseURL        *url.URL                                  // The base URL which API requests will be made against
	userAgent      string                                    // What user agent we will use in the API requests
}

// Do sends the req to the Encore application adding the authorization token as required.
func (b *baseClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", b.userAgent)

	// If a authorization token generator is present, call it and add the returned token to the request
	if b.tokenGenerator != nil {
		if token, err := b.tokenGenerator(req.Context()); err != nil {
			return nil, fmt.Errorf("unable to create authorization token for api request: %w", err)
		} else if token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}
	}

	// Merge the base URL and the API URL
	req.URL = b.baseURL.ResolveReference(req.URL)
	req.Host = req.URL.Host

	// Finally, make the request via the configured HTTP Client
	return b.httpClient.Do(req)
}

// callAPI is used by each generated API method to actually make request and decode the responses
func callAPI[Response any](ctx context.Context, client *baseClient, method, path string, body any) (Response, error) {
	var response Response

	// Encode the API body
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return response, fmt.Errorf("unable to marshal api request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, path, bodyReader)
	if err != nil {
		return response, fmt.Errorf("unable to create api request: %w", err)
	}

	// Make the request via the base client
	rawResponse, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("api request failed: %w", err)
	}
	defer func() {
		_ = rawResponse.Body.Close()
	}()

	// Decode the response
	if err := json.NewDecoder(rawResponse.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("api request failed: %w", err)
	}
	return response, nil
}
