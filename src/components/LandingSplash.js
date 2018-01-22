import React from "react"
import Link from "gatsby-link"
import styled from "styled-components"
import { TransitionGroup } from "react-transition-group"
import { MotionValue } from "popmotion-react"
import { tween, transform } from "popmotion"
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
  // margin-bottom: ${styles.shared.spacing * 1.5}px;
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

// Motion
const stateChangeHandlers = {
  componentWillAppear: ({ value, complete }) => tween({
    from: value.get(),
    to: 1,
  }).start({
    update: value,
    complete,
  })
}

const LandingSplash = () => (
  <TransitionGroup>
    <MotionValue key="landing-motion" onStateChange={stateChangeHandlers}>
      {(popmotionState) => {
        console.log(popmotionState)
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
                <LandingSplashButton>Github</LandingSplashButton>
              </LandingSplashButtonGroup>
            </LandingSplashInnerDiv>
          </LandingSplashContainerDiv>
        )
      }}
    </MotionValue>
  </TransitionGroup>
)

export default LandingSplash
