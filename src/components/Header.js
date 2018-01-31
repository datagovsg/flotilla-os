import React from "react"
import Link from "gatsby-link"
import styled from "styled-components"
import { BookOpen, Github } from "react-feather"
import Logo from "../assets/favicon.png"
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
  z-index: ${styles.zIndex.header},
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
const HeaderLinkStyles = `${styles.mixins.flex("row", "nowrap", "flex-start", "center", true)}`
const HeaderRouterLink = styled(Link)`
  ${HeaderLinkStyles}
`
const HeaderLink = styled.a`
  ${HeaderLinkStyles}
`
const LogoImg = styled.img`
  width: 40px;
  height: 40px;
  border-radius: 6px;
  margin: 0;
  margin-right: ${styles.shared.spacing}px;
  transform: translateY(-3px);
`

const Header = () => (
  <HeaderContainerDiv>
    <HeaderInnerDiv>
      <HeaderLinksDiv>
        <LogoImg src={Logo} alt="stitchfix-logo" />
        <Link to="/">Flotilla</Link>
      </HeaderLinksDiv>
      <HeaderLinksDiv>
        <HeaderRouterLink activeClassName="active" to="/docs/philosophy">
          <BookOpen />
          Documentation
        </HeaderRouterLink>
        <HeaderLink target="_blank" href="https://github.com/stitchfix/flotilla-os">
          <Github />
          Github
        </HeaderLink>
      </HeaderLinksDiv>
    </HeaderInnerDiv>
  </HeaderContainerDiv>
)

export default Header
