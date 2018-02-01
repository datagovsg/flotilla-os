import React from "react"
import PropTypes from "prop-types"
import cn from "classnames"

const Page = ({ children, isLanding }) => {
  const innerclass = cn({
    "fl-page-inner": true,
    "is-landing": isLanding,
  })

  return (
    <div className="fl-page-container">
      <div className={innerclass}>
        {children}
      </div>
    </div>
  )
}

Page.propTypes = {
  // Whether or not this is the landing page, which requires some different styles.
  isLanding: PropTypes.bool.isRequired,
}
Page.defaultProps = {
  isLanding: false,
}

export default Page
