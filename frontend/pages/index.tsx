import { NewspaperIcon, PhoneIcon, SupportIcon } from '@heroicons/react/outline'
import { DefaultClient } from '../client/default'
import BlogPostList from '../components/BlogCardList'
import { SEO } from '../components/SEO'
import { social } from '../components/social'
import Page from '../components/Page'
import { InferGetStaticPropsType } from 'next'

import {  GetStaticProps } from 'next'
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

function Home({ posts }: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <div>
      <SEO title={undefined /* defaults to "Brian Ketelsen" */} description="Head in the clouds" />

      {/* Profile section */}
      <Page title='Brian Ketelsen' hero_text="Howdy! Thanks for stopping by. Inside you'll find articles, tutorials, technical reference material and maybe even a rant or two :)" subtitle='Head in the clouds'/>

      <div className="flex justify-center space-x-6">
        {social.map((item) => (
          <a
            key={item.name}
            rel="nofollow"
            href={item.href}
            className="text-primary hover:text-secondary"
          >
            <span className="sr-only">{item.name}</span>
            <item.icon className="w-6 h-6" aria-hidden="true" />
          </a>
        ))}
      </div>

      <h2 className="flex justify-center mt-10 text-xl font-bold tracking-tight title-font text-base-content sm:text-2xl">
        Recent Articles
      </h2>
      {!posts ? (
        <div className="text-neutral-content">Loading...</div>
      ) : (
        <BlogPostList posts={posts} />
      )}
    </div>
  )
}
export const getStaticProps: GetStaticProps = async () => {
  const res = await DefaultClient.blog.GetBlogPosts({ limit: 6, offset: 0  })
  const posts = res.blog_posts
  //const pageRes = await DefaultClient.blog.GetPage('index')
  //const page = pageRes
  return {
    props: {
      posts,
     // page,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export default Home
