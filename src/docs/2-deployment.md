---
path: "/docs/deploying"
title: "Deploying"
group: "deploying"
---
In a production deployment you'll want multiple instances of the flotilla service running and postgres running elsewhere (eg. Amazon RDS). In this case the most salient detail configuration detail is the `DATABASE_URL`.

### Docker based deploy (simple)

The simplest way to deploy for very light usage is to avoid a reverse proxy and deploy directly with docker.

1. Build and tag an image for flotilla using the `Dockerfile` provided in this repo:

  ```
  docker build -t <your repo name>/flotilla:<version tag>
  ```
2. Run this image wherever you deploy your services:

  ```
  docker run -e DATABASE_URL=<your db url> -e FLOTILLA_MODE=prod -p 3000:3000 ...<other standard docker run args>
  ```

  > Note: Flotilla uses [viper](https://github.com/spf13/viper) for configuration so you can override any of the default configuration under `conf/` using run time environment variables passed to `docker run`

  See [docker run](https://docs.docker.com/engine/reference/run/) for more details

### Docker based deploy (advanced)

In an environment with more heavy usage you almost certainly want a reverse proxy.
