import type { NextPage } from 'next';
import Head from 'next/head';
import useSWR from 'swr';
import { bytes } from '../../client/client';
import BytesList from '../../components/BytesList';

const BytesIndex: NextPage = () => {
  const {data: bytes, error} = useSWR<{Bytes: bytes.Byte[]}>("/bytes", {refreshInterval: 1000})
  
  return (
    <div>
      <Head>
        <title>Bytes | Brian Ketelsen</title>
      </Head>

      <div>
        <p className="text-base lg:text-lg tracking-tight text-neutral-400">Quick dopamine hit</p>
        <h1 className="text-3xl font-extrabold tracking-tight text-neutral-900 md:text-4xl">Bytes</h1>
        <p className="mt-6 mb-9 text-xl text-neutral-500">
          I found it so you don&apos;t have to
        </p>
      </div>

      <section>
        {!bytes ? (
          <div className="text-neutral-400">Loading...</div>
        ) : (
          <BytesList bytes={bytes.Bytes} />
        )}
      </section>


    </div>
  )
}

export default BytesIndex
