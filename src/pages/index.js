import React from "react"
import Link from "gatsby-link"
import styled from "styled-components"
import styles from "../constants/styles"
import colors from "../constants/colors"
import LandingSplash from "../components/LandingSplash"
import LandingFeatures from "../components/LandingFeatures"

const IndexPage = () => (
  <div>
    {/*<LandingSplash />*/}
    <LandingFeatures />
    <LandingSplash />
  </div>
)

export default IndexPage
