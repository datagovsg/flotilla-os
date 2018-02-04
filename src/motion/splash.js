import { tween, styler, easing, stagger, timeline } from "popmotion"

const splash = () => {
  const splashHeader = styler(document.querySelector("#splashHeader"))
  const splashParagraph = styler(document.querySelector("#splashParagraph"))
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
      { ...sharedMotion, track: "splashButtonsTrack", },
      `${staggerDuration}`,
      { ...sharedMotion, track: "splashImageTrack", },
    ]).start(({ splashHeaderTrack, splashParagraphTrack, splashButtonsTrack, splashImageTrack }) => {
      splashHeader.set(splashHeaderTrack)
      splashParagraph.set(splashParagraphTrack)
      splashButtons.set(splashButtonsTrack)
      splashImage.set(splashImageTrack)
    })
  })
}

export default splash
