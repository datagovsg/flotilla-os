import React, { Component } from "react"
import Link from "gatsby-link"
import { ExternalLink } from "react-feather"
import mapDocsToLinks from "../utils/mapDocsToLinks"

const DocsSidebar = props => {
  const docs = mapDocsToLinks(props.docs)
  return (
    <div className="fl-docs-sidebar">
      {!!docs &&
        docs.map(p => (
          <Link
            className="fl-docs-sidebar-link"
            activeClassName="active"
            key={p.path}
            to={p.path}
          >
            {p.title}
          </Link>
        ))}
      <a
        className="fl-docs-sidebar-link"
        href="swagger.html"
        target="_blank"
        rel="noopener noreferrer"
      >
        API <ExternalLink size={16} />
      </a>
    </div>
  )
}

export default DocsSidebar
