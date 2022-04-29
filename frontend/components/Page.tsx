import { FC } from "react"
import { blog } from "../client/client";

const Page: FC<{page: blog.Page}> = ({page}) => {
  return (
    <div>
      <p className="text-base font-medium tracking-tight text-center lg:text-lg text-base-content">{page.Subtitle}</p>
      <h1 className="text-2xl font-extrabold tracking-tight text-center text-base-content md:text-3xl lg:text-4xl">{page.Title}</h1>
      <p className="my-6 text-xl text-center text-base-content ">{page.HeroText}</p>
    </div>
  )
}

export default Page
