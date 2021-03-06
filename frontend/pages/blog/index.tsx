import type { NextPage } from 'next'
import Head from 'next/head'
import { useEffect, useState } from 'react'
import { blog } from '../../client/client'
import { DefaultClient } from '../../client/default'
import BlogPostList from '../../components/BlogPostList'
import { SEO } from '../../components/SEO'
import { InferGetStaticPropsType } from 'next'
import Page from '../../components/Page'

import { GetStaticPaths, GetStaticProps } from 'next'
function BlogIndex({ posts }: InferGetStaticPropsType<typeof getStaticProps>) {


  return (
    <div>
      <SEO
        title="Blog"
        description="I wrote this"
      />
      <Page title='Articles' hero_text='It was just too long for a tweet' subtitle='I wrote this'/>

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
export const getStaticProps: GetStaticProps = async () => {

  const res = await DefaultClient.blog.GetBlogPosts({ limit: 100, offset: 0 })
  const posts = res.blog_posts
  return {
    props: {
      posts,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export default BlogIndex
