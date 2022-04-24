import { FC } from "react"
import { blog } from "../client/client";

const Page: FC<{page: blog.Page}> = ({page}) => {
  return (
    <div>
      <p className="text-base lg:text-lg tracking-tight text-neutral-400">Head in the clouds</p>
      <h1 className="text-3xl font-extrabold tracking-tight text-neutral-900 md:text-4xl">Brian Ketelsen</h1>
      <p className="mt-6 mb-9 text-xl text-neutral-500">
        Howdy! Thanks for stopping by. Inside you&apos;ll find articles, tutorials, technical reference material and maybe even a rant or two :)
      </p>
    </div>
  )
}

export default Page
