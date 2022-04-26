import { bytes } from '../../client/client';
import { DefaultClient } from '../../client/default';
import BytesList from '../../components/BytesList';
import { SEO } from '../../components/SEO';
import { InferGetStaticPropsType } from 'next'

import {  GetStaticProps } from 'next'
function BytesIndex({bytes}: InferGetStaticPropsType<typeof getStaticProps>) {

  
  return (
    <div>
      <SEO
        title="Bytes"
        description="Quick dopamine hits"
      />

      <div>
        <p className="text-base lg:text-lg tracking-tight text-neutral-400">Quick dopamine hits</p>
        <h1 className="text-3xl font-extrabold tracking-tight text-neutral-900 md:text-4xl">Bytes</h1>
        <p className="mt-6 mb-9 text-xl text-neutral-500">
          I found it so you don&apos;t have to
        </p>
      </div>

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
