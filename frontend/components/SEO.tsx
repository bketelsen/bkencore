import * as React from "react";
import {DefaultSeo, NextSeo, NextSeoProps, DefaultSeoProps} from "next-seo";
import {useRouter} from "next/router";

export interface Props extends NextSeoProps {
  title?: string;
  description?: string;
  image?: string;
}

const siteURL = process.env.NEXT_PUBLIC_VERCEL_URL ? `https://${process.env.NEXT_PUBLIC_VERCEL_URL}` : "http://localhost:3000"

const title = "Brian Ketelsen"
const description = "Brian Ketelsen is a community leader and open source enthusiast focusing on Go, Rust, Web Development and WASM."
const image = `${siteURL}/assets/card.png`

const config: DefaultSeoProps = {
  title,
  description,
  openGraph: {
    type:      "website",
    url:       "https://brian.dev",
    site_name: title,
    images:    [{url: image}],
  },
  twitter:   {
    site:     "@bketelsen",
    handle:   "@bketelsen",
    cardType: "summary_large_image",
  },
};

export const SEO: React.FC<Props> = ({title, image, openGraph, ...props}) => {
  const router = useRouter()
  const path = router.asPath
  const url = `${siteURL}${path}`
  if (title) {
    title += " | Brian Ketelsen"
  }

  // merge in the defaults for open graph to what the page provided
  openGraph = openGraph ?? {}
  openGraph.url = url
  openGraph.images = (image) ? [{url: image}] : [];

  return (
    <>
      <DefaultSeo {...config} />

      <NextSeo
        canonical={url}
        openGraph={openGraph}
        title={title}
        {...props}
      />
    </>
  );
};
