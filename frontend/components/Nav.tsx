import Link from 'next/link'
import { useRouter } from 'next/router'
import { FC } from 'react'

interface NavItem {
  name: string
  href: string
}

const navigation: NavItem[] = [
  { name: 'Blog', href: '/blog' },
  { name: 'Bytes', href: '/bytes' },
]

const Nav: FC = (props) => {
  const router = useRouter()
  const current = (it: NavItem) => router.asPath.startsWith(it.href)

  return (
    <div className="shadow-lg  navbar bg-base-100">
      <div className="navbar-start">
        <div className="dropdown">
          <label tabIndex={0} className="btn btn-ghost lg:hidden">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4 6h16M4 12h8m-8 6h16"
              />
            </svg>
          </label>
          <ul
            tabIndex={0}
            className="p-2 mt-3 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-52"
          >
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
                <a>About Me</a>
              </Link>
            </li>
          </ul>
        </div>
        <Link href={'/'}>
          <a className="text-xl normal-case btn btn-ghost">brian.dev</a>
        </Link>
      </div>
      <div className="hidden navbar-center lg:flex">
        <ul className="p-0 menu menu-horizontal">
          {navigation.map((item) => (
            <li key={item.name}>
              <Link href={item.href}>
                <a>{item.name}</a>
              </Link>
            </li>
          ))}
        </ul>
      </div>
      <div className="navbar-end">
        <Link href={'/about'}>
          <a className="btn">About Me</a>
        </Link>
      </div>
    </div>
  )
}

export default Nav
