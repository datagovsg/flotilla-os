import React, { Component } from "react"
import styled from "styled-components"
import { tween, styler, easing, stagger, timeline } from "popmotion"
import styles from "../constants/styles"
import colors from "../constants/colors"
import LandingPageSection from "./LandingPageSection"
import withScrollAnimation from "./withScrollAnimation"

const LandingFeaturesDiv = styled.div`
  ${styles.mixins.flex("row", "nowrap", "space-between", "center")}
  width: 100%;
  ${styles.mixins.mediaQuery({
    max: styles.breakpoints.small,
    styles: `
      ${styles.mixins.flex("column", "nowrap", "flex-start", "stretch")}
    `
  })}
`

console.log(styles.mixins.mediaQuery({
    max: styles.breakpoints.small,
    styles: `
      ${styles.mixins.flex("column", "nowrap", "flex-start", "stretch")}
    `
  }))

const LandingFeatureDiv = styled.div`
  ${styles.mixins.flex("column", "nowrap", "center", "center", true)}
  color: ${colors.light_gray_4};
  flex: 1;
`

class LandingFeatures extends Component {
  render() {
    return (
      <LandingPageSection
        background={colors.gray_3}
      >
        <LandingFeaturesDiv>
          <LandingFeatureDiv>Feature 1</LandingFeatureDiv>
          <LandingFeatureDiv>Feature 2</LandingFeatureDiv>
          <LandingFeatureDiv>Feature 3</LandingFeatureDiv>
        </LandingFeaturesDiv>
      </LandingPageSection>
    )
  }
}

export default LandingFeatures
