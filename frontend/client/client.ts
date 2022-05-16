export default class Client {
    blog: blog.ServiceClient
    bytes: bytes.ServiceClient
    email: email.ServiceClient
    url: url.ServiceClient

    constructor(environment: string = "staging", token?: string) {
        const base = new BaseClient(environment, token)
        this.blog = new blog.ServiceClient(base)
        this.bytes = new bytes.ServiceClient(base)
        this.email = new email.ServiceClient(base)
        this.url = new url.ServiceClient(base)
    }
}

export namespace blog {
    export interface BlogPostFull {
        id: string
        uuid: string
        title: string
        slug: string
        html: string
        plaintext: string
        feature_image: string
        featured: boolean
        status: string
        visibility: string
        email_recipient_filter: string
        created_at: string
        updated_at: string
        published_at: string
        custom_excerpt: string
        canonical_url: string
        url: string
        excerpt: string
        reading_time: number
        og_image: string
        og_title: string
        og_description: string
        twitter_image: string
        twitter_title: string
        twitter_description: string
        meta_title: string
        meta_description: string
        feature_image_alt: string
        feature_image_caption: string
        primary_tag: Tag
        tags: Tag[]
    }

    export interface Category {
        category: string
        summary: string
    }

    export interface GetBlogPostsParams {
        limit: number
        offset: number
    }

    export interface GetBlogPostsResponse {
        count: number
        blog_posts: BlogPostFull[]
    }

    export interface GetCategoriesResponse {
        count: number
        categories: Category[]
    }

    export interface GetTagsResponse {
        count: number
        tags: Tag[]
    }

    export interface PageFull {
        id: string
        uuid: string
        title: string
        slug: string
        html: string
        plaintext: string
        feature_image: string
        featured: boolean
        status: string
        visibility: string
        email_recipient_filter: string
        created_at: string
        updated_at: string
        published_at: string
        custom_excerpt: string
        canonical_url: string
        url: string
        excerpt: string
        reading_time: number
        og_image: string
        og_title: string
        og_description: string
        twitter_image: string
        twitter_title: string
        twitter_description: string
        meta_title: string
        meta_description: string
        feature_image_alt: string
        feature_image_caption: string
        primary_tag: Tag
        tags: Tag[]
    }

    export interface Tag {
        slug_name: string
        slug: string
        slug_description: string
        feature_image: string
        visibility: string
        og_image: string
        og_title: string
        og_description: string
        twitter_image: string
        twitter_title: string
        twitter_description: string
        meta_title: string
        meta_description: string
        accent_color: string
        created_at: string
        updated_at: string
        slug_url: string
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * CreateTag creates a new blog post.
         */
        public CreateCategory(params: Category): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog.CreateCategory`, params)
        }

        /**
         * GetBlogPost retrieves a blog post by slug.
         */
        public GetBlogPost(slug: string): Promise<BlogPostFull> {
            return this.baseClient.do<BlogPostFull>("GET", `/blog/${slug}`)
        }

        /**
         * GetBlogPosts retrieves a list of blog posts with
         * optional limit and offset.
         */
        public GetBlogPosts(params: GetBlogPostsParams): Promise<GetBlogPostsResponse> {
            const query: any[] = [
                "limit", params.limit,
                "offset", params.offset,
            ]
            return this.baseClient.do<GetBlogPostsResponse>("GET", `/blog?${encodeQuery(query)}`)
        }

        /**
         * GetCategories retrieves a list of categories
         */
        public GetCategories(): Promise<GetCategoriesResponse> {
            return this.baseClient.do<GetCategoriesResponse>("GET", `/category`)
        }

        /**
         * GetCategory retrieves a category by slug.
         */
        public GetCategory(category: string): Promise<Category> {
            return this.baseClient.do<Category>("GET", `/category/${category}`)
        }

        /**
         * GetPage retrieves a page by slug.
         */
        public GetPage(slug: string): Promise<PageFull> {
            return this.baseClient.do<PageFull>("GET", `/page/${slug}`)
        }

        /**
         * GetTag retrieves a tag by slug.
         */
        public GetTag(slug: string): Promise<Tag> {
            return this.baseClient.do<Tag>("GET", `/tag/${slug}`)
        }

        /**
         * GetTags retrieves a list of tags
         */
        public GetTags(): Promise<GetTagsResponse> {
            return this.baseClient.do<GetTagsResponse>("GET", `/tag`)
        }

        /**
         * GetTagsBySlug retrieves a list of tags for a post
         */
        public GetTagsByPage(slug: string): Promise<GetTagsResponse> {
            return this.baseClient.do<GetTagsResponse>("GET", `/tagsbypage/${slug}`)
        }

        /**
         * GetTagsBySlug retrieves a list of tags for a post
         */
        public GetTagsByPost(slug: string): Promise<GetTagsResponse> {
            return this.baseClient.do<GetTagsResponse>("GET", `/tagsbypost/${slug}`)
        }

        /**
         * Post receives incoming post CRUD webhooks from ghost.
         */
        public PageHook(): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog.PageHook`)
        }

        /**
         * Post receives incoming post CRUD webhooks from ghost.
         */
        public PostHook(): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog.PostHook`)
        }

        /**
         * Post receives incoming post CRUD webhooks from ghost.
         */
        public TagHook(): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog.TagHook`)
        }
    }
}

export namespace bytes {
    export interface Byte {
        id: number
        title: string
        summary: string
        url: string
        created: string
    }

    export interface ListParams {
        limit: number
        offset: number
    }

    export interface ListResponse {
        bytes: Byte[]
    }

    export interface PromoteParams {
        /**
         * Schedule decides how the promotion should be scheduled.
         * Valid values are "auto" for scheduling it at a suitable time
         * based on the current posting schedule, and "now" to schedule it immediately.
         */
        Schedule: ScheduleType
    }

    export interface PublishParams {
        title: string
        summary: string
        url: string
    }

    export interface PublishResponse {
        id: number
    }

    export type ScheduleType = string

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * Get retrieves a byte.
         */
        public Get(id: number): Promise<Byte> {
            return this.baseClient.do<Byte>("GET", `/bytes/${id}`)
        }

        /**
         * List lists published bytes.
         */
        public List(params: ListParams): Promise<ListResponse> {
            const query: any[] = [
                "limit", params.limit,
                "offset", params.offset,
            ]
            return this.baseClient.do<ListResponse>("GET", `/bytes?${encodeQuery(query)}`)
        }

        /**
         * Promote schedules the promotion a byte.
         */
        public Promote(id: number, params: PromoteParams): Promise<void> {
            return this.baseClient.doVoid("POST", `/bytes/${id}/promote`, params)
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
        email: string
    }

    export interface UnsubscribeParams {
        /**
         * Token is the unsubscribe token in to the email.
         */
        token: string
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

export namespace url {
    export interface GetListResponse {
        count: number
        urls: URL[]
    }

    export interface ShortenParams {
        /**
         * the URL to shorten
         */
        url: string
    }

    export interface URL {
        /**
         * short-form URL id
         */
        id: string

        /**
         * original URL, in long form
         */
        url: string

        /**
         * short URL
         */
        short_url: string
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
