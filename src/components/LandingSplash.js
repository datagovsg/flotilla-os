import React, { Component } from "react"
import Link from "gatsby-link"
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
            Run ETL, model training jobs and productions jobs with Flotilla.
          </div>
          <div id="landingSplashButtonGroup">
            <Link className="pl-button pl-intent-primary" to="/usage/quick-start">Get Started</Link>
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
