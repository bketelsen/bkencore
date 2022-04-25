import { NewspaperIcon, PhoneIcon, SupportIcon } from '@heroicons/react/outline'
import type { NextPage } from 'next'
import { useEffect, useState } from 'react'
import { blog } from '../client/client'
import { DefaultClient } from '../client/default'
import BlogPostList from '../components/BlogPostList'
import { SEO } from '../components/SEO'
import { social } from '../components/social'

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

const Home: NextPage = () => {
  const [posts, setPosts] = useState<blog.BlogPost[]>()
  useEffect(() => {
    const fetch = async() => {
      const p = await DefaultClient.blog.GetBlogPosts({Limit: 10, Offset: 0})
      setPosts(p.BlogPosts ?? [])
    }
    fetch()
  }, [])

  return (
    <div>
      <SEO
        title={undefined /* defaults to "Brian Ketelsen" */}
        description="Head in the clouds"
      />

      {/* Profile section */}
      <section className="flex flex-col justify-center max-w-[65ch] mx-auto pb-20 lg:pb-28">
        <p className="text-base lg:text-lg font-medium tracking-tight text-neutral-400 text-center">Head in the clouds</p>
        <h1 className="text-2xl font-extrabold tracking-tight text-neutral-900 md:text-3xl lg:text-4xl text-center">Brian Ketelsen</h1>
        <p className="my-6 text-xl text-neutral-500 text-center">
          Howdy! Thanks for stopping by. Inside you&apos;ll find articles, tutorials, technical reference material and maybe even a rant or two :)
        </p>
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

export default Home
