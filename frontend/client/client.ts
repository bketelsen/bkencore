export default class Client {
    blog: blog.ServiceClient
    bytes: bytes.ServiceClient
    email: email.ServiceClient
    url: url.ServiceClient

    constructor(environment: string = "prod", token?: string) {
        const base = new BaseClient(environment, token)
        this.blog = new blog.ServiceClient(base)
        this.bytes = new bytes.ServiceClient(base)
        this.email = new email.ServiceClient(base)
        this.url = new url.ServiceClient(base)
    }
}

export namespace blog {
    export interface BlogPost {
        slug: string
        created_at: string
        modified_at: string
        published: boolean
        title: string
        summary: string
        body: string
        body_rendered: string
        featured_image: string
        category: Category
        tags: Tag[]
    }

    export interface Category {
        category: string
        summary: string
    }

    export interface CreateBlogPostParams {
        slug: string
        created_at: string
        modified_at: string
        published: boolean
        title: string
        summary: string
        body: string
        featured_image: string
        category: string
        tags: string[]
    }

    export interface CreatePageParams {
        published: boolean
        title: string
        subtitle: string
        hero_text: string
        summary: string
        body: string
        /**
         * empty string means no image
         */
        featured_image: string
    }

    export interface GetBlogPostsParams {
        limit: number
        offset: number
        category: string
        tag: string
    }

    export interface GetBlogPostsResponse {
        count: number
        blog_posts: BlogPost[]
    }

    export interface GetCategoriesResponse {
        count: number
        categories: Category[]
    }

    export interface GetTagsResponse {
        count: number
        tags: Tag[]
    }

    export interface Page {
        slug: string
        created_at: string
        modified_at: string
        published: boolean
        title: string
        subtitle: string
        hero_text: string
        summary: string
        body: string
        body_rendered: string
        featured_image: string
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

    export interface Tag {
        tag: string
        summary: string
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
         * CreateTag creates a new blog post.
         */
        public CreateCategory(params: Category): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog.CreateCategory`, params)
        }

        /**
         * CreatePage creates a new page, or updates it if it already exists.
         */
        public CreatePage(slug: string, params: CreatePageParams): Promise<void> {
            return this.baseClient.doVoid("PUT", `/page/${slug}`, params)
        }

        /**
         * CreateTag creates a new blog post.
         */
        public CreateTag(params: Tag): Promise<void> {
            return this.baseClient.doVoid("POST", `/blog.CreateTag`, params)
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
                "limit", params.limit,
                "offset", params.offset,
                "category", params.category,
                "tag", params.tag,
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
        public GetPage(slug: string): Promise<Page> {
            return this.baseClient.do<Page>("GET", `/page/${slug}`)
        }

        /**
         * GetTag retrieves a tag by slug.
         */
        public GetTag(tag: string): Promise<Tag> {
            return this.baseClient.do<Tag>("GET", `/tag/${tag}`)
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
        public GetTagsBySlug(slug: string): Promise<GetTagsResponse> {
            return this.baseClient.do<GetTagsResponse>("GET", `/tagbyslug/${slug}`)
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
