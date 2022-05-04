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
function LiamPage({ page, mdx }: InferGetStaticPropsType<typeof getStaticProps>) {
  const Component = useMemo(() => getMDXComponent(mdx.mdxSource), [mdx.mdxSource])
  return (
    <div className="max-w-3xl mx-auto">
      <SEO title={page?.Title} description={page?.Summary} />

      {!page ? (
        'Loading...'
      ) : (
        <>
          <Page page={page} />
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
  const slug = 'liam'
  const page = await DefaultClient.blog.GetPage(slug as string)
  const mdx = await getMdx(page)
  return {
    props: {
      page,
      mdx,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 60 seconds
    revalidate: 60, // In seconds
  }
}

export default LiamPage
