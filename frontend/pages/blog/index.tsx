import type { NextPage } from 'next'
import Head from 'next/head'
import { useEffect, useState } from 'react'
import { blog } from '../../client/client'
import { DefaultClient } from '../../client/default'
import BlogPostList from '../../components/BlogPostList'
import { SEO } from '../../components/SEO'

const BlogIndex: NextPage = () => {
  const [posts, setPosts] = useState<blog.BlogPost[]>()
  
  useEffect(() => {
    const fetch = async() => {
      const p = await DefaultClient.blog.GetBlogPosts({Limit: 100, Offset: 0})
      setPosts(p.BlogPosts ?? [])
    }
    fetch()
  }, [])

  return (
    <div>
      <SEO
        title="Blog"
        description="I wrote this"
      />

      <div>
        <p className="text-base lg:text-lg tracking-tight text-neutral-400">I wrote this</p>
        <h1 className="text-3xl font-extrabold tracking-tight text-neutral-900 md:text-4xl">Blog</h1>
        <p className="mt-6 mb-9 text-xl text-neutral-500">
          It was just too long for a twitter thread.
        </p>
      </div>

      <section>
        {!posts ? (
          <div className="text-neutral-400">Loading...</div>
        ) : (
          <BlogPostList posts={posts} />
        )}
      </section>

    </div>
  )
}

export default BlogIndex
