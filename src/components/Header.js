import React from "react"
import Link from "gatsby-link"
import { ExternalLink } from "react-feather"
import Logo from "../assets/favicon.png"

const Header = () => (
  <div className="fl-header-container" id="header">
    <div className="fl-header-inner">
      <div className="fl-header-section">
        <img className="fl-header-logo" src={Logo} alt="stitchfix-logo" />
        <Link className="fl-header-link" to="/">Flotilla</Link>
      </div>
      <div className="fl-header-section">
        <Link className="fl-header-link" activeClassName="active" to="/docs/introduction">
          Docs
        </Link>
        <a className="fl-header-link" target="_blank" rel="noopener noreferrer" href="https://github.com/stitchfix/flotilla-os">
          Github&nbsp;<ExternalLink size={16} />
        </a>
      </div>
    </div>
  </div>
)

export default Header
