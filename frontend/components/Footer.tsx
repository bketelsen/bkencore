/* This example requires Tailwind CSS v2.0+ */
import { FC } from 'react'
import NewsletterSignup from './NewsletterSignup'
import { social } from './social'

const Footer: FC = (props) => {
  return (
    <footer>
      <NewsletterSignup />

      <div className="pt-4 pb-12 px-4 sm:px-6 md:flex md:items-center md:justify-between">
        <div className="flex justify-center space-x-6 md:order-2">
          {social.map((item) => (
            <a key={item.name} rel="nofollow" href={item.href} className="text-gray-400 hover:text-gray-500">
              <span className="sr-only">{item.name}</span>
              <item.icon className="h-6 w-6" aria-hidden="true" />
            </a>
          ))}
        </div>
        <div className="mt-8 md:mt-0 md:order-1">
          <p className="text-center text-base text-gray-400">
            &copy; 2022 Brian Ketelsen 
            →
            <a href="https://twitter.com/bketelsen" className="ml-1 hover-underline" target="_blank" rel="noopener noreferrer">@bketelsen</a>
          </p>
        </div>
      </div>
    </footer>
  )
}

export default Footer
