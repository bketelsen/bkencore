import type { NextPage } from 'next'
import Head from 'next/head'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import { blog } from '../../client/client'
import { DefaultClient } from '../../client/default'
import { timeToRead } from '../../components/BlogPostList'
import { SEO } from '../../components/SEO'
import { InferGetStaticPropsType } from 'next'

import { GetStaticPaths, GetStaticProps } from 'next'
import { ParsedUrlQuery } from 'querystring'

interface IParams extends ParsedUrlQuery {
  slug: string
}
function BlogPost({post}: InferGetStaticPropsType<typeof getStaticProps>) {
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

export  const getStaticProps: GetStaticProps = async(context)=>{
  const { slug } = context.params as IParams
  const post = await DefaultClient.blog.GetBlogPost(slug as string)

  return {
    props: {
      post,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export async function getStaticPaths() {
 const posts= await DefaultClient.blog.GetBlogPosts({Limit: 10, Offset:0})
const slugs = posts.BlogPosts.map((post) =>  ({ params: {slug: post.Slug}}))
  return {
    paths: slugs,
    fallback: "blocking" // See the "fallback" section below
  };
}

export default BlogPost
