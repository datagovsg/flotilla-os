// Note: while most of the components animated by the `splash()` function are
// located at src/components/LandingSplash.js, we're also animating some of the
// elements in the <Header> component (src/components/Header.js)
import { tween, styler, easing, stagger, timeline, keyframes } from "popmotion"

const splash = () => {
  const header = styler(document.querySelector("#header"))
  const splashHeader = styler(document.querySelector("#splashHeader"))
  const splashParagraph = styler(document.querySelector("#splashParagraph"))
  const splashCode = styler(document.querySelector("#splashCode"))
  const splashButtons = styler(document.querySelector("#splashButtons"))
  const splashImage = styler(document.querySelector("#splashImage"))

  const staggerDuration = -350
  const motionDuration = 500
  const initialYPos = 120

  const sharedMotion = {
    from: { opacity: 0, y: initialYPos },
    to: { opacity: 1, y: 0 },
    ease: easing.ease,
    duration: motionDuration,
  }

  window.addEventListener("DOMContentLoaded", () => {
    timeline([
      { ...sharedMotion, track: "splashHeaderTrack", },
      `${staggerDuration}`,
      { ...sharedMotion, track: "splashParagraphTrack", },
      `${staggerDuration}`,
      { ...sharedMotion, track: "splashCodeTrack", },
      `${staggerDuration}`,
      { ...sharedMotion, track: "splashButtonsTrack", },
      `${staggerDuration}`,
      { ...sharedMotion, track: "splashImageTrack", },
      `${staggerDuration}`,
      {
        from: { opacity: 0, y: -80 },
        to: { opacity: 1, y: 0 },
        ease: easing.backInOut,
        duration: 1000,
        track: "headerTrack"
      },
    ]).start(({ headerTrack, splashHeaderTrack, splashParagraphTrack, splashCodeTrack, splashButtonsTrack, splashImageTrack }) => {
      splashHeader.set(splashHeaderTrack)
      splashParagraph.set(splashParagraphTrack)
      splashCode.set(splashCodeTrack)
      splashButtons.set(splashButtonsTrack)
      splashImage.set(splashImageTrack)
      header.set(headerTrack)
    })
  })
}

export default splash
