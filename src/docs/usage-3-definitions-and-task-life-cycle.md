---
path: "/usage/definitions-and-task-life-cycle"
title: "Definitions and Task Life Cycle"
group: "usage"
index: "3"
---

### Definitions
| Name | Definition |
| ---- | ---------- |
| `task` | A definition of a task that can be executed to create a `run` |
| `run` | An instance of a task |

### Task Life Cycle

When executed, a task's run goes through several transitions

1. `QUEUED` - this is the first phase of a run and means the run is currently queued and waiting to be allocated to a cluster
2. `PENDING` - every `worker.submit_interval` (defined in the config) the submit worker pulls from the queues and submits them for execution. At this point, if the cluster associated with the run has resources, the run gets allocated to the cluster and transitions to the `PENDING` status. For the default execution engine this stage encapsulates the process of pulling the docker image and starting the container. It can take several minutes depending on whether the image is cached and how large the image is.
3. `RUNNING` - Once the run starts on a particular execution host it transitions to this stage. At this point logs should become available.
4. `STOPPED` - A run enters this stage when it finishes execution. This can mean it either succeeded or failed depending on the existence of an `exit_code` and the value of that exit code.
5. `NEEDS_RETRY` - on occassion, due to host level characteristics (full disk, too many open files, timeouts pulling image, etc) the run exits with a null exit code without ever being executed. In this case the reason is analyzed to determine if the run is retriable. If it is, the task transitions to this status and is allocated to the appropriate execution queue again, and will repeat the lifecycle.

#### Normal Lifecycle

`QUEUED` --> `PENDING` --> `RUNNING` --> `STOPPED`

#### Retry Lifecycle

... --> `PENDING` --> `STOPPED` --> `NEEDS_RETRY` --> `QUEUED` --> ...

