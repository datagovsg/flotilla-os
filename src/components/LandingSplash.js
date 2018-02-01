import React, { Component } from "react"
import Link from "gatsby-link"
import LinkButton from "./LinkButton"
import LandingPageSection from "./LandingPageSection"
import withScrollAnimation from "./withScrollAnimation"
import splashMotion from "../motion/splash"

class LandingSplash extends Component {
  render() {
    return (
      <LandingPageSection id="splash" innerRef={this.props.innerRef}>
        <div className="fl-landing-section-content fl-splash">
          <div className="fl-splash-tagline" id="landingSplashTagline">
            You don't need data engineers. You need Flotilla.
          </div>
          <div className="fl-splash-large-text" id="landingSplashLargeText">
            Some words about running containers and data science and stuff.
          </div>
          <div id="landingSplashButtonGroup">
            <LinkButton to="/docs/quickstart">Get Started</LinkButton>
          </div>
        </div>
      </LandingPageSection>
    )
  }
}

export default withScrollAnimation({
  elRef: "x",
  animate: splashMotion,
  heightRenderedBeforeAnimationStarts: 100
})(LandingSplash)
