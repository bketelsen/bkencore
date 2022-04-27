import { NewspaperIcon, PhoneIcon, SupportIcon } from '@heroicons/react/outline'
import type { NextPage } from 'next'
import { useEffect, useState } from 'react'
import { blog } from '../client/client'
import { DefaultClient } from '../client/default'
import BlogPostList from '../components/BlogPostList'
import { SEO } from '../components/SEO'
import { social } from '../components/social'
import Page from '../components/Page'
import { InferGetStaticPropsType } from 'next'

import { GetStaticPaths, GetStaticProps } from 'next'
const links = [
  {
    name: 'Blog',
    href: '/blog',
    description: 'It was just too long for a twitter thread',
    icon: PhoneIcon,
  },
  {
    name: 'Bytes',
    href: '/bytes',
    description: 'Quick hit news and interesting articles',
    icon: SupportIcon,
  },
  {
    name: 'About',
    href: '/about',
    description: 'Everything you need to know',
    icon: NewspaperIcon,
  },
]

function Home({posts, page}: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <div>
      <SEO
        title={undefined /* defaults to "Brian Ketelsen" */}
        description="Head in the clouds"
      />

      {/* Profile section */}
      <section className="flex flex-col justify-center max-w-[65ch] mx-auto pb-20 lg:pb-28">
      <Page page={page} />
        <div className="flex justify-center space-x-6">
          {social.map((item) => (
            <a key={item.name} rel="nofollow" href={item.href} className="text-gray-400 hover:text-gray-500">
              <span className="sr-only">{item.name}</span>
              <item.icon className="h-6 w-6" aria-hidden="true" />
            </a>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-xl tracking-tight font-extrabold text-gray-900 sm:text-2xl">Recent blog posts</h2>
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

  const res = await DefaultClient.blog.GetBlogPosts({Limit: 5, Offset:0})
  const posts = res.BlogPosts
    const pageRes = await DefaultClient.blog.GetPage("index")
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
export default Home
