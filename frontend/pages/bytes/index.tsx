import { bytes } from '../../client/client';
import { DefaultClient } from '../../client/default';
import BytesList from '../../components/BytesList';
import { SEO } from '../../components/SEO';
import { InferGetStaticPropsType } from 'next'
import Page from '../../components/Page'

import {  GetStaticProps } from 'next'
function BytesIndex({bytes, page}: InferGetStaticPropsType<typeof getStaticProps>) {


  return (
    <div>
      <SEO
        title="Bytes"
        description="Quick dopamine hits"
      />
<Page page={page} />


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

  const res = await DefaultClient.bytes.List({Offset:0, Limit: 20})
  const bytes = res.Bytes
        const pageRes = await DefaultClient.blog.GetPage("bytes")
    const page = pageRes
  return {
    props: {
      bytes,
      page,
    },
    // Next.js will attempt to re-generate the page:
    // - When a request comes in
    // - At most once every 10 seconds
    revalidate: 60, // In seconds
  }
}
export default BytesIndex
