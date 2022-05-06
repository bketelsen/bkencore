import { FC } from 'react'
import { blog } from '../client/client'

const Page: FC<{ page: blog.Page }> = ({ page }) => {
  return (
    <div>
      <p className="text-base font-medium tracking-tight text-center lg:text-lg text-secondary">
        {page.subtitle}
      </p>
      <h1 className="text-2xl font-extrabold tracking-tight text-center text-primary md:text-3xl lg:text-4xl">
        {page.title}
      </h1>
      <p className="my-6 text-xl text-center text-base-content ">{page.hero_text}</p>
    </div>
  )
}

export default Page
