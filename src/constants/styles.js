import { css } from "styled-components"
import colors from "./colors"

const shared = {
  maxWidth: 1200,
  color: colors.light_gray_1,
  spacing: 12,
  fontFamily: `"MaisonNeue-Book"`,
  primaryColor: colors.blue_0,
}

const breakpoints = {
  small: 480,
  medium: 800,
  large: 1200,
}

const link = {
  color: shared.primaryColor,
  hoverColor: colors.blue_3,
}

const app = {
  background: colors.light_gray_3,
}

const header = {
  height: 100,
  background: colors.black_0,
}

const sidebar = {
  width: 240,
  background: colors.light_gray_1,
}

const containerComponentSharedStyles = () => `
  width: 100vw;
`

const innerComponentSharedStyles = () => `
  width: 100%;
  max-width: ${shared.maxWidth}px;
  padding: ${shared.spacing * 2}px;
`

const flex = (direction, wrap, justify, align, withChildMargin = false, childMargin = shared.spacing) => {
  let childMarginStr = ""
  let alignKey = "align-items"

  if (wrap === "wrap") {
    alignKey = "align-content"
  }

  if (withChildMargin) {
    const marginDirection = direction === "row" ? "right" : "bottom"
    childMarginStr = `& > * { margin-${marginDirection}: ${childMargin}px; }`
  }

  return `display: flex; flex-flow: ${direction} ${wrap}; justify-content: ${justify}; ${alignKey}: ${align}; ${childMarginStr}`
}

const mediaQuery = ({
  min = 0,
  max = 1000000,
  styles = ""
} = {}) => `
  @media (min-width: ${min}px) and (max-width: ${max}px) {
    ${styles}
  }
`

export default {
  shared,
  breakpoints,
  header,
  app,
  link,
  sidebar,
  mixins: {
    containerComponentSharedStyles,
    innerComponentSharedStyles,
    flex,
    mediaQuery,
  }
}
