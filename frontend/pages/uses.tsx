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
function NowPage({ page }: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <div className="max-w-3xl mx-auto">
      <SEO title={page?.title} description={page?.excerpt} />

      {!page ? (
        'Loading...'
      ) : (
        <>
      <Page title='Uses' hero_text="These are the tools that power my workflows" subtitle="Behind the curtain"/>
        </>
      )}
      <div className="max-w-3xl mx-auto mt-6 prose ">
      <div dangerouslySetInnerHTML={{ __html: page.html }} />

      </div>
    </div>
  )
}

export const getStaticProps: GetStaticProps = async (context) => {
  const slug = 'uses'
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
