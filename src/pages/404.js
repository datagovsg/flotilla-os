import React from "react"
import Link from "gatsby-link"

const NotFoundPage = () => (
  <div>
    <h1>This route doesn't exist 😔</h1>
    <Link to="/docs/introduction">View the docs instead?</Link>
  </div>
)

export default NotFoundPage
