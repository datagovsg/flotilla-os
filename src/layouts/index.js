import React from "react"
import PropTypes from "prop-types"
import Helmet from "react-helmet"
import Header from "../components/Header"
import DocsSidebar from "../components/DocsSidebar"
import Page from "../components/Page"
import "./index.scss"

const TemplateWrapper = (props) => {
  const { children, location } = props
  const shouldRenderSidebar = location.pathname.startsWith("/usage") || location.pathname.startsWith("/api")
  const isLanding = location.pathname === "/"

  return (
    <div>
      <Helmet
        title="Flotilla | Stitch Fix"
        meta={[
          { name: "description", content: "Sample" },
          { name: "keywords", content: "sample, something" },
        ]}
      />
      <Header />
      <Page isLanding={isLanding}>
        {shouldRenderSidebar && <DocsSidebar posts={props.data} />}
        {children()}
      </Page>
    </div>
  )
}

TemplateWrapper.propTypes = {
  children: PropTypes.func,
}

export default TemplateWrapper

export const query = graphql`
  query DocsQuery {
    allMarkdownRemark(sort: { order: ASC, fields: [frontmatter___index] }) {
      edges {
        node {
          id
          excerpt(pruneLength: 250)
          frontmatter {
            path
            title
            group
            index
          }
        }
      }
    }
  }
`
