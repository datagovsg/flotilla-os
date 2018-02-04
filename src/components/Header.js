import React from "react"
import Link from "gatsby-link"
import Logo from "../assets/favicon.png"

const Header = () => (
  <div className="fl-header-container" id="header">
    <div className="fl-header-inner">
      <div className="fl-header-section">
        <img className="fl-header-logo" src={Logo} alt="stitchfix-logo" />
        <Link className="fl-header-link" to="/">Flotilla</Link>
      </div>
      <div className="fl-header-section">
        <Link className="fl-header-link" activeClassName="active" to="/usage/introduction">
          Usage
        </Link>
        <a className="fl-header-link" target="_blank" href="https://github.com/stitchfix/flotilla-os">
          Github
        </a>
      </div>
    </div>
  </div>
)

export default Header
