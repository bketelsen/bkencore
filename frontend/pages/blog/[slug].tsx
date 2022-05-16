import Link from 'next/link'
import { DefaultClient } from '../../client/default'
import { timeToRead } from '../../components/BlogPostList'
import { SEO } from '../../components/SEO'
import { InferGetStaticPropsType } from 'next'
import Image from '@/components/Image'
import { GetStaticProps } from 'next'
import { DateTime } from 'luxon'

import { ParsedUrlQuery } from 'querystring'
export const MDXComponents = {
  // Image,
  // TOCInline,
}
interface IParams extends ParsedUrlQuery {
  slug: string
}
function BlogPost({ post }: InferGetStaticPropsType<typeof getStaticProps>) {
  const created = DateTime.fromISO(post.published_at)

  return (
    <div className="max-w-3xl mx-auto">
      <SEO title={post?.title} description={post?.summary} />

      {!post ? (
        'Loading...'
      ) : (
        <>
          {post.feature_image && (
            <div>
              <Image
                className="object-fill w-full h-auto max-w-3xl mx-auto mt-6 mb-6 rounded-md"
                src={post.feature_image}
                height={'500'}
                width={'800'}
                alt={post.title}
              />
            </div>
          )}
          <div className="text-primary">
            <Link href="/blog">
              <a className="hover:underline hover:decoration-neutral-300 font-sm">Blog</a>
            </Link>{' '}
            /
          </div>
          <h1 className="text-4xl font-extrabold text-primary">{post.title}</h1>
          <div className="mt-3 mb-3 text-base text-secondary">
            {' '}
            <time dateTime={post.published_at}>{created.toFormat('d LLL yyyy')}</time> -{' '}
            {timeToRead(post.html || '')}
          </div>

        </>
      )}
      <div className="max-w-3xl mx-auto mt-6 prose prose-xl ">

        <div dangerouslySetInnerHTML={{ __html: post.html }} />

      </div>
    </div>
  )
}

export const getStaticProps: GetStaticProps = async (context) => {
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
  const posts = await DefaultClient.blog.GetBlogPosts({ limit: 100, offset: 0 })
  const slugs = posts.blog_posts.map((post) => ({ params: { slug: post.slug } }))
  return {
    paths: slugs,
    fallback: 'blocking', // See the "fallback" section below
  }
}

export default BlogPost
