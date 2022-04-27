import { DateTime } from "luxon"
import Link from "next/link"
import { FC } from "react"
import { blog } from "../client/client"

const BlogPostList: FC<{ posts: blog.BlogPost[] }> = ({ posts }) => (
  <>
    {posts.map((post) => {
      const created = DateTime.fromISO(post.CreatedAt)
      const modified = DateTime.fromISO(post.ModifiedAt)
      return (
        <div key={post.Slug} className="pt-8">
          <Link href={"/blog/" + post.Slug}>
            <a className="block hover-underline text-xl font-semibold text-neutral-900">
              {post.Title}
            </a>
          </Link>
          <p className="mt-1 text-sm text-neutral-500">
            <time dateTime={post.CreatedAt}>{created.toFormat("d LLL yyyy")}</time>
            <span className="px-2 text-primary-500">·</span>
            <span>{timeToRead(post.Body)}</span>

          </p>
          <p className="mt-2 text-base text-gray-600">{post.Summary}</p>
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
