import { FC } from 'react'
import Link from 'next/link'

import Footer from './Footer'
import Nav from './Nav'
interface NavItem {
  name: string
  href: string
}

const navigation: NavItem[] = [
  { name: 'Blog', href: '/blog' },
  { name: 'Bytes', href: '/bytes' },
]
const Layout: FC<React.PropsWithChildren<{}>> = (props) => {
  return (
    <div className="drawer">
      <input id="my-drawer-3" type="checkbox" className="drawer-toggle" />
      <div className="flex flex-col drawer-content">
        <div className="max-w-5xl mx-auto navbar bg-base-300">
          <div className="flex-none lg:hidden">
            <label htmlFor="my-drawer-3" className="btn btn-square btn-ghost">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="inline-block w-6 h-6 stroke-current"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
            </label>
          </div>
          <div className="flex-1 px-2 mx-2">
            <Link href={'/'}>
              <a>brian.dev</a>
            </Link></div>
          <div className="flex-none hidden lg:block">
            <ul className="menu menu-horizontal">
              {navigation.map((item) => (
                <li key={item.name}>
                  <Link href={item.href}>
                    <a className="hover-underline">{item.name}</a>
                  </Link>
                </li>
              ))}
              <li>
                <Link href={'/about'}>
                  <a className="hover-underline">About</a>
                </Link>
              </li>
            </ul>
          </div>
        </div>
        <div className="flex flex-col w-full max-w-5xl px-4 mx-auto mt-10 lg:px-0">
          {props.children}
        </div>
        <div className="flex-none">
        <Footer />
      </div>
      </div>
      <div className="drawer-side">
        <label htmlFor="my-drawer-3" className="drawer-overlay"></label>
        <ul className="p-4 overflow-y-auto menu w-80 bg-base-100">
          <li>
            <Link href={'/'}>
              <a>Home</a>
            </Link>
          </li>
          {navigation.map((item) => (
            <li key={item.name}>
              <Link href={item.href}>
                <a>{item.name}</a>
              </Link>
            </li>
          ))}
          <li>
            <Link href={'/about'}>
              <a>About</a>
            </Link>
          </li>

        </ul>

      </div>
    </div>
  )
}

export default Layout
