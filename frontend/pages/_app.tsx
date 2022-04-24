import '../styles/globals.css'
import type { AppProps } from 'next/app'
import Layout from '../components/Layout'
import { SWRConfig } from 'swr'

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <SWRConfig
      value={{
        fetcher: (resource, init) => fetch(`http://localhost:4000${resource}`, init).then(res => res.json())
      }}
    >
      <Layout>
        <Component {...pageProps} />
      </Layout>

    </SWRConfig>
  )
} 

export default MyApp
