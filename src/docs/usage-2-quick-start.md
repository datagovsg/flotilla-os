---
path: "/usage/quick-start"
title: "Quick Start"
group: "usage"
index: "2"
---

### Minimal Assumptions

Before we can do _anything_ there's some *prerequistes* that must be met.

1. Flotilla by default uses AWS. You must have an AWS account and the credentials available to you in a way that standard AWS tools can access. That is, the standard credential provider chain. This means one of:
	1. Environment variables
	2. A shared credentials file
	3. IAM role
2. Flotilla uses AWS's Elastic Continer Service (ECS) as the execution backend. However, Flotilla does not manage ECS clusters. There must be at least one cluster defined in AWS's ECS service available to you and it must have at least one task node. Most typically this is the `default` cluster and examples will assume this going forward.

### Starting the service locally

You can run the service locally (which will still leverage AWS resources) using the [docker-compose](https://docs.docker.com/compose/) tool. From inside the repo run:

```
docker-compose up -d
```

You'll notice it builds the code in the repo and starts the flotilla service as well as the default postgres backend.

Verify the service is running by making a `GET` request with cURL (or navigating to in a web browser) the url `http://localhost:3000/api/v1/task`. A 200OK response means things are good!

> Note: The default configuration under `conf` and in the `docker-compose.yml` assume port 3000. You'll have to change it in both places if you don't want to use port 3000 locally.

### Using the UI

Flotilla has a simple, easy to use UI. Here's some example images for basic usage.

#### Define a task with the UI

The UI allows you to quickly create new tasks.

![Define Task](https://user-images.githubusercontent.com/166823/35579884-999e0b04-059c-11e8-8ec7-65a14f7ae647.png "Create New Task")


#### Launch a task with UI
You can run tasks you've created with the UI as well. Once you've ran a task the run will transition from `Queued` to `Pending` to `Running` before it finishes and shows `Success` or `Failed` (see [Task Life Cycle](#definitions-and-task-life-cycle)). Once a task is in the `Running` state the logs should be visible.



1. Launch

   ![Run Task](https://user-images.githubusercontent.com/166823/35579954-cc158d8c-059c-11e8-857c-0c0c7e8cd7b9.png "Run Task")

2. Queued --> Pending

   ![Queued Task](https://user-images.githubusercontent.com/166823/35579975-e1e3bb20-059c-11e8-87d5-5c78f8aa96f4.png "Queued Task")

   ![Pending Task](https://user-images.githubusercontent.com/166823/35579998-eff41368-059c-11e8-8fb7-d5c217998a5d.png "Pending Task")
3. View logs

   ![Running Task](https://user-images.githubusercontent.com/166823/35580026-038ae348-059d-11e8-95e4-f0150400a1a8.png "Running Task")

   ![Finished Task](https://user-images.githubusercontent.com/166823/35580037-1455ea10-059d-11e8-92da-dd1249dcf40d.png "Finished Task")


### Basic API Usage

#### Defining your first task
Before you can run a task you first need to define it. We'll use the example hello world task definition. Here's what that looks like:

> hello-world.json
>
```
{
  "alias": "hello-flotilla",
  "group_name": "examples",
  "image": "ubuntu:latest",
  "memory": 512,
  "env": [
    {
      "name": "USERNAME",
      "value": "_fill_me_in_"
    }
  ],
  "command": "echo \"hello ${USERNAME}\""
}
```

It's a simple task that runs in the default ubuntu image, prints your username to the logs, and exits.

> Note: While you can use non-public images and images in your own registries with flotilla, credentials for accessing those images must exist on the ECS hosts. This is outside the scope of this doc. See the AWS [documentation](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/private-auth.html).


Let's define it:


```
curl -XPOST localhost:3000/api/v1/task --data @examples/hello-world.json
```

You'll notice that if you visit the initial url again `http://localhost:3000/api/v1/task` the newly defined definition will be in the list.

#### Running your first task

This is the fun part. You'll make a `PUT` request to the execution endpoint for the task you just defined and specify any environment variables.

```
curl -XPUT localhost:3000/api/v1/task/alias/hello-flotilla/execute -d '{
  "cluster":"default",
  "env":[
    {"name":"USERNAME","value":"yourusername"}
  ],
  "run_tags":{"owner_id":"youruser"}
}'
```
> Note: `run_tags` is defined as a way for all runs to have a ownership injected for visibility and is *required*.

You'll get a response that contains a `run_id` field. You can check the status of your task at `http://localhost:3000/api/v1/history/<run_id>`

```
curl -XGET localhost:3000/api/v1/history/<run_id>

{
  "instance": {
    "dns_name": "<dns-host-of-task-node>",
    "instance_id": "<instance-id-of-task-node>"
  },
  "run_id": "<run_id>",
  "definition_id": "<definition_id>",
  "alias": "hello-flotilla",
  "image": "ubuntu:latest",
  "cluster": "default",
  "status": "PENDING",
  "env": [
    {
      "name": "FLOTILLA_RUN_OWNER_ID",
      "value": "youruser"
    },
    {
      "name": "FLOTILLA_SERVER_MODE",
      "value": "dev"
    },
    {
      "name": "FLOTILLA_RUN_ID",
      "value": "<run_id>"
    },
    {
      "name": "USERNAME",
      "value": "yourusername"
    }
  ]
}
```

and you can get the logs for your task at `http://localhost:3000/api/v1/<run_id>/logs`. You will not see any logs until your task is at least in the `RUNNING` state.

```
curl -XGET localhost:3000/api/v1/<run_id>/logs

{
  "last_seen":"<last_seen_token_used_for_paging>",
  "log":"+ set -e\n+ echo 'hello yourusername'\nhello yourusername"
}
```

