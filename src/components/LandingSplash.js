import React, { Component } from "react"
import Link from "gatsby-link"
import styled from "styled-components"
import { tween, styler, easing, stagger, timeline } from "popmotion"
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
const LandingSplashTagline = styled.div.attrs({
  id: "landingSplashTagline",
})`
  font-size: 20px;
`
const LandingSplashLargeText = styled.div.attrs({
  id: "landingSplashLargeText",
})`
  font-size: 52px;
  line-height: 1.125;
  max-width: 600px;
`
const LandingSplashButtonGroup = styled.div.attrs({
  id: "landingSplashButtonGroup",
})`
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

// Animations

class LandingSplash extends Component {
  componentDidMount() {
    const landingSplashTagline = styler(document.querySelector("#landingSplashTagline"))
    const landingSplashLargeText = styler(document.querySelector("#landingSplashLargeText"))
    const landingSplashButtonGroup = styler(document.querySelector("#landingSplashButtonGroup"))

    const staggerDuration = -300
    const motionDuration = 500
    const initialYPos = 80

    const sharedMotion = {
      from: { opacity: 0, y: initialYPos },
      to: { opacity: 1, y: 0 },
      ease: easing.ease,
      duration: motionDuration,
    }

    timeline([
      { ...sharedMotion, track: "tagline", },
      `${staggerDuration}`,
      { ...sharedMotion, track: "largeText", },
      `${staggerDuration}`,
      { ...sharedMotion, track: "buttonGroup", },
    ]).start(({ tagline, largeText, buttonGroup }) => {
      landingSplashTagline.set(tagline)
      landingSplashLargeText.set(largeText)
      landingSplashButtonGroup.set(buttonGroup)
    })
  }
  render() {
    return (
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
  }
}

export default LandingSplash
