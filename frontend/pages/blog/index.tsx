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
function BlogIndex({posts, page}: InferGetStaticPropsType<typeof getStaticProps>) {


  return (
    <div>
      <SEO
        title="Blog"
        description="I wrote this"
      />

<Page page={page} />


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
export  const getStaticProps: GetStaticProps = async()=>{

  const res = await DefaultClient.blog.GetBlogPosts({limit: 100, offset:0, category:'',tag:''})
  const posts = res.blog_posts
      const pageRes = await DefaultClient.blog.GetPage("blog")
    const page = pageRes
  return {
    props: {
      posts,
      page,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export default BlogIndex
