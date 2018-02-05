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
          <h1 className="fl-splash-text fl-splash-header" id="splashHeader">
            asdf
          </h1>
          <p className="fl-splash-text" id="splashParagraph">
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Non vero
            blanditiis laborum voluptatum, voluptatibus perspiciatis. Repellat
            quibusdam deleniti similique quos ut cum ipsa adipisci animi ipsum.
            Eligendi ex, id unde!
          </p>
          <pre className="fl-splash-code" id="splashCode">
            <code>git clone git@github.com:stitchfix/flotilla-os.git</code>
            <code>docker-compose up --build -d</code>
          </pre>
          <div className="fl-splash-buttons" id="splashButtons">
            <Link
              className="pl-button pl-intent-primary"
              to="/docs/introduction"
            >
              Get Started
            </Link>
          </div>
          <img
            className="fl-splash-image"
            src="https://user-images.githubusercontent.com/166823/35580026-038ae348-059d-11e8-95e4-f0150400a1a8.png"
            alt="flotilla-ui-image"
            id="splashImage"
          />
        </div>
      </LandingPageSection>
    )
  }
}

export default withScrollAnimation({
  elRef: "x",
  animate: splashMotion,
  heightRenderedBeforeAnimationStarts: 100,
})(LandingSplash)
