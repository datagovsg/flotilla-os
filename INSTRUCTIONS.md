# Flotilla Landing Page

The landing page is built with [Gatsby](https://www.gatsbyjs.org/) and React.

## Development
```
npm install
gatsby develop
```

## Generating Docs
```
# Source the README.md file from the master branch
git checkout master -- README.md

# The `write-docs` command will split the README into individual files in
# src/docs/
npm run write-docs
```

## Deploying
```
# This command will generate the static build files and deploy them to the
# gh-pages branch. The changes will be available at
# https://stitchfix.github.io/flotilla-os
npm run deploy
```
