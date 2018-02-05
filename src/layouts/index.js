import React from "react"
import PropTypes from "prop-types"
import Helmet from "react-helmet"
import Header from "../components/Header"
import Footer from "../components/Footer"
import DocsSidebar from "../components/DocsSidebar"
import Page from "../components/Page"
import "./index.scss"

const TemplateWrapper = props => {
  const { children, location } = props
  const shouldRenderSidebar = location.pathname.startsWith("/docs")
  const isLanding = location.pathname === "/"

  return (
    <div>
      <Helmet>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta property="og:title" content="Flotilla | Stitch Fix Open Source" />
        <meta
          property="og:description"
          content="Flotilla is a self-service framework that dramatically simplifies the process of defining and executing containerized jobs."
        />
        <meta property="og:type" content="website" />
        <meta
          property="og:url"
          content="http://github.stitchfix.io/flotilla-os"
        />
        <meta
          property="og:image"
          content="https://user-images.githubusercontent.com/166823/35580026-038ae348-059d-11e8-95e4-f0150400a1a8.png"
        />
        <meta property="og:image:type" content="image/png" />
        <meta property="og:image:width" content="2880" />
        <meta property="og:image:height" content="1800" />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:site" content="@stitchfix_algo" />
        <meta
          name="twitter:image"
          content="https://user-images.githubusercontent.com/166823/35580026-038ae348-059d-11e8-95e4-f0150400a1a8.png"
        />
        <title>Flotilla | Stitch Fix Open Source</title>
      </Helmet>
      <Header />
      <Page isLanding={isLanding}>
        {shouldRenderSidebar && <DocsSidebar docs={props.data} />}
        {children()}
      </Page>
      <Footer docs={props.data} />
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
