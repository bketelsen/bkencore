import { DateTime } from "luxon"
import Link from "next/link"
import { FC } from "react"
import { blog } from "../client/client"

const BlogPostList: FC<{ posts: blog.BlogPostFull[]  }> = ({ posts }) => (
  <>
    {posts.map((post) => {
      const created = DateTime.fromISO(post.published_at)
      return (
        <div key={post.slug} className="pt-8">
          <Link href={"/blog/" + post.slug}>
            <a className="block text-xl font-semibold text-base-content hover-underline ">
              {post.title}
            </a>
          </Link>
          <div className="my-2 badge badge-lg badge-secondary">{post.primary_tag.slug_name.toUpperCase()}</div>
          <p className="mt-1 text-sm text-primary">
            <time dateTime={post.published_at}>{created.toFormat("d LLL yyyy")}</time>
            <span className="px-2 text-primary">·</span>
            <span>{timeToRead(post.html|| "")}</span>
            {post.tags && post.tags.map((tag) => {
                    return (
                      <div key={tag.slug_name} className="ml-4 badge badge-accent">{'#' + tag.slug_name.toUpperCase()}</div>
                    )
                  })}
          </p>
          <p className="mt-2 text-base text-base-content ">{post.excerpt}</p>
        </div>
      )
    })}
  </>
)

export function timeToRead(str: string): string {
  const wpm = 225 // average adult reading speed (words per minute)
  const words = str.trim().split(/\s+/).length;
  const timeToRead = Math.ceil(words / wpm)
  return `${timeToRead} min${timeToRead !== 1 ? "s" : ""}`
}

export default BlogPostList
