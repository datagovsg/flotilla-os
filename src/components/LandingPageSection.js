import React from "react"
import styled from "styled-components"
import styles from "../constants/styles"
import colors from "../constants/colors"

const LandingPageSectionContainerDiv = styled.div`
  ${styles.mixins.containerComponentSharedStyles()}
  ${styles.mixins.flex("row", "nowrap", "center", "flex-start")}
  background: ${props => !!props.background ? props.background : styles.app.background};
  ${props => !!props.height ? `height: ${props.height}px;` : ""}
`
const LandingPageSectionInnerDiv = styled.div`
  ${styles.mixins.innerComponentSharedStyles()}
  height: 100%;
`

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
