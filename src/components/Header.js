import React from "react"
import Link from "gatsby-link"
import Logo from "../assets/favicon.png"

const Header = () => (
  <div className="fl-header-container">
    <div className="fl-header-inner">
      <div className="fl-header-section">
        <img className="fl-header-logo" src={Logo} alt="stitchfix-logo" />
        <Link className="fl-header-link" to="/">Flotilla</Link>
      </div>
      <div className="fl-header-section with-horizontal-child-margin">
        <Link className="fl-header-link" activeClassName="active" to="/usage/philosophy">
          Documentation
        </Link>
        <a className="fl-header-link" target="_blank" href="https://github.com/stitchfix/flotilla-os">
          Github
        </a>
      </div>
    </div>
  </div>
)

export default Header
