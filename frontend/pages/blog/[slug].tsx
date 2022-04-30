import Link from 'next/link'
import { DefaultClient } from '../../client/default'
import { timeToRead } from '../../components/BlogPostList'
import { SEO } from '../../components/SEO'
import { InferGetStaticPropsType } from 'next'
import { useMemo } from 'react'
import Image from '@/components/Image'
import CustomLink from '@/components/Link'
import TOCInline from '@/components/TOCInline'
import Pre from '@/components/Pre'
import { GetStaticProps } from 'next'
import { getMdx } from '@/lib/mdx'
import { getMDXComponent, MDXContentProps } from 'mdx-bundler/client'

import { ParsedUrlQuery } from 'querystring'
export const MDXComponents = {
  // Image,
  // TOCInline,
}
interface IParams extends ParsedUrlQuery {
  slug: string
}
function BlogPost({ post, mdx }: InferGetStaticPropsType<typeof getStaticProps>) {
  const Component = useMemo(() => getMDXComponent(mdx.mdxSource), [mdx.mdxSource])
  return (
    <div className="max-w-3xl mx-auto">
      <SEO title={post?.Title} description={post?.Summary} />

      {!post ? (
        'Loading...'
      ) : (
        <>
          <div className="text-base-content">
            <Link href="/blog">
              <a className="hover:underline hover:decoration-neutral-300 font-sm">Blog</a>
            </Link>{' '}
            /
          </div>
          <h1 className="text-4xl font-extrabold text-base-content">{post.Title}</h1>
          <div className="mt-3 text-base text-secondary">{timeToRead(post.Body)}</div>
          {post.FeaturedImage && (
            <div>
              <Image
                className="object-fill w-full h-auto max-w-3xl mt-6 mb-6 rounded-md"
                src={post.FeaturedImage}
                height={'500'}
                width={'800'}
                alt={post.Title}
              />
            </div>
          )}
        </>
      )}
      <div className="max-w-3xl mx-auto mt-6 prose ">
        <Component
          components={{
            Image,
            TOCInline,
            a: CustomLink,
            pre: Pre,
          }}
        />
      </div>
    </div>
  )
}

export const getStaticProps: GetStaticProps = async (context) => {
  const { slug } = context.params as IParams
  const post = await DefaultClient.blog.GetBlogPost(slug as string)
  const mdx = await getMdx(post)
  return {
    props: {
      post,
      mdx,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export async function getStaticPaths() {
  const posts = await DefaultClient.blog.GetBlogPosts({ Limit: 10, Offset: 0 })
  const slugs = posts.BlogPosts.map((post) => ({ params: { slug: post.Slug } }))
  return {
    paths: slugs,
    fallback: 'blocking', // See the "fallback" section below
  }
}

export default BlogPost
