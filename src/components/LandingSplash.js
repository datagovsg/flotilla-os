import React from "react"
import Link from "gatsby-link"
import styled from "styled-components"
import styles from "../constants/styles"
import colors from "../constants/colors"

// Styles
const LandingSplashContainerDiv = styled.div`
  ${styles.mixins.containerComponentSharedStyles()}
  ${styles.mixins.flex("row", "nowrap", "center", "flex-start")}
  background: ${colors.gray_4};
`
const LandingSplashInnerDiv = styled.div`
  ${styles.mixins.innerComponentSharedStyles()}
  ${styles.mixins.flex("column", "nowrap", "center", "flex-start", true, styles.shared.spacing * 2)}
  height: 600px;
`
const LandingSplashTagline = styled.div`
  font-size: 20px;
`
const LandingSplashLargeText = styled.div`
  font-size: 52px;
  line-height: 1.125;
  max-width: 600px;
`
const LandingSplashButtonGroup = styled.div`
  ${styles.mixins.flex("row", "nowrap", "flex-start", "center", true)}
`
const LandingSplashButton = styled(Link)`
  ${styles.mixins.flex("row", "nowrap", "center", "center", true)}
  outline: none;
  background: transparent;
  -webkit-appearance: none;
  border: none;
  background: ${colors.blue_0};
  color: white !important;
  padding: ${styles.shared.spacing}px ${styles.shared.spacing * 2}px;
  line-height: 1;
`

const LandingSplash = () => (
  <LandingSplashContainerDiv>
    <LandingSplashInnerDiv>
      <LandingSplashTagline>
        You don't need data engineers. You need Flotilla.
      </LandingSplashTagline>
      <LandingSplashLargeText>
        Some words about running containers and data science and stuff.
      </LandingSplashLargeText>
      <LandingSplashButtonGroup>
        <LandingSplashButton to="/docs/quickstart">Documentation</LandingSplashButton>
        <LandingSplashButton to="/">Github</LandingSplashButton>
      </LandingSplashButtonGroup>
    </LandingSplashInnerDiv>
  </LandingSplashContainerDiv>
)

export default LandingSplash
