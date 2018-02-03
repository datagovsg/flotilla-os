import React from "react"

const LandingPageSection = (props) => {
  return (
    <div className="fl-landing-section-container" id={props.id} ref={props.innerRef}>
      <div className="fl-landing-section-inner">
        {props.children}
      </div>
    </div>
  )
}

export default LandingPageSection
