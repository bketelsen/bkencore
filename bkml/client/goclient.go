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
	Blog  BlogClient
	Bytes BytesClient
	Email EmailClient
	Url   UrlClient
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
		userAgent:  "devweek-k65i-Generated-Client (Encore/v1.0.2)",
	}

	// Apply any given options
	for _, option := range options {
		if err := option(base); err != nil {
			return nil, fmt.Errorf("unable to apply client option: %w", err)
		}
	}

	return &Client{
		Blog:  &blogClient{base},
		Bytes: &bytesClient{base},
		Email: &emailClient{base},
		Url:   &urlClient{base},
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
	Slug          string       `json:"slug"`
	CreatedAt     time.Time    `json:"created_at" qs:"created_at"`
	ModifiedAt    time.Time    `json:"modified_at" qs:"modified_at"`
	Published     bool         `json:"published"`
	Title         string       `json:"title"`
	Summary       string       `json:"summary"`
	Body          string       `json:"body"`
	BodyRendered  string       `json:"body_rendered" qs:"body_rendered"`
	FeaturedImage string       `json:"featured_image" qs:"featured_image"`
	Category      BlogCategory `json:"category"`
	Tags          []BlogTag    `json:"tags"`
}

type BlogCategory struct {
	Category string `json:"category"`
	Summary  string `json:"summary"`
}

type BlogCreateBlogPostParams struct {
	Slug          string   `json:"slug"`
	CreatedAt     string   `json:"created_at" qs:"created_at"`
	ModifiedAt    string   `json:"modified_at" qs:"modified_at"`
	Published     bool     `json:"published"`
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	Body          string   `json:"body"`
	FeaturedImage string   `json:"featured_image" qs:"featured_image"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
}

type BlogCreatePageParams struct {
	Published     bool   `json:"published"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	HeroText      string `json:"hero_text" qs:"hero_text"`
	Summary       string `json:"summary"`
	Body          string `json:"body"`
	FeaturedImage string `json:"featured_image" qs:"featured_image"` // empty string means no image
}

type BlogGetBlogPostsParams struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
}

type BlogGetBlogPostsResponse struct {
	Count     int            `json:"count"`
	BlogPosts []BlogBlogPost `json:"blog_posts" qs:"blog_posts"`
}

type BlogGetCategoriesResponse struct {
	Count      int            `json:"count"`
	Categories []BlogCategory `json:"categories"`
}

type BlogGetTagsResponse struct {
	Count int       `json:"count"`
	Tags  []BlogTag `json:"tags"`
}

type BlogPage struct {
	Slug          string    `json:"slug"`
	CreatedAt     time.Time `json:"created_at" qs:"created_at"`
	ModifiedAt    time.Time `json:"modified_at" qs:"modified_at"`
	Published     bool      `json:"published"`
	Title         string    `json:"title"`
	Subtitle      string    `json:"subtitle"`
	HeroText      string    `json:"hero_text" qs:"hero_text"`
	Summary       string    `json:"summary"`
	Body          string    `json:"body"`
	BodyRendered  string    `json:"body_rendered" qs:"body_rendered"`
	FeaturedImage string    `json:"featured_image" qs:"featured_image"`
}

type BlogPromoteParams struct {

	// Schedule decides how the promotion should be scheduled.
	// Valid values are "auto" for scheduling it at a suitable time
	// based on the current posting schedule, and "now" to schedule it immediately.
	Schedule BlogScheduleType
}

type BlogScheduleType = string

type BlogTag struct {
	Tag     string `json:"tag"`
	Summary string `json:"summary"`
}

// BlogClient Provides you access to call public and authenticated APIs on blog. The concrete implementation is blogClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type BlogClient interface {
	// CreateBlogPost creates a new blog post.
	CreateBlogPost(ctx context.Context, params BlogCreateBlogPostParams) error

	// CreateTag creates a new blog post.
	CreateCategory(ctx context.Context, params BlogCategory) error

	// CreatePage creates a new page, or updates it if it already exists.
	CreatePage(ctx context.Context, slug string, params BlogCreatePageParams) error

	// CreateTag creates a new blog post.
	CreateTag(ctx context.Context, params BlogTag) error

	// GetBlogPost retrieves a blog post by slug.
	GetBlogPost(ctx context.Context, slug string) (BlogBlogPost, error)

	// GetBlogPosts retrieves a list of blog posts with
	// optional limit and offset.
	GetBlogPosts(ctx context.Context, params BlogGetBlogPostsParams) (BlogGetBlogPostsResponse, error)

	// GetCategories retrieves a list of categories
	GetCategories(ctx context.Context) (BlogGetCategoriesResponse, error)

	// GetCategory retrieves a category by slug.
	GetCategory(ctx context.Context, category string) (BlogCategory, error)

	// GetPage retrieves a page by slug.
	GetPage(ctx context.Context, slug string) (BlogPage, error)

	// GetTag retrieves a tag by slug.
	GetTag(ctx context.Context, tag string) (BlogTag, error)

	// GetTags retrieves a list of tags
	GetTags(ctx context.Context) (BlogGetTagsResponse, error)

	// GetTagsBySlug retrieves a list of tags for a post
	GetTagsBySlug(ctx context.Context, slug string) (BlogGetTagsResponse, error)

	// Promote schedules the promotion a blog post.
	Promote(ctx context.Context, slug string, params BlogPromoteParams) error
}

type blogClient struct {
	base *baseClient
}

var _ BlogClient = (*blogClient)(nil)

// CreateBlogPost creates a new blog post.
func (c *blogClient) CreateBlogPost(ctx context.Context, params BlogCreateBlogPostParams) error {
	return callAPI(ctx, c.base, "POST", "/blog.CreateBlogPost", params, nil)
}

// CreateTag creates a new blog post.
func (c *blogClient) CreateCategory(ctx context.Context, params BlogCategory) error {
	return callAPI(ctx, c.base, "POST", "/blog.CreateCategory", params, nil)
}

// CreatePage creates a new page, or updates it if it already exists.
func (c *blogClient) CreatePage(ctx context.Context, slug string, params BlogCreatePageParams) error {
	return callAPI(ctx, c.base, "PUT", fmt.Sprintf("/page/%s", slug), params, nil)
}

// CreateTag creates a new blog post.
func (c *blogClient) CreateTag(ctx context.Context, params BlogTag) error {
	return callAPI(ctx, c.base, "POST", "/blog.CreateTag", params, nil)
}

// GetBlogPost retrieves a blog post by slug.
func (c *blogClient) GetBlogPost(ctx context.Context, slug string) (resp BlogBlogPost, err error) {
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/blog/%s", slug), nil, &resp)
	return resp, err
}

// GetBlogPosts retrieves a list of blog posts with
// optional limit and offset.
func (c *blogClient) GetBlogPosts(ctx context.Context, params BlogGetBlogPostsParams) (resp BlogGetBlogPostsResponse, err error) {
	queryString := url.Values{
		"category": []string{fmt.Sprint(params.Category)},
		"limit":    []string{fmt.Sprint(params.Limit)},
		"offset":   []string{fmt.Sprint(params.Offset)},
		"tag":      []string{fmt.Sprint(params.Tag)},
	}
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/blog?%s", queryString.Encode()), nil, &resp)
	return resp, err
}

// GetCategories retrieves a list of categories
func (c *blogClient) GetCategories(ctx context.Context) (resp BlogGetCategoriesResponse, err error) {
	err = callAPI(ctx, c.base, "GET", "/category", nil, &resp)
	return resp, err
}

// GetCategory retrieves a category by slug.
func (c *blogClient) GetCategory(ctx context.Context, category string) (resp BlogCategory, err error) {
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/category/%s", category), nil, &resp)
	return resp, err
}

// GetPage retrieves a page by slug.
func (c *blogClient) GetPage(ctx context.Context, slug string) (resp BlogPage, err error) {
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/page/%s", slug), nil, &resp)
	return resp, err
}

// GetTag retrieves a tag by slug.
func (c *blogClient) GetTag(ctx context.Context, tag string) (resp BlogTag, err error) {
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/tag/%s", tag), nil, &resp)
	return resp, err
}

// GetTags retrieves a list of tags
func (c *blogClient) GetTags(ctx context.Context) (resp BlogGetTagsResponse, err error) {
	err = callAPI(ctx, c.base, "GET", "/tag", nil, &resp)
	return resp, err
}

// GetTagsBySlug retrieves a list of tags for a post
func (c *blogClient) GetTagsBySlug(ctx context.Context, slug string) (resp BlogGetTagsResponse, err error) {
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/tagbyslug/%s", slug), nil, &resp)
	return resp, err
}

// Promote schedules the promotion a blog post.
func (c *blogClient) Promote(ctx context.Context, slug string, params BlogPromoteParams) error {
	return callAPI(ctx, c.base, "POST", fmt.Sprintf("/blog/%s/promote", slug), params, nil)
}

type BytesByte struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}

type BytesListParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type BytesListResponse struct {
	Bytes []BytesByte `json:"bytes"`
}

type BytesPromoteParams struct {

	// Schedule decides how the promotion should be scheduled.
	// Valid values are "auto" for scheduling it at a suitable time
	// based on the current posting schedule, and "now" to schedule it immediately.
	Schedule BytesScheduleType
}

type BytesPublishParams struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	URL     string `json:"url"`
}

type BytesPublishResponse struct {
	ID int64 `json:"id"`
}

type BytesScheduleType = string

// BytesClient Provides you access to call public and authenticated APIs on bytes. The concrete implementation is bytesClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type BytesClient interface {
	// Get retrieves a byte.
	Get(ctx context.Context, id int64) (BytesByte, error)

	// List lists published bytes.
	List(ctx context.Context, params BytesListParams) (BytesListResponse, error)

	// Promote schedules the promotion a byte.
	Promote(ctx context.Context, id int64, params BytesPromoteParams) error

	// Publish publishes a byte.
	Publish(ctx context.Context, params BytesPublishParams) (BytesPublishResponse, error)
}

type bytesClient struct {
	base *baseClient
}

var _ BytesClient = (*bytesClient)(nil)

// Get retrieves a byte.
func (c *bytesClient) Get(ctx context.Context, id int64) (resp BytesByte, err error) {
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/bytes/%d", id), nil, &resp)
	return resp, err
}

// List lists published bytes.
func (c *bytesClient) List(ctx context.Context, params BytesListParams) (resp BytesListResponse, err error) {
	queryString := url.Values{
		"limit":  []string{fmt.Sprint(params.Limit)},
		"offset": []string{fmt.Sprint(params.Offset)},
	}
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/bytes?%s", queryString.Encode()), nil, &resp)
	return resp, err
}

// Promote schedules the promotion a byte.
func (c *bytesClient) Promote(ctx context.Context, id int64, params BytesPromoteParams) error {
	return callAPI(ctx, c.base, "POST", fmt.Sprintf("/bytes/%d/promote", id), params, nil)
}

// Publish publishes a byte.
func (c *bytesClient) Publish(ctx context.Context, params BytesPublishParams) (resp BytesPublishResponse, err error) {
	err = callAPI(ctx, c.base, "POST", "/bytes", params, &resp)
	return resp, err
}

type EmailSubscribeParams struct {
	Email string `json:"email"`
}

type EmailUnsubscribeParams struct {
	Token string `json:"token"` // Token is the unsubscribe token in to the email.
}

// EmailClient Provides you access to call public and authenticated APIs on email. The concrete implementation is emailClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type EmailClient interface {
	// Subscribe subscribes to the email newsletter for a given email.
	Subscribe(ctx context.Context, params EmailSubscribeParams) error

	// Unsubscribe unsubscribes the user from the email list.
	Unsubscribe(ctx context.Context, params EmailUnsubscribeParams) error
}

type emailClient struct {
	base *baseClient
}

var _ EmailClient = (*emailClient)(nil)

// Subscribe subscribes to the email newsletter for a given email.
func (c *emailClient) Subscribe(ctx context.Context, params EmailSubscribeParams) error {
	return callAPI(ctx, c.base, "POST", "/email/subscribe", params, nil)
}

// Unsubscribe unsubscribes the user from the email list.
func (c *emailClient) Unsubscribe(ctx context.Context, params EmailUnsubscribeParams) error {
	return callAPI(ctx, c.base, "POST", "/email/unsubscribe", params, nil)
}

type UrlGetListResponse struct {
	Count int      `json:"count"`
	URLS  []UrlURL `json:"urls"`
}

type UrlShortenParams struct {
	URL string `json:"url"` // the URL to shorten
}

type UrlURL struct {
	ID       string `json:"id"`                       // short-form URL id
	URL      string `json:"url"`                      // original URL, in long form
	ShortURL string `json:"short_url" qs:"short_url"` // short URL
}

// UrlClient Provides you access to call public and authenticated APIs on url. The concrete implementation is urlClient.
// It is setup as an interface allowing you to use GoMock to create mock implementations during tests.
type UrlClient interface {
	// Get retrieves the original URL for the id.
	Get(ctx context.Context, id string) (UrlURL, error)

	// List retrieves all shortened URLs
	List(ctx context.Context) (UrlGetListResponse, error)

	// Shorten shortens a URL.
	Shorten(ctx context.Context, params UrlShortenParams) (UrlURL, error)
}

type urlClient struct {
	base *baseClient
}

var _ UrlClient = (*urlClient)(nil)

// Get retrieves the original URL for the id.
func (c *urlClient) Get(ctx context.Context, id string) (resp UrlURL, err error) {
	err = callAPI(ctx, c.base, "GET", fmt.Sprintf("/url/%s", id), nil, &resp)
	return resp, err
}

// List retrieves all shortened URLs
func (c *urlClient) List(ctx context.Context) (resp UrlGetListResponse, err error) {
	err = callAPI(ctx, c.base, "GET", "/url", nil, &resp)
	return resp, err
}

// Shorten shortens a URL.
func (c *urlClient) Shorten(ctx context.Context, params UrlShortenParams) (resp UrlURL, err error) {
	err = callAPI(ctx, c.base, "POST", "/url", params, &resp)
	return resp, err
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
func callAPI(ctx context.Context, client *baseClient, method, path string, body, resp any) error {
	// Encode the API body
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Make the request via the base client
	rawResponse, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		_ = rawResponse.Body.Close()
	}()
	if rawResponse.StatusCode >= 400 {
		return fmt.Errorf("got error response: %s", rawResponse.Status)
	}

	// Decode the response
	if resp != nil {
		if err := json.NewDecoder(rawResponse.Body).Decode(resp); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}
