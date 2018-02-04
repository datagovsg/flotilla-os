import React, { Component } from "react"
import Link from "gatsby-link"
import LandingPageSection from "./LandingPageSection"
import withScrollAnimation from "./withScrollAnimation"
import splashMotion from "../motion/splash"
import SplashImage from "../assets/flot-img-1.png"

class LandingSplash extends Component {
  render() {
    return (
      <LandingPageSection id="splash" innerRef={this.props.innerRef}>
        <div className="fl-landing-section-content fl-splash">
          <h1 className="fl-splash-text fl-splash-header" id="splashHeader">Self-service data science.</h1>
          <p className="fl-splash-text" id="splashParagraph">Flotilla is a self-service framework that dramatically simplifies the process of defining and executing containerized jobs. This means you get to focus on the work you're doing rather than how to do it.</p>
          <div className="fl-splash-buttons" id="splashButtons">
            <Link className="pl-button pl-intent-primary" to="/usage/introduction">Get Started</Link>
          </div>
          <img src={SplashImage} alt="flotilla-ui-image" id="splashImage" />
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
