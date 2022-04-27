import { FC } from "react"
import { blog } from "../client/client";

const Page: FC<{page: blog.Page}> = ({page}) => {
  return (
    <div>
      <p className="text-base lg:text-lg font-medium tracking-tight text-neutral-400 text-center">{page.Subtitle}</p>
      <h1 className="text-2xl font-extrabold tracking-tight text-neutral-900 md:text-3xl lg:text-4xl text-center">{page.Title}</h1>
      <p className="my-6 text-xl text-neutral-500 text-center">{page.HeroText}</p>
    </div>
  )
}

export default Page
