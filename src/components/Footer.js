import React from "react"
import Link from "gatsby-link"
import { ExternalLink } from "react-feather"
import mapDocsToLinks from "../utils/mapDocsToLinks"

const Footer = props => {
  const docs = mapDocsToLinks(props.docs)
  return (
    <div className="fl-footer-container">
      <div className="fl-footer-inner">
        <div className="fl-footer-col">
          <h4>Stitch Fix</h4>
        </div>
        <div className="fl-footer-col">
          <h4>Docs</h4>
          {!!docs && docs.map(p => (
            <Link
              className="fl-footer-link"
              activeClassName="active"
              key={p.path}
              to={p.path}
            >
              {p.title}
            </Link>
          ))}
          <hr/>
          <a className="fl-footer-link" target="_blank" href="https://github.com/stitchfix/flotilla-os" rel="noopener noreferrer">
            Github <ExternalLink size={16} />
          </a>
        </div>
        <div className="fl-footer-col">


        </div>
      </div>
    </div>
  )
}

export default Footer
