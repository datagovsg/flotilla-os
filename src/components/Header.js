import React from "react"
import Link from "gatsby-link"
import styled from "styled-components"
import styles from "../constants/styles"
import colors from "../constants/colors"

const HeaderContainerDiv = styled.div`
  ${styles.mixins.containerComponentSharedStyles()}
  ${styles.mixins.flex("row", "nowrap", "center", "center")}
  height: ${styles.header.height}px;
  background: ${styles.header.background};
  color: ${styles.shared.color};
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  a {
    color: ${styles.shared.color} !important;
    border-bottom: 2px solid transparent;
    &:hover,
    &.active {
      color: ${styles.link.hoverColor} !important;
      border-bottom-color: ${styles.shared.primaryColor};
    }
  }
`
const HeaderInnerDiv = styled.div`
  ${styles.mixins.flex("row", "nowrap", "space-between", "center")}
  ${styles.mixins.innerComponentSharedStyles()}
`
const HeaderLinksDiv = styled.div`
  ${styles.mixins.flex("row", "nowrap", "flex-start", "center", true, styles.shared.spacing * 1.5)}
`

const Header = () => (
  <HeaderContainerDiv>
    <HeaderInnerDiv>
      <div>
        <Link to="/">Flotilla</Link>
      </div>
      <HeaderLinksDiv>
        <Link activeClassName="active" to="/docs/philosophy">Documentation</Link>
        <a target="_blank" href="https://github.com/stitchfix/flotilla-os">Github</a>
      </HeaderLinksDiv>
    </HeaderInnerDiv>
  </HeaderContainerDiv>
)

export default Header
