import type { NextPage } from 'next'
import Head from 'next/head'
import { useEffect, useState } from 'react'
import { blog, DefaultClient } from '../../client'

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
      <Head>
        <title>Blog | Brian Ketelsen</title>
      </Head>

      <div className="relative bg-gray-50 pt-16 pb-20 px-4 sm:px-6 lg:pt-24 lg:pb-28 lg:px-8">
        <div className="absolute inset-0">
          <div className="bg-white h-1/3 sm:h-2/3" />
        </div>
        <div className="relative max-w-7xl mx-auto">
          <div className="text-center">
            <h2 className="text-3xl tracking-tight font-extrabold text-gray-900 sm:text-4xl">From the blog</h2>
            <p className="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
              Lorem ipsum dolor sit amet consectetur, adipisicing elit. Ipsa libero labore natus atque, ducimus sed.
            </p>
          </div>
          <div className="mt-12 max-w-lg mx-auto grid gap-5 lg:grid-cols-3 lg:max-w-none">
            {posts === undefined ? "Loading..." : <>
              {posts.map((post) => (
                <div key={post.Slug} className="flex flex-col rounded-lg shadow-lg overflow-hidden">
                  <div className="flex-shrink-0">
                    {/* <img className="h-48 w-full object-cover" src={post.imageUrl} alt="" /> */}
                  </div>
                  <div className="flex-1 bg-white p-6 flex flex-col justify-between">
                    <div className="flex-1">
                      <p className="text-sm font-medium text-indigo-600">
                        {/* <a href={post.category.href} className="hover:underline">
                        {post.category.name}
                      </a> */}
                      </p>
                      <a href={'/blog/' + post.Slug} className="block mt-2">
                        <p className="text-xl font-semibold text-gray-900">{post.Title}</p>
                        <p className="mt-3 text-base text-gray-500">{post.Summary}</p>
                      </a>
                    </div>
                    <div className="mt-6 flex items-center">
                      <div className="flex-shrink-0">
                        {/* <a href={post.author.href}>
                        <span className="sr-only">{post.author.name}</span>
                        <img className="h-10 w-10 rounded-full" src={post.author.imageUrl} alt="" />
                      </a> */}
                      </div>
                      <div className="ml-3">
                        <p className="text-sm font-medium text-gray-900">
                          {/* <a href={post.author.href} className="hover:underline">
                          {post.author.name}
                        </a> */}
                        </p>
                        <div className="flex space-x-1 text-sm text-gray-500">
                          <time dateTime={post.CreatedAt}>{post.CreatedAt}</time>
                          <span aria-hidden="true">&middot;</span>
                          {/* <span>{post.readingTime} read</span> */}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </>}
          </div>
        </div>
      </div>

    </div>
  )
}

export default BlogIndex
