import type { NextPage } from 'next'
import Head from 'next/head'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import { blog, DefaultClient } from '../../client'

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
      <Head>
        <title>Create Next App</title>
      </Head>

      {!post ? "Loading..." : (
        <div className="relative py-16 bg-white overflow-hidden">
          <div className="hidden lg:block lg:absolute lg:inset-y-0 lg:h-full lg:w-full">
            <div className="relative h-full text-lg max-w-prose mx-auto" aria-hidden="true">
              <svg
                className="absolute top-12 left-full transform translate-x-32"
                width={404}
                height={384}
                fill="none"
                viewBox="0 0 404 384"
              >
                <defs>
                  <pattern
                    id="74b3fd99-0a6f-4271-bef2-e80eeafdf357"
                    x={0}
                    y={0}
                    width={20}
                    height={20}
                    patternUnits="userSpaceOnUse"
                  >
                    <rect x={0} y={0} width={4} height={4} className="text-gray-200" fill="currentColor" />
                  </pattern>
                </defs>
                <rect width={404} height={384} fill="url(#74b3fd99-0a6f-4271-bef2-e80eeafdf357)" />
              </svg>
              <svg
                className="absolute top-1/2 right-full transform -translate-y-1/2 -translate-x-32"
                width={404}
                height={384}
                fill="none"
                viewBox="0 0 404 384"
              >
                <defs>
                  <pattern
                    id="f210dbf6-a58d-4871-961e-36d5016a0f49"
                    x={0}
                    y={0}
                    width={20}
                    height={20}
                    patternUnits="userSpaceOnUse"
                  >
                    <rect x={0} y={0} width={4} height={4} className="text-gray-200" fill="currentColor" />
                  </pattern>
                </defs>
                <rect width={404} height={384} fill="url(#f210dbf6-a58d-4871-961e-36d5016a0f49)" />
              </svg>

            </div>
          </div>
          <div className="relative px-4 sm:px-6 lg:px-8">
            <div className="text-lg max-w-prose mx-auto">
              <h1>
                <span className="mt-2 block text-3xl text-center leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl">
                  {post.Title}
                </span>
              </h1>
              <p className="mt-8 text-xl text-gray-500 leading-8">
                {post.Summary}
              </p>
            </div>
            <div className="mt-6 prose prose-indigo prose-lg text-gray-500 mx-auto"
              dangerouslySetInnerHTML={{ __html: post.BodyRendered }} />
          </div>
        </div>
      )}
    </div>
  )
}

export default BlogPost
