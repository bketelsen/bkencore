import { DateTime } from "luxon"
import Link from "next/link"
import { FC } from "react"
import { blog } from "../client/client"

const BlogPostList: FC<{ posts: blog.BlogPost[] }> = ({ posts }) => (
  <>
    {posts.map((post) => {
      const created = DateTime.fromISO(post.created_at)
      const modified = DateTime.fromISO(post.modified_at)
      return (
        <div key={post.Slug} className="pt-8">
          <Link href={"/blog/" + post.Slug}>
            <a className="block text-xl font-semibold text-base-content hover-underline ">
              {post.Title}
            </a>
          </Link>
          <div className="my-2 badge badge-lg badge-secondary">{post.Category.Category.toUpperCase()}</div>

          <p className="mt-1 text-sm text-secondary">
            <time dateTime={post.created_at}>{created.toFormat("d LLL yyyy")}</time>
            <span className="px-2 text-primary">·</span>
            <span>{timeToRead(post.Body)}</span>
            {post.Tags && post.Tags.map((tag) => {
                    return (
                      <div key={tag.Tag} className="ml-4 badge badge-accent">{'#' + tag.Tag.toUpperCase()}</div>
                    )
                  })}
          </p>
          <p className="mt-2 text-base text-base-content ">{post.Summary}</p>
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
