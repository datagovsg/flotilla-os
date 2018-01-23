import React, { Component } from "react"
import { isEmpty } from "lodash"
import Link from "gatsby-link"
import styled from "styled-components"
import styles from "../constants/styles"
import colors from "../constants/colors"

const DocsSidebarContainerDiv = styled.div`
  width: ${styles.sidebar.width}px;
  min-width: ${styles.sidebar.width}px;
  background: ${styles.sidebar.background};
  margin-right: ${styles.shared.spacing * 2}px;
`
const DocsSidebarInnerDiv = styled.div`
  ${styles.mixins.flex("column", "nowrap", "flex-start", "stretch")}
`
const DocsSidebarLink = styled(Link)`
  ${styles.mixins.flex("row", "nowrap", "flex-start", "center")}
  height: 40px;
  padding: ${styles.shared.spacing}px;
  &.active,
  &:hover {
    color: ${colors.black_0} !important;
    background: ${colors.gray_3};
  }
`

const mapPostsToTitles = (posts) => posts.allMarkdownRemark.edges.map(p => p.node.frontmatter)

const DocsSidebar = (props) => {
  const posts = mapPostsToTitles(props.posts)
  return (
    <DocsSidebarContainerDiv>
      <DocsSidebarInnerDiv>
        {!!posts && posts.map(p => (
          <DocsSidebarLink
            activeClassName="active"
            key={p.path}
            to={p.path}
          >
            {p.title}
          </DocsSidebarLink>
        ))}
      </DocsSidebarInnerDiv>
    </DocsSidebarContainerDiv>
  )
}

export default DocsSidebar
