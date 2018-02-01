import React, { Component } from "react"
import Link from "gatsby-link"

const mapPostsToTitles = (posts) => posts.allMarkdownRemark.edges.map(p => p.node.frontmatter)

const DocsSidebar = (props) => {
  const posts = mapPostsToTitles(props.posts)
  return (
    <div className="fl-docs-sidebar">
      {!!posts && posts.map(p => (
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
