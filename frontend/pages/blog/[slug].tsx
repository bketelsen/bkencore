import Link from 'next/link'
import { DefaultClient } from '../../client/default'
import { timeToRead } from '../../components/BlogPostList'
import { SEO } from '../../components/SEO'
import { InferGetStaticPropsType } from 'next'
import { useMemo } from 'react'
import Image from '@/components/Image'
import CustomLink from '@/components/Link'
import Pre from '@/components/Pre'
import { GetStaticProps } from 'next'
import { getMdx } from '@/lib/mdx'
import { getMDXComponent, MDXContentProps } from 'mdx-bundler/client'
import { DateTime } from 'luxon'

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
  const created = DateTime.fromISO(post.created_at)

  return (
    <div className="max-w-3xl mx-auto">
      <SEO title={post?.title} description={post?.summary} />

      {!post ? (
        'Loading...'
      ) : (
        <>
          <div className="text-secondary">
            <Link href="/blog">
              <a className="hover:underline hover:decoration-neutral-300 font-sm">Blog</a>
            </Link>{' '}
            /
          </div>
          <h1 className="text-4xl font-extrabold text-primary">{post.title}</h1>
          <div className="mt-3 mb-3 text-base text-secondary">
            {' '}
            <time dateTime={post.created_at}>{created.toFormat('d LLL yyyy')}</time> -{' '}
            {timeToRead(post.body)}
          </div>
          {post.featured_image && (
            <div>
              <Image
                className="object-fill w-full h-auto max-w-3xl mt-6 mb-6 rounded-md"
                src={post.featured_image}
                height={'500'}
                width={'800'}
                alt={post.title}
              />
            </div>
          )}
        </>
      )}
      <div className="max-w-3xl mx-auto mt-6 prose ">
        <Component
          components={{
            Image,
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
  const posts = await DefaultClient.blog.GetBlogPosts({ limit: 10, offset: 0 , category:'', tag:''})
  const slugs = posts.blog_posts.map((post) => ({ params: { slug: post.slug } }))
  return {
    paths: slugs,
    fallback: 'blocking', // See the "fallback" section below
  }
}

export default BlogPost
