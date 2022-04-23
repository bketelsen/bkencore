import { FC } from "react"
import Nav from "./Nav"

const Layout: FC<React.PropsWithChildren<{}>> = (props) => {
  return (
    <div>
      <Nav />
      {props.children}
    </div>
  )
}

export default Layout
