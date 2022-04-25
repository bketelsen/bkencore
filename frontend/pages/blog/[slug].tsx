import type { NextPage } from 'next'
import Head from 'next/head'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import { blog } from '../../client/client'
import { DefaultClient } from '../../client/default'
import { timeToRead } from '../../components/BlogPostList'
import { SEO } from '../../components/SEO'

const BlogPost: NextPage = () => {
  const router = useRouter()
  const {slug} = router.query
  const [post, setPost] = useState<blog.BlogPost>()

  useEffect(() => {
    if (!slug) { return }

    const fetch = async() => {
      const p = await DefaultClient.blog.GetBlogPost(slug as string)
      setPost(p)
    }
    fetch()
  }, [slug])


  return (
    <div>
      <SEO
        title={post?.Title}
        description={post?.Summary}
      />

      {!post ? "Loading..." : <>
        <div className="text-neutral-500">
          <Link href="/blog"><a className="hover:underline hover:decoration-neutral-300 font-sm">Blog</a></Link> /
        </div>
        <h1 className="text-4xl font-extrabold text-neutral-900">{post.Title}</h1>
        <div className="mt-3 text-base text-neutral-500">{timeToRead(post.Body)}</div>
        {post.FeaturedImage && <img className="mt-6 mb-6 rounded-md w-full h-auto max-w-prose" src={post.FeaturedImage} />}
        <div className="mt-6 prose prose-indigo text-gray-500 max-w-prose"
          dangerouslySetInnerHTML={{ __html: post.BodyRendered }} />
      </>}
    </div>
  )
}

export default BlogPost
