import { bytes } from '../../client/client';
import { DefaultClient } from '../../client/default';
import BytesList from '../../components/BytesList';
import { SEO } from '../../components/SEO';
import { InferGetStaticPropsType } from 'next'
import Page from '../../components/Page'

import {  GetStaticProps } from 'next'
function BytesIndex({bytes}: InferGetStaticPropsType<typeof getStaticProps>) {


  return (
    <div>
      <SEO
        title="Bytes"
        description="Quick dopamine hits"
      />
      <Page title='Bytes' hero_text="I found it so you don't have to" subtitle='Quick Dopamine Hits'/>

      <section>
        {!bytes ? (
          <div className="text-neutral-400">Loading...</div>
        ) : (
          <BytesList bytes={bytes} />
        )}
      </section>


    </div>
  )
}
export  const getStaticProps: GetStaticProps = async()=>{

  const res = await DefaultClient.bytes.List({ offset: 0, limit: 20 })
   const bytes = res.bytes

  return {
    props: {
      bytes,

    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export default BytesIndex
