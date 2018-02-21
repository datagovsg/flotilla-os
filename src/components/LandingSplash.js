import React, { Component } from "react"
import Link from "gatsby-link"
import LandingPageSection from "./LandingPageSection"
import withScrollAnimation from "./withScrollAnimation"
import splashMotion from "../motion/splash"

class LandingSplash extends Component {
  componentDidMount() {
    splashMotion()
  }
  render() {
    return (
      <LandingPageSection id="splash" innerRef={this.props.innerRef}>
        <div className="fl-landing-section-content fl-splash">
          <h1
            className="fl-splash-text fl-splash-header fl-initial-opacity-zero"
            id="splashHeader"
          >
            Own your work end-to-end.
          </h1>
          <p
            className="fl-splash-text fl-initial-opacity-zero"
            id="splashParagraph"
          >
            Flotilla is a self-service batch job execution framework that dramatically simplifies the process of defining and executing containerized jobs. Focus on the work you're doing rather than how to do it.
          </p>
          <pre
            className="fl-splash-code fl-initial-opacity-zero"
            id="splashCode"
          >
            <code>git clone git@github.com:stitchfix/flotilla-os.git</code>
            <code>docker-compose up --build -d</code>
          </pre>
          <div
            className="fl-splash-buttons fl-initial-opacity-zero"
            id="splashButtons"
          >
            <Link
              className="pl-button pl-intent-primary"
              to="/docs/introduction"
            >
              Get Started
            </Link>
          </div>
          <img
            className="fl-splash-image fl-initial-opacity-zero"
            src="https://user-images.githubusercontent.com/10807627/36499493-4842176c-16f6-11e8-9467-a345987bd407.png"
            alt="flotilla-ui-image"
            id="splashImage"
          />
        </div>
      </LandingPageSection>
    )
  }
}

export default LandingSplash
