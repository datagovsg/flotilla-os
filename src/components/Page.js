import React from "react"
import PropTypes from "prop-types"
import styled from "styled-components"
import styles from "../constants/styles"

const PageContainerDiv = styled.div`
  ${styles.mixins.containerComponentSharedStyles()}
  ${styles.mixins.flex("row", "nowrap", "center", "flex-start")}
`
const PageInnerDiv = styled.div`
  ${styles.mixins.flex("row", "nowrap", "flex-start", "flex-start")}
  ${styles.mixins.innerComponentSharedStyles()}
  margin-top: ${styles.header.height}px;
  padding: ${props => props.isLanding ? 0 : styles.shared.spacing * 2}px;
  max-width: ${props => props.isLanding ? "100vw" : styles.shared.maxWidth};
`

const Page = ({ children, isLanding }) => (
  <PageContainerDiv>
    <PageInnerDiv isLanding={isLanding}>
      {children}
    </PageInnerDiv>
  </PageContainerDiv>
)

Page.propTypes = {
  // Whether or not this is the landing page, which requires some different styles.
  isLanding: PropTypes.bool.isRequired,
}
Page.defaultProps = {
  isLanding: false,
}

export default Page
