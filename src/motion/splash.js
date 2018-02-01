import { tween, styler, easing, stagger, timeline } from "popmotion"

const splash = () => {
  const landingSplashTagline = styler(document.querySelector("#landingSplashTagline"))
  const landingSplashLargeText = styler(document.querySelector("#landingSplashLargeText"))
  const landingSplashButtonGroup = styler(document.querySelector("#landingSplashButtonGroup"))

  const staggerDuration = -350
  const motionDuration = 500
  const initialYPos = 120

  const sharedMotion = {
    from: { opacity: 0, y: initialYPos },
    to: { opacity: 1, y: 0 },
    ease: easing.ease,
    duration: motionDuration,
  }

  timeline([
    { ...sharedMotion, track: "tagline", },
    `${staggerDuration}`,
    { ...sharedMotion, track: "largeText", },
    `${staggerDuration}`,
    { ...sharedMotion, track: "buttonGroup", },
    `${-500}`,
  ]).start(({ tagline, largeText, buttonGroup }) => {
    landingSplashTagline.set(tagline)
    landingSplashLargeText.set(largeText)
    landingSplashButtonGroup.set(buttonGroup)
  })
}

export default splash
