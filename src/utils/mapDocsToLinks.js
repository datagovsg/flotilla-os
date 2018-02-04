export default function mapDocsToLinks(docs) {
  return docs.allMarkdownRemark.edges.map(p => p.node.frontmatter)
}
