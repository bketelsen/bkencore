/* This example requires Tailwind CSS v2.0+ */
import { Disclosure } from '@headlessui/react'
import { MenuIcon, XIcon } from '@heroicons/react/outline'
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
  { name: 'About', href: '/about' },
]

function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ')
}

const Nav: FC = (props) => {
  const router = useRouter()
  const current = (it: NavItem) => router.asPath.startsWith(it.href)

  return (
    <Disclosure as="nav" className="py-6">
      {({ open }) => (
        <>
          <div className="relative flex items-center justify-between h-16">
            <div className="absolute inset-y-0 left-0 flex items-center sm:hidden">
              {/* Mobile menu button*/}
              <Disclosure.Button className="inline-flex items-center justify-center p-2 rounded-md text-base-content hover:neutral-focus hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white">
                <span className="sr-only">Open main menu</span>
                {open ? (
                  <XIcon className="block w-6 h-6" aria-hidden="true" />
                ) : (
                  <MenuIcon className="block w-6 h-6" aria-hidden="true" />
                )}
              </Disclosure.Button>
            </div>

            <div className="flex items-center justify-center flex-1 text-lg font-semibold sm:items-stretch sm:justify-between text-primary">

              <div className="flex items-center flex-shrink-0">
                <Link href="/"><a className="hover-underline">brian.dev</a></Link>
              </div>

              <div className="hidden sm:block sm:ml-6">
                <div className="flex space-x-4">
                  {navigation.map((item) => (
                    <a
                      key={item.name}
                      href={item.href}
                      className={"hover-underline px-3 py-2"}
                      aria-current={current(item) ? 'page' : undefined}
                    >
                      {item.name}
                    </a>
                  ))}

                </div>

              </div>
            </div>
            
          </div>

          <Disclosure.Panel className="sm:hidden">
            <div className="pt-2 pb-3 space-y-1">
              {navigation.map((item) => (
                <Disclosure.Button
                  key={item.name}
                  as="a"
                  href={item.href}
                  className={classNames(
                    current(item) ? 'bg-gray-900 text-white' : 'dark:text-gray text-gray-800 hover:bg-gray-700 hover:text-white',
                    'block px-3 py-2 rounded-md text-base font-medium'
                  )}
                  aria-current={current(item) ? 'page' : undefined}
                >
                  {item.name}
                </Disclosure.Button>
              ))}
            </div>
          </Disclosure.Panel>
        </>
      )}
    </Disclosure>
  )
}

export default Nav
