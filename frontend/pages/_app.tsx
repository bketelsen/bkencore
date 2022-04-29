import '@/styles/globals.css' 
import '@/styles/prism.css'

import type { AppProps } from 'next/app'
import Layout from '../components/Layout'

var baseURL = "http://localhost:4000"
const env = process.env.NEXT_PUBLIC_ENCORE_ENV ?? "staging"
switch (env) {
  case "local":
      baseURL = "http://localhost:4000"
      break
  case "staging":
      baseURL = "https://api.brian.dev"
      break
  default:
      baseURL = `https://devweek-k65i.encoreapi.com/${env}`
  }
console.log(baseURL)
function MyApp({ Component, pageProps }: AppProps) {
  return (
      <Layout>
        <Component {...pageProps} />
      </Layout>

  )
} 

export default MyApp
