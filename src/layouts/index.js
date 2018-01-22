import React from "react"
import PropTypes from "prop-types"
import Helmet from "react-helmet"
import styled from "styled-components"
import Header from "../components/Header"
import DocsSidebar from "../components/DocsSidebar"
import Page from "../components/Page"
import styles from "../constants/styles"
import "./index.css"

const AppDiv = styled.div`
  font-family: ${styles.shared.fontFamily};
  background: ${styles.app.background};
  width: 100vw;
  min-height: 100vh;

  a {
    color: ${styles.link.color};
    text-decoration: none;
    &:hover {
      color: ${styles.link.hoverColor};
    }
  }
`
const AppInnerDiv = styled.div`
  ${styles.mixins.flex("row", "nowrap", "flex-start", "flex-start")}
`

const TemplateWrapper = (props) => {
  const { children, location } = props
  const shouldRenderSidebar = location.pathname.startsWith("/docs")
  const isLanding = location.pathname === "/"

  return (
    <AppDiv>
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
    </AppDiv>
  )
}

TemplateWrapper.propTypes = {
  children: PropTypes.func,
}

export default TemplateWrapper

export const query = graphql`
  query DocsQuery {
    allMarkdownRemark {
      edges {
        node {
          id
          excerpt(pruneLength: 250)
          frontmatter {
            path
            title
            group
          }
        }
      }
    }
  }
`
