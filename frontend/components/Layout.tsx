import { FC } from 'react'
import Footer from './Footer'
import Nav from './Nav'

const Layout: FC<React.PropsWithChildren<{}>> = (props) => {
  return (
    <div className="flex flex-col w-full max-w-5xl px-4 mx-auto lg:px-0">
      <div className="flex-none mb-10">
        <Nav />
      </div>
      <div className="justify-center flex-grow">{props.children}</div>
      <div className="flex-none">
        <Footer />
      </div>
    </div>
  )
}

export default Layout
