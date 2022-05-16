import { DateTime } from 'luxon'
import Link from 'next/link'
import { FC } from 'react'
import { blog } from '../client/client'
import Image from './Image'

const BlogCardList: FC<{ posts: blog.BlogPostFull[] }> = ({ posts }) => (
  <>
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-2">
      {posts.map((post) => {
        const created = DateTime.fromISO(post.created_at)
        return (
          <div key={post.slug} className="p-4">
<div className="transition duration-300 ease-in-out shadow-xl card card-compact bg-base-300 hover:scale-105">
              {post.feature_image && (
                <Image
                  className="object-fill w-full h-auto max-w-3xl mx-auto mt-6 mb-6 rounded-md"
                  alt={post.title}
                  src={post.feature_image}
                  height={225}
                  width={400}

                />
              )}
              <div className="card-body">
                <div className="badge badge-lg badge-primary">{post.primary_tag && post.primary_tag.slug_name}</div>

                <h2 className="card-title">{post.title}</h2>

                <p>{post.excerpt}</p>
                <div className="justify-center py-4 card-actions">
                  <Link href={'/blog/' + post.slug}>
                    <a role="button" className="btn btn-primary">
                      Read More
                    </a>
                  </Link>
                </div>
                <div className="flex flex-col justify-between text-sm">

                  <p className="text-base-content">
                    <time dateTime={post.created_at}>{created.toFormat('d LLL yyyy')}</time> -
                    Reading Time {timeToRead(post.html)}
                    <br />
                    {post.tags && post.tags.map((tag) => {
                    return (
                      <span key={tag.slug_name} className="mt-2 mr-2 badge badge-accent">{'#' + tag.slug_name.toUpperCase()}</span>

                    )
                  })}
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
