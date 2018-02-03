// Splits README.md, writes each chunk to src/docs/, which is then used by
// ./gatsby-node.js to create pages.
const fs = require("fs")
const path = require("path")

//
// Helper functions
//
// Slice the "## "
const getCleanTitle = str => str.slice(3)

// Lowercase; replace spaces with "-"
const getPath = str => str.toLowerCase().replace(/\s+/g, "-")

// Append frontmatter - necessary for Gatsby to query.
const withFrontmatter = (line, i, group) => {
  const path = `/${group}/${getPath(getCleanTitle(line))}`
  const title = getCleanTitle(line)

  return `---\npath: "${path}"\ntitle: "${title}"\ngroup: "${group}"\nindex: "${i}"\n---\n`
}

// The filename.
const getFileName = (line, index, group) => {
  return `${group}-${index}-${getPath(getCleanTitle(line))}`
}

Promise.resolve()
  .then(() => {
    // Read it.
    const readme = fs.readFileSync("README.md", "utf-8")

    // Split it.
    const lines = readme.split("\n")

    // Group lines into chunks, delimited by "## ".
    const pages = lines.reduce((acc, line, i) => {
      if (line.startsWith("## ")) {
        acc = [
          ...acc,
          {
            filename: getFileName(line, acc.length, "usage"),
            text: withFrontmatter(line, acc.length, "usage")
          }
        ]
      } else if (acc.length > 0) {
        const index = acc.length - 1
        acc[index].text += `${line}\n`
      }
      return acc
    }, [])

    return pages
  })
  .then((pages) => {
    // Write the files to src/docs/$filename.md
    pages.forEach(page => {
      const pagepath = path.resolve("src", "docs", `${page.filename}.md`)
      fs.writeFile(pagepath, page.text, (err) => {
        if (err) return console.error(err)

        console.log(`${page.filename}.md was written to ${pagepath}`)
      })
    })
  })
