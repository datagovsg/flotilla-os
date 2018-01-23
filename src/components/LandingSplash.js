import React, { Component } from "react"
import Link from "gatsby-link"
import styled from "styled-components"
import { tween, styler, easing, stagger, timeline } from "popmotion"
import styles from "../constants/styles"
import colors from "../constants/colors"
import LandingPageSection from "./LandingPageSection"
import withScrollAnimation from "./withScrollAnimation"

// Styles
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

const animate = () => {
  const landingSplashTagline = styler(document.querySelector("#landingSplashTagline"))
  const landingSplashLargeText = styler(document.querySelector("#landingSplashLargeText"))
  const landingSplashButtonGroup = styler(document.querySelector("#landingSplashButtonGroup"))

  const staggerDuration = -350
  const motionDuration = 500
  const initialYPos = 120

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

class LandingSplash extends Component {
  render() {
    return (
      <LandingPageSection
        height={600}
        background={colors.gray_4}
        innerRef={this.props.innerRef}
      >
        <LandingSplashTagline>
          You don't need data engineers. You need Flotilla.
        </LandingSplashTagline>
        <LandingSplashLargeText>
          Some words about running containers and data science and stuff.
        </LandingSplashLargeText>
        <LandingSplashButtonGroup>
          <LandingSplashButton to="/docs/quickstart">Get Started</LandingSplashButton>
        </LandingSplashButtonGroup>
      </LandingPageSection>
    )
  }
}

export default withScrollAnimation({
  elRef: "x",
  animate,
  heightRenderedBeforeAnimationStarts: 100
})(LandingSplash)
