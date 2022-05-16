import { DefaultClient } from '../client/default'
import { SEO } from '../components/SEO'
import { InferGetStaticPropsType } from 'next'
import { useMemo } from 'react'
import Image from '@/components/Image'
import CustomLink from '@/components/Link'
import Pre from '@/components/Pre'
import { GetStaticProps } from 'next'
import { getMdx } from '@/lib/mdx'
import { getMDXComponent, MDXContentProps } from 'mdx-bundler/client'
import Page from '../components/Page'

import { ParsedUrlQuery } from 'querystring'
export const MDXComponents = {
  // Image,
  // TOCInline,
}
interface IParams extends ParsedUrlQuery {
  slug: string
}
function NowPage({ page}: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <div className="max-w-3xl mx-auto">
      <SEO title={page?.Title} description={page?.Summary} />

      {!page ? (
        'Loading...'
      ) : (
        <>
      <Page title='Now' hero_text="Catch up if we haven't talked in a bit" subtitle="What I'm Doing Now"/>
        </>
      )}
      <div className="max-w-3xl mx-auto mt-6 prose ">
      <div className="max-w-3xl mx-auto mt-6 prose prose-xl">
      <div dangerouslySetInnerHTML={{ __html: page.html }} />

      </div>
      </div>
    </div>
  )
}

export const getStaticProps: GetStaticProps = async (context) => {
  const slug = 'now'
  const page = await DefaultClient.blog.GetPage(slug as string)
  return {
    props: {
      page,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 60 seconds
    revalidate: 60, // In seconds
  }
}

export default NowPage
