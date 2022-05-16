import { SEO } from '../components/SEO'
import { DefaultClient } from '../client/default'
import Page from '../components/Page'
import Image from '../components/Image'

import { InferGetStaticPropsType } from 'next'
import { GetStaticProps } from 'next'

function About({ page }: InferGetStaticPropsType<typeof getStaticProps>) {
  return (
    <div className="max-w-3xl mx-auto ">
      <SEO title="About Me" description="More than you probably need to know" />

      <Page title='About Me' hero_text="Here's some relevant information about me and my background." subtitle='Oversharing as a Service'/>

      <Image
        className="w-full h-auto max-w-3xl mt-6 mb-6 rounded-md"
        src="/static/images/brian.jpg"
        width={800}
        height={534}
        alt={'Brian Ketelsen'}
      />
      <div className="max-w-3xl mx-auto mt-6 prose prose-xl">
      <div dangerouslySetInnerHTML={{ __html: page.html }} />

      </div>
    </div>
  )
}
export const getStaticProps: GetStaticProps = async () => {
  const pageRes = await DefaultClient.blog.GetPage('about')
  const page = pageRes
  return {
    props: {
      page,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export default About
