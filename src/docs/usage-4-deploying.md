---
path: "/usage/deploying"
title: "Deploying"
group: "usage"
index: "4"
---

In a production deployment you'll want multiple instances of the flotilla service running and postgres running elsewhere (eg. Amazon RDS). In this case the most salient detail configuration detail is the `DATABASE_URL`.

### Docker based deploy

The simplest way to deploy for very light usage is to avoid a reverse proxy and deploy directly with docker.

1. Build and tag an image for flotilla using the `Dockerfile` provided in this repo:

	```
	docker build -t <your repo name>/flotilla:<version tag>
	```
2. Run this image wherever you deploy your services:

	```
	docker run -e DATABASE_URL=<your db url> -e FLOTILLA_MODE=prod -p 3000:3000 ...<other standard docker run args>
	```

	> Notes:
	> -----
	> * Flotilla uses [viper](https://github.com/spf13/viper) for configuration so you can override any of the default configuration under `conf/` using run time environment variables passed to `docker run`
	> * In most realistic deploys you'll likely want to configure a reverse proxy to sit in front of the flotilla container. See the docs [here](https://hub.docker.com/_/nginx/)


	See [docker run](https://docs.docker.com/engine/reference/run/) for more details

### Configuration In Detail

The variables in `conf/config.yml` are sensible defaults. Most should be left alone unless you're developing flotilla itself. However, there are a few you may want to change in a production environment.

| Variable Name | Description |
| ------------- | ----------- |
| `worker.retry_interval` | Run frequency of the retry worker |
| `worker.submit_interval` | Poll frequency of the submit worker |
| `worker.status_interval` | Poll frequency of the status update worker |
| `http.server.read_timeout_seconds` | Sets read timeout in seconds for the http server |
| `http.server.write_timeout_seconds` | Sets the write timeout in seconds for the http server |
| `http.server.listen_address` | The port for the http server to listen on |
| `owner_id_var` | Which environment variable containing ownership information to inject into the runtime of jobs |
| `enabled_workers` | This variable is a list of the workers that run. Use this to control what workers run when using a multi-container deployment strategy. Valid list items include (`retry`, `submit`, and `status`) |
| `log.namespace` | For the default ECS execution engine setup this is the `log-group` to use |
| `log.retention_days` | For the default ECS execution engine this is the number of days to retain logs |
| `log.driver.options.*` | For the default ECS execution engine these map to the `awslogs` driver options [here](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html) |
| `queue.namespace` | For the default ECS execution engine this is the prefix used for SQS to determine which queues to pull job launch messages from |
| `queue.retention_seconds` | For the default ECS execution engine this configures how long a message will stay in an SQS queue without being consumed |
| `queue.process_time` | For the default ECS execution engine configures the length of time allowed to process a job launch message |
| `queue.status` | For the default ECS execution engine this configures which SQS queue to route ECS cluster status updates to |
| `queue.status_rule` | For the default ECS execution engine this configures the name of the rule for routing ECS cluster status updates |



