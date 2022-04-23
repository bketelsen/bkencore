import { FC } from "react"
import Footer from "./Footer"
import Nav from "./Nav"

const Layout: FC<React.PropsWithChildren<{}>> = (props) => {
  return (
    <div className="flex flex-col h-full">
      <div className="flex-none">
        <Nav />
      </div>
      <div className="flex-grow">
        {props.children}
      </div>
      <div className="flex-none">
        <Footer />
      </div>
    </div>
  )
}

export default Layout
