// Splits README.md, writes each chunk to src/docs/, which is then used by
// ./gatsby-node.js to create pages.
const fs = require("fs")
const path = require("path")

//
// Helper functions
//
// Append frontmatter - necessary for Gatsby to query.
const withFrontmatter = (text) => {
  return `---\npath: "/api"\ntitle: "API Documentation"\ngroup: "api"\nindex: "0"\n---\n${text}\n`
}

// The filename.
const getFileName = (index, group) => {
  return `${group}-${index}`
}

Promise.resolve()
  .then(() => {
    const swagpath = path.resolve("docs", "swagger-spec.md")
    // Read it.
    const swag = fs.readFileSync(swagpath, "utf-8")

    return ({
      filename: getFileName(0, "api"),
      text: withFrontmatter(swag),
    })
  })
  .then((page) => {
    const pagepath = path.resolve("src", "docs", `${page.filename}.md`)
    fs.writeFile(pagepath, page.text, (err) => {
      if (err) return console.error(err)

      console.log(`${page.filename}.md was written to ${pagepath}`)
    })
  })
