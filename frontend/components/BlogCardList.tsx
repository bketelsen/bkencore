import { DateTime } from 'luxon'
import Link from 'next/link'
import { FC } from 'react'
import { blog } from '../client/client'
import Image from './Image'

const BlogCardList: FC<{ posts: blog.BlogPost[] }> = ({ posts }) => (
  <>
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {posts.map((post) => {
        const created = DateTime.fromISO(post.CreatedAt)
        return (
          <div key={post.Slug} className="p-4">
            <div className="shadow-xl hover:shadow-2xl card bg-base-100">
              {post.FeaturedImage && (
                <Image
                  className="rounded-xl"
                  alt={post.Title}
                  src={post.FeaturedImage}
                  height={225}
                  width={400}
                />
              )}
              <div className="card-body">
                <div className="badge badge-lg">Article</div>

                <h2 className="card-title">{post.Title}</h2>

                <p>{post.Summary}</p>
                <div className="justify-center py-4 card-actions">
                  <Link href={'/blog/' + post.Slug}>
                    <a role="button" className="btn btn-secondary">
                      Read More
                    </a>
                  </Link>
                </div>
                <div className="flex flex-col justify-between text-sm">
                  <p className="text-gray-400 dark:text-gray-300">
                    <time dateTime={post.CreatedAt}>{created.toFormat('d LLL yyyy')}</time> -
                    Reading Time {timeToRead(post.Body)}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  </>
)

export function timeToRead(str: string): string {
  const wpm = 225 // average adult reading speed (words per minute)
  const words = str.trim().split(/\s+/).length
  const timeToRead = Math.ceil(words / wpm)
  return `${timeToRead} min${timeToRead !== 1 ? 's' : ''}`
}

export default BlogCardList
