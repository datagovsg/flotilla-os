import React from "react"
import styled from "styled-components"
import styles from "../constants/styles"

const DocContainerDiv = styled.div`
  max-width: 100%;
`
const DocDiv = styled.div`
  word-break: break-all;
`

export default function Template({ data }) {
  const { markdownRemark } = data
  const { frontmatter, html } = markdownRemark
  return (
    <DocContainerDiv>
      <DocDiv>
        <h1>{frontmatter.title}</h1>
        <h2>{frontmatter.date}</h2>
        <div
          className="blog-post-content"
          dangerouslySetInnerHTML={{ __html: html }}
        />
      </DocDiv>
    </DocContainerDiv>
  )
}

export const pageQuery = graphql`
  query BlogPostByPath($path: String!) {
    markdownRemark(frontmatter: { path: { eq: $path } }) {
      html
      frontmatter {
        path
        title
        group
      }
    }
  }
`
