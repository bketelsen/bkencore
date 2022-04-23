export default class Client {
    blog: blog.ServiceClient
    email: email.ServiceClient
    hello: hello.ServiceClient
    twitter: twitter.ServiceClient
    url: url.ServiceClient

    constructor(environment: string = "prod", token?: string) {
        const base = new BaseClient(environment, token)
        this.blog = new blog.ServiceClient(base)
        this.email = new email.ServiceClient(base)
        this.hello = new hello.ServiceClient(base)
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
    }

    export interface CreateBlogPostParams {
        Slug: string
        Published: boolean
        Title: string
        Summary: string
        Body: string
    }

    export interface GetBlogPostsParams {
        Limit: number
        Offset: number
    }

    export interface GetBlogPostsResponse {
        Count: number
        BlogPosts: BlogPost[]
    }

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
    }
}

export namespace email {
    export interface CreateTemplateParams {
        /**
         * sender email
         */
        Sender: string

        /**
         * subject line to use
         */
        Subject: string

        /**
         * plaintext body
         */
        BodyText: string

        /**
         * html body
         */
        BodyHTML: string
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
         * CreateTemplate creates an email template.
         * If the template with that id already exists it is updated.
         */
        public CreateTemplate(id: string, params: CreateTemplateParams): Promise<void> {
            return this.baseClient.doVoid("PUT", `/email/templates/${id}`, params)
        }

        /**
         * Unsubscribe unsubscribes the user from the email list.
         */
        public Unsubscribe(params: UnsubscribeParams): Promise<void> {
            return this.baseClient.doVoid("POST", `/email/unsubscribe`, params)
        }
    }
}

export namespace hello {
    export interface Response {
        Message: string
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * This is a simple REST API that responds with a personalized greeting.
         * To call it, run in your terminal:
         * 
         *     curl http://localhost:4000/hello/World
         */
        public World(name: string): Promise<Response> {
            return this.baseClient.do<Response>("GET", `/hello/${name}`)
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
         * Tweet sends a tweet using the Twitter API.
         */
        public Tweet(params: TweetParams): Promise<TweetResponse> {
            return this.baseClient.do<TweetResponse>("POST", `/twitter/tweet`, params)
        }
    }
}

export namespace url {
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
         * complete URL, in long form
         */
        URL: string
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
        if (environment === "local") {
            this.baseURL = "http://localhost:4000"
        } else {
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
