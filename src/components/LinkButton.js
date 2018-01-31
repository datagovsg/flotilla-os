import React from "react"
import { Link } from "react-router-dom"
import styled from "styled-components"
import styles from "../constants/styles"
import colors from "../constants/colors"

const LinkButton = styled(Link)`
  ${styles.mixins.flex("row", "nowrap", "center", "center", true)}
  outline: none;
  -webkit-appearance: none;
  border: none;
  background: ${colors.blue_0};
  color: white !important;
  padding: ${styles.shared.spacing}px ${styles.shared.spacing * 2}px;
  line-height: 1;
  &:hover {
    background: ${colors.blue_2};
  }
`

export default LinkButton
