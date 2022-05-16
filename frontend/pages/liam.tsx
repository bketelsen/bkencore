import { DefaultClient } from '../client/default'
import { SEO } from '../components/SEO'
import { InferGetStaticPropsType } from 'next'
import { GetStaticProps } from 'next'
import Page from '../components/Page'

import { ParsedUrlQuery } from 'querystring'
export const MDXComponents = {
  // Image,
  // TOCInline,
}
interface IParams extends ParsedUrlQuery {
  slug: string
}
function LiamPage({ page}: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <div className="max-w-3xl mx-auto">
      <SEO title={page?.Title} description={page?.Summary} />

      {!page ? (
        'Loading...'
      ) : (
        <>
          <Page title='Liam' subtitle='Living with Muenke Syndrome' hero_text='The continuing adventures of Liam Walker' />
        </>
      )}
      <div className="max-w-3xl mx-auto mt-6 prose ">
      <div dangerouslySetInnerHTML={{ __html: page.html }} />

      </div>
    </div>
  )
}

export const getStaticProps: GetStaticProps = async (context) => {
  const slug = 'liam'
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

export default LiamPage
