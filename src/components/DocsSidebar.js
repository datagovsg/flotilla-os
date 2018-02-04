import React, { Component } from "react"
import Link from "gatsby-link"
import mapDocsToLinks from "../utils/mapDocsToLinks"

const DocsSidebar = (props) => {
  const docs = mapDocsToLinks(props.docs)
  return (
    <div className="fl-docs-sidebar">
      {!!docs && docs.map(p => (
        <Link
          className="fl-docs-sidebar-link"
          activeClassName="active"
          key={p.path}
          to={p.path}
        >
          {p.title}
        </Link>
      ))}
    </div>
  )
}

export default DocsSidebar
