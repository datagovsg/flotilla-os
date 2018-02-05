import React, { Component } from "react"
import { has } from "lodash"

const validateExistenceAndType = (object, key, type) =>
  has(object, key) && typeof object[key] === type
const printErr = (key, type) => {
  console.error(
    `The opts object passed to the initial call of withScrollAnimation() ` +
      `must have a [${key}] property of type [${type}].`
  )
}

const validateOpts = opts => {
  const keysToCheck = [
    { key: "elRef", type: "string" },
    { key: "animate", type: "function" },
    { key: "heightRenderedBeforeAnimationStarts", type: "number" },
  ]

  keysToCheck.forEach(({ key, type }) => {
    if (!validateExistenceAndType(opts, key, type)) {
      printErr(key, type)
    }
  })
}

const withScrollAnimation = (opts = {}) => Unwrapped => {
  validateOpts(opts)

  const { elRef, animate, heightRenderedBeforeAnimationStarts } = opts

  return class Wrapped extends Component {
    static displayName = `withScrollAnimation(${Unwrapped.displayName ||
      "UnwrappedComponent"})`
    constructor(props) {
      super(props)
      this.scrollListener = this.scrollListener.bind(this)
    }
    state = {
      hasAnimated: false,
    }
    componentDidMount() {
      if (this.shouldAnimate()) {
        this.setState({ hasAnimated: true }, () => {
          animate()
        })
      } else {
        window.addEventListener("scroll", this.scrollListener)
      }
    }
    componentDidUpdate(prevProps, prevState) {
      if (!prevState.hasAnimated && !!this.state.hasAnimated) {
        window.removeEventListener("scroll", this.scrollListener)
      }
    }
    scrollListener() {
      if (this.shouldAnimate()) {
        this.setState({ hasAnimated: true }, () => {
          animate()
        })
      }
    }
    shouldAnimate() {
      if (this.state.hasAnimated) return false

      const viewportHeight = window.innerHeight
      const elTop = this[elRef].getBoundingClientRect().top

      if (viewportHeight - elTop >= heightRenderedBeforeAnimationStarts) {
        return true
      }

      return false
    }
    render() {
      return (
        <Unwrapped
          {...this.props}
          innerRef={r => {
            this[elRef] = r
          }}
        />
      )
    }
  }
}

export default withScrollAnimation
