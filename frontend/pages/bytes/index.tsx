import type { NextPage } from 'next'
import Head from 'next/head'
import useSWR from 'swr'

interface Byte {
  ID: number;
  Title: string;
  Summary: string;
  URL: string;
}

const BytesIndex: NextPage = () => {
  const {data: bytes, error} = useSWR<{Bytes: Byte[]}>("/bytes", {refreshInterval: 1000})
  
  return (
    <div>
      <Head>
        <title>Bytes | Brian Ketelsen</title>
      </Head>

      <div className="relative bg-gray-50 pt-16 pb-20 px-4 sm:px-6 lg:pt-24 lg:pb-28 lg:px-8">
        <div className="absolute inset-0">
          <div className="bg-white h-1/3 sm:h-2/3" />
        </div>
        <div className="relative max-w-7xl mx-auto">
          <div className="text-center">
            <h2 className="text-3xl tracking-tight font-extrabold text-gray-900 sm:text-4xl">Bytes</h2>
            <p className="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
              Lorem ipsum dolor sit amet consectetur, adipisicing elit. Ipsa libero labore natus atque, ducimus sed.
            </p>
          </div>
          <div className="mt-12 max-w-lg mx-auto grid gap-5 lg:grid-cols-3 lg:max-w-none">
            {error ? (
                <div>Failed to load bytes: {error}</div>
              ) : bytes === undefined ? (
                <div>Loading...</div>
              ) : <>
                {(bytes.Bytes ?? []).map((b) => (
                  <div key={b.ID} className="flex flex-col rounded-lg shadow-lg overflow-hidden">
                    <div className="flex-shrink-0">
                      {/* <img className="h-48 w-full object-cover" src={post.imageUrl} alt="" /> */}
                    </div>
                    <div className="flex-1 bg-white p-6 flex flex-col justify-between">
                      <div className="flex-1">
                        <p className="text-sm font-medium text-indigo-600">
                          {/* <a href={post.category.href} className="hover:underline">
                          {post.category.name}
                        </a> */}
                        </p>
                        <a href={b.URL} rel="nofollow" className="block mt-2">
                          <p className="text-xl font-semibold text-gray-900">{b.Title}</p>
                          <p className="mt-3 text-base text-gray-500">{b.Summary}</p>
                        </a>
                      </div>
                      
                    </div>
                  </div>
                ))}
              </>}
            </div>
        </div>
      </div>

    </div>
  )
}

export default BytesIndex
