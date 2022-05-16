import { FC } from 'react'
import { blog } from '../client/client'

const Page: FC<{ title: string, subtitle: string, hero_text:string}> = ({ title,subtitle,hero_text }) => {
  return (
    <div>
      <p className="text-base font-medium tracking-tight text-center lg:text-lg text-primary">
        {subtitle}
      </p>
      <h1 className="text-2xl font-extrabold tracking-tight text-center text-primary md:text-3xl lg:text-4xl">
        {title}
      </h1>
      <p className="my-6 text-xl text-center text-base-content ">{hero_text}</p>
    </div>
  )
}

export default Page
