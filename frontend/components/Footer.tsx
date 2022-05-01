/* This example requires Tailwind CSS v2.0+ */
import { FC } from 'react'
import NewsletterSignup from './NewsletterSignup'
import { social } from './social'
import Link from 'next/link'

const Footer: FC = (props) => {
  return (
    <>
      <NewsletterSignup />

      <footer className="p-10 mt-10 rounded lg:mt-20 footer footer-center bg-base-200 text-base-content">
        <div className="grid grid-flow-col gap-4">
          <Link href={'/about'}>
            <a className="link link-hover hover-underline">About</a>
          </Link>
          <Link href={'/uses'}>
            <a className="link link-hover hover-underline ">Uses</a>
          </Link>
          <Link href={'/now'}>
            <a className="link link-hover hover-underline">Now</a>
          </Link>
        </div>
        <div>
          <div className="grid grid-flow-col gap-4">
            {social.map((item) => (
              <a
                key={item.name}
                rel="nofollow"
                href={item.href}
                className="text-gray-400 hover:text-gray-500"
              >
                <span className="sr-only">{item.name}</span>
                <item.icon className="w-6 h-6" aria-hidden="true" />
              </a>
            ))}
          </div>
        </div>
        <div>
          <p>
            Copyright © 2022 - Brian Ketelsen →
            <a
              href="https://twitter.com/bketelsen"
              className="ml-1 hover-underline"
              target="_blank"
              rel="noopener noreferrer"
            >
              @bketelsen
            </a>
          </p>
        </div>
      </footer>
    </>
  )
}

export default Footer
