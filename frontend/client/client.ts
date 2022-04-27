export default class Client {
    blog: blog.ServiceClient
    bytes: bytes.ServiceClient
    email: email.ServiceClient
    twitter: twitter.ServiceClient
    url: url.ServiceClient

    constructor(environment: string = "staging", token?: string) {
        const base = new BaseClient(environment, token)
        this.blog = new blog.ServiceClient(base)
        this.bytes = new bytes.ServiceClient(base)
        this.email = new email.ServiceClient(base)
        this.twitter = new twitter.ServiceClient(base)
        this.url = new url.ServiceClient(base)
    }
}

export namespace blog {
    export interface BlogPost {
        Slug: string
        CreatedAt: string
        Published: boolean
        ModifiedAt: string
        Title: string
        Summary: string
        Body: string
        BodyRendered: string
        FeaturedImage: string
    }

    export interface CreateBlogPostParams {
        Slug: string
        Published: boolean
        Title: string
        Summary: string
        Body: string
        FeaturedImage: string
    }

    export interface CreatePageParams {
        Published: boolean
        Title: string
        Subtitle: string
        HeroText: string
        Summary: string
        Body: string
        /**
         * empty string means no image
         */
        FeaturedImage: string
    }

    export interface GetBlogPostsParams {
        Limit: number
        Offset: number
    }

    export interface GetBlogPostsResponse {
        Count: number
        BlogPosts: BlogPost[]
    }

    export interface Page {
        Slug: string
        CreatedAt: string
        ModifiedAt: string
        Published: boolean
        Title: string
        Subtitle: string
        HeroText: string
        Summary: string
        Body: string
        BodyRendered: string
        /**
         * emty string means no image
         */
        FeaturedImage: string
    }

    export interface PromoteParams {
        /**
         * Schedule decides how the promotion should be scheduled.
         * Valid values are "auto" for scheduling it at a suitable time
         * based on the current posting schedule, and "now" to schedule it immediately.
         */
        Schedule: ScheduleType
    }

    export type ScheduleType = string

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * CreateBlogPost creates a new blog post.
         */
        public CreateBlogPost(params: CreateBlogPostParams): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog.CreateBlogPost`, params)
        }

        /**
         * CreatePage creates a new page, or updates it if it already exists.
         */
        public CreatePage(slug: string, params: CreatePageParams): Promise<void> {
            return this.baseClient.doVoid("PUT", `/page/${slug}`, params)
        }

        /**
         * GetBlogPost retrieves a blog post by slug.
         */
        public GetBlogPost(slug: string): Promise<BlogPost> {
            return this.baseClient.do<BlogPost>("GET", `/blog/${slug}`)
        }

        /**
         * GetBlogPosts retrieves a list of blog posts with
         * optional limit and offset.
         */
        public GetBlogPosts(params: GetBlogPostsParams): Promise<GetBlogPostsResponse> {
            const query: any[] = [
                "limit", params.Limit,
                "offset", params.Offset,
            ]
            return this.baseClient.do<GetBlogPostsResponse>("GET", `/blog?${encodeQuery(query)}`)
        }

        /**
         * GetPage retrieves a page by slug.
         */
        public GetPage(slug: string): Promise<Page> {
            return this.baseClient.do<Page>("GET", `/page/${slug}`)
        }

        /**
         * Promote schedules the promotion a blog post.
         */
        public Promote(slug: string, params: PromoteParams): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog/${slug}/promote`, params)
        }
    }
}

export namespace bytes {
    export interface Byte {
        ID: number
        Title: string
        Summary: string
        URL: string
        Created: string
    }

    export interface ListParams {
        Limit: number
        Offset: number
    }

    export interface ListResponse {
        Bytes: Byte[]
    }

    export interface PublishParams {
        Title: string
        Summary: string
        URL: string
    }

    export interface PublishResponse {
        ID: number
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * List lists published bytes.
         */
        public List(params: ListParams): Promise<ListResponse> {
            const query: any[] = [
                "limit", params.Limit,
                "offset", params.Offset,
            ]
            return this.baseClient.do<ListResponse>("GET", `/bytes?${encodeQuery(query)}`)
        }

        /**
         * Publish publishes a byte.
         */
        public Publish(params: PublishParams): Promise<PublishResponse> {
            return this.baseClient.do<PublishResponse>("POST", `/bytes`, params)
        }
    }
}

export namespace email {
    export interface SubscribeParams {
        Email: string
    }

    export interface UnsubscribeParams {
        /**
         * Token is the unsubscribe token in to the email.
         */
        Token: string
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * Subscribe subscribes to the email newsletter for a given email.
         */
        public Subscribe(params: SubscribeParams): Promise<void> {
            return this.baseClient.doVoid("POST", `/email/subscribe`, params)
        }

        /**
         * Unsubscribe unsubscribes the user from the email list.
         */
        public Unsubscribe(params: UnsubscribeParams): Promise<void> {
            return this.baseClient.doVoid("POST", `/email/unsubscribe`, params)
        }
    }
}

export namespace twitter {
    export interface TweetParams {
        /**
         * Text is the text to tweet.
         */
        Text: string
    }

    export interface TweetResponse {
        /**
         * ID is the tweet id.
         */
        ID: string
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * OAuthBegin begins an OAuth handshake.
         */
        public OAuthBegin(): Promise<void> {
            return this.baseClient.doVoid("GET", `/twitter/oauth/begin`)
        }

        /**
         * OAuthToken retrieves an OAuth token.
         */
        public OAuthToken(): Promise<void> {
            return this.baseClient.doVoid("GET", `/twitter/oauth/token`)
        }

        /**
         * Tweet writes a mock tweet to the database.
         */
        public Tweet(params: TweetParams): Promise<TweetResponse> {
            return this.baseClient.do<TweetResponse>("POST", `/twitter/tweet`, params)
        }

        /**
         * Tweet sends a tweet using the Twitter API.
         */
        public TweetForReal(params: TweetParams): Promise<TweetResponse> {
            return this.baseClient.do<TweetResponse>("POST", `/twitter/tweet/for-real`, params)
        }
    }
}

export namespace url {
    export interface GetListResponse {
        Count: number
        URLS: URL[]
    }

    export interface ShortenParams {
        /**
         * the URL to shorten
         */
        URL: string
    }

    export interface URL {
        /**
         * short-form URL id
         */
        ID: string

        /**
         * original URL, in long form
         */
        URL: string

        /**
         * short URL
         */
        ShortURL: string
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * Get retrieves the original URL for the id.
         */
        public Get(id: string): Promise<URL> {
            return this.baseClient.do<URL>("GET", `/url/${id}`)
        }

        /**
         * List retrieves all shortened URLs
         */
        public List(): Promise<GetListResponse> {
            return this.baseClient.do<GetListResponse>("GET", `/url`)
        }

        /**
         * Shorten shortens a URL.
         */
        public Shorten(params: ShortenParams): Promise<URL> {
            return this.baseClient.do<URL>("POST", `/url`, params)
        }
    }
}

class BaseClient {
    baseURL: string
    headers: {[key: string]: string}

    constructor(environment: string, token?: string) {
        this.headers = {"Content-Type": "application/json"}
        if (token !== undefined) {
            this.headers["Authorization"] = "Bearer " + token
        }

        switch (environment) {
        case "local":
            this.baseURL = "http://localhost:4000"
            break
        case "staging":
            this.baseURL = "https://api.brian.dev"
            break
        default:
            this.baseURL = `https://devweek-k65i.encoreapi.com/${environment}`
        }
    }

    public async do<T>(method: string, path: string, req?: any): Promise<T> {
        let response = await fetch(this.baseURL + path, {
            method: method,
            headers: this.headers,
            body: req !== undefined ? JSON.stringify(req) : undefined
        })
        if (!response.ok) {
            let body = await response.text()
            throw new Error("request failed: " + body)
        }
        return <T>(await response.json())
    }

    public async doVoid(method: string, path: string, req?: any): Promise<void> {
        let response = await fetch(this.baseURL + path, {
            method: method,
            headers: this.headers,
            body: req !== undefined ? JSON.stringify(req) : undefined
        })
        if (!response.ok) {
            let body = await response.text()
            throw new Error("request failed: " + body)
        }
        await response.text()
    }
}

function encodeQuery(parts: any[]): string {
    const pairs = []
    for (let i = 0; i < parts.length; i += 2) {
        const key = parts[i]
        let val = parts[i+1]
        if (!Array.isArray(val)) {
            val = [val]
        }
        for (const v of val) {
            pairs.push(`${key}=${encodeURIComponent(v)}`)
        }
    }
    return pairs.join("&")
}
