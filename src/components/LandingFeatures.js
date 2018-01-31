import React, { Component } from "react"
import styled from "styled-components"
import { tween, styler, easing, stagger, timeline } from "popmotion"
import styles from "../constants/styles"
import colors from "../constants/colors"
import LandingPageSection from "./LandingPageSection"
import withScrollAnimation from "./withScrollAnimation"
import LinkButton from "./LinkButton"

const Features = styled.div`
  ${styles.mixins.flex("column", "nowrap", "flex-start", "flex-start")}
  width: 100%;
`

const PrimaryFeature = styled.div`
  ${styles.mixins.flex("column", "nowrap", "center", "flex-start", true, styles.shared.spacing * 1.5)}
  flex: 1;
  max-width: 600px;
  padding: ${styles.shared.spacing * 2}px 0;
`
const PrimaryFeatureText = styled.h4`
  line-height: 1.5;
`
const FeaturesList = styled.ul`
  max-width: 600px;
`

class LandingFeatures extends Component {
  render() {
    return (
      <LandingPageSection
        background={colors.light_gray_3}
      >
        <Features>
          <PrimaryFeature>
            <h2>What is Flotilla?</h2>
            <PrimaryFeatureText>Flotilla is a self-service framework that dramatically simplifies the process of defining and executing containerized jobs. This means you get to focus on the work you're doing rather than how to do it.</PrimaryFeatureText>
          </PrimaryFeature>
          <FeaturesList>
            <li>Define containerized jobs by allowing you to specify exactly what command to run, what image to run that command in, and what resources that command needs to run.</li>
            <li>Run any previously defined job and access its logs, status, and exit code.</li>
            <li>View and edit job definitions with a flexible UI.</li>
            <li>Run jobs and view execution history and logs within the UI.</li>
            <li>Use the complete REST API for definitions, jobs, and logs to build your own custom workflows.</li>
          </FeaturesList>
          <LinkButton to="/docs/quickstart">View Quickstart Guide</LinkButton>
        </Features>
      </LandingPageSection>
    )
  }
}

export default LandingFeatures
