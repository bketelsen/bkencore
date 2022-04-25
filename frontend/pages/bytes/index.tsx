import type { NextPage } from 'next';
import { useEffect, useState } from 'react';
import { bytes } from '../../client/client';
import { DefaultClient } from '../../client/default';
import BytesList from '../../components/BytesList';
import { SEO } from '../../components/SEO';

const BytesIndex: NextPage = () => {
  const [bytes, setBytes] = useState<bytes.Byte[]>()
  
  useEffect(() => {
    const fetch = async() => {
      const p = await DefaultClient.bytes.List({Limit: 100, Offset: 0})
      setBytes(p.Bytes ?? [])
    }
    fetch()
  }, [])
  
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

export default BytesIndex
