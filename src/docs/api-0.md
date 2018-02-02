---
path: "/api"
title: "API Documentation"
group: "api"
index: "0"
---
Flotilla
========
Flotilla server. Define your own integrations, clients, and workflows using the Flotilla API

**Version:** 1.0.0

**License:** [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

[Find out more about flotilla](https://github.com/stitchfix/flotilla-os)
### /v1/task/{definition_id}
---
##### ***PUT***
**Summary:** Update task definition

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of task definition to update | Yes | string |
| body | body | Definition updates | Yes | [Definition](#definition) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | successful operation | [Definition](#definition) |
| 400 | Malformed request |
| 404 | Definition not found |

##### ***GET***
**Summary:** Get task definition by id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of task definition to get | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | successful operation | [Definition](#definition) |
| 404 | Definition not found |

##### ***DELETE***
**Summary:** Delete task definition by id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of task definition to delete | Yes | string |

**Responses**

| Code | Description |
| ---- | ----------- |
| 200 | successful operation |
| 404 | Definition not found |

### /v1/task/alias/{alias}/execute
---
##### ***PUT***
**Summary:** Create new task run by task alias

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| alias | path | Alias to task definition to launch | Yes | string |
| body | body | Launch parameters | Yes | [LaunchRequestV2](#launchrequestv2) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [Run](#run) |
| 400 | Malformed request. Most typically this means an invalid image is used or the cluster cannot support running the job |
| 404 | Definition not found |

### /v1/task
---
##### ***POST***
**Summary:** Create a new task definition. Execute existing definitions with the execute endpoint

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| body | body | Definition object to create | Yes | [Definition](#definition) |

**Responses**

| Code | Description |
| ---- | ----------- |
| 200 | Definition created successfully |
| 400 | Malformed input |
| 409 | Another definition with the same alias exists |

##### ***GET***
**Summary:** List existing task definitions

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| limit | query | Limit number of results to this many | No | integer |
| offset | query | Page offset | No | integer |
| group_name | query | Filter to definitions with the specified group_name. Strict match. | No | string |
| alias | query | Filter to those definitions with the specified alias. Substring match. | No | string |
| sort_by | query | Field on task definitions to sort by | No | string |
| order | query | Sort order; (asc, desc) | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | successful operation | [DefinitionList](#definitionlist) |

### /v4/task/{definition_id}/execute
---
##### ***PUT***
**Summary:** Create new task run by definition id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of task to launch | Yes | string |
| body | body | Launch parameters | Yes | [LaunchRequestV4](#launchrequestv4) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [Run](#run) |
| 400 | Malformed request. Most typically this means an invalid image is used or the cluster cannot support running the job |
| 404 | Definition not found |

### /groups
---
##### ***GET***
**Summary:** Get list of group names in use

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| limit | query | Limit results to this many | No | integer |
| offset | query | Page offset | No | integer |
| name | query | Search for group names with this value | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [GroupsList](#groupslist) |

### /v1/task/alias/{alias}
---
##### ***GET***
**Summary:** Get task definition by alias

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| alias | path | Alias to task definition to get | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | successful operation | [Definition](#definition) |
| 404 | Definition not found |

### /v1/task/{definition_id}/history
---
##### ***GET***
**Summary:** List all runs for a given task definition id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of runs to fetch | Yes | string |
| limit | query | Limit number of results to this many | No | integer |
| offset | query | Page offset | No | integer |
| cluster_name | query | Filter to runs with the specified cluster_name. Strict match. | No | string |
| status | query | Filter to those definitions with the specified status. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [RunList](#runlist) |

### /v1/history
---
##### ***GET***
**Summary:** List all runs

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| limit | query | Limit number of results to this many | No | integer |
| offset | query | Page offset | No | integer |
| cluster_name | query | Filter to runs with the specified cluster_name. Strict match. | No | string |
| status | query | Filter to those definitions with the specified status. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [RunList](#runlist) |

### /v1/{run_id}/logs
---
##### ***GET***
**Summary:** Get logs for a run

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| run_id | path | Run id of run to get logs for | Yes | string |
| last_seen | query | Last seen token (from a previous call to this endpoint). Allows only fetching logs that come after the specified `last_seen` token. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [RunLogs](#runlogs) |
| 404 | No run with that run_id found |

### /v1/task/{definition_id}/execute
---
##### ***PUT***
**Summary:** Create new task run by definition id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of task definition to launch | Yes | string |
| body | body | Launch parameters | Yes | [LaunchRequestV1](#launchrequestv1) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | successful operation | [Run](#run) |
| 404 | Definition not found |

### /v2/task/{definition_id}/execute
---
##### ***PUT***
**Summary:** Create new task run by definition id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of task to launch | Yes | string |
| body | body | Launch parameters | Yes | [LaunchRequestV2](#launchrequestv2) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [Run](#run) |
| 400 | Malformed request. Most typically this means an invalid image is used or the cluster cannot support running the job |
| 404 | Definition not found |

### /v1/task/{definition_id}/history/{run_id}
---
##### ***DELETE***
**Summary:** Terminate a run

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| definition_id | path | Definition id of run to terminate | Yes | string |
| run_id | path | Run id of run to terminate | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [Run](#run) |
| 404 | No run found |

### /v1/history/{run_id}
---
##### ***GET***
**Summary:** Get run by run id

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| run_id | path | Run id of run to get | Yes | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [Run](#run) |
| 404 | No run found |

### /tags
---
##### ***GET***
**Summary:** Get list of tags in use

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| limit | query | Limit results to this many | No | integer |
| offset | query | Page offset | No | integer |
| name | query | Search for tags with this value | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [TagsList](#tagslist) |

### /clusters
---
##### ***GET***
**Summary:** Get list of clusters

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Successful operation | [ClustersList](#clusterslist) |

### Models
---

### Definition  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| tags | [ string ] |  | No |
| image | string |  | No |
| group_name | string |  | No |
| definition_id | string |  | No |
| alias | string |  | No |
| command | string |  | No |
| env | [ [EnvVar](#envvar) ] |  | No |
| memory | integer | memory in units of megabytes | No |

### Run  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| status | string |  | No |
| cluster | string |  | No |
| finished_at | dateTime |  | No |
| image | string |  | No |
| run_id | string |  | No |
| exit_code | integer |  | No |
| group_name | string |  | No |
| definition_id | string |  | No |
| instance | object |  | No |
| alias | string |  | No |
| env | [ [EnvVar](#envvar) ] |  | No |
| started_at | dateTime |  | No |

### RunList  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| status | string |  | No |
| history | [ [Run](#run) ] |  | No |
| limit | integer |  | No |
| offset | integer |  | No |
| total | integer |  | No |
| order | string |  | No |
| sort_by | string |  | No |

### LaunchRequestV1  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cluster | string |  | No |
| env | [ [EnvVar](#envvar) ] |  | No |

### DefinitionList  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| limit | integer |  | No |
| offset | integer |  | No |
| definitions | [ [Definition](#definition) ] |  | No |
| total | integer |  | No |
| order | string |  | No |
| sort_by | string |  | No |

### GroupsList  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| total | integer |  | No |
| limit | integer |  | No |
| groups | [ string ] |  | No |
| offset | integer |  | No |

### LaunchRequestV2  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cluster | string |  | No |
| run_tags | object |  | No |
| env | [ [EnvVar](#envvar) ] |  | No |

### ClustersList  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| clusters | [ string ] |  | No |

### RunLogs  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| log | string |  | No |
| last_seen | string |  | No |

### EnvVar  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| value | string |  | No |

### TagsList  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| total | integer |  | No |
| limit | integer |  | No |
| offset | integer |  | No |
| tags | [ string ] |  | No |

### LaunchRequestV4  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cluster | string |  | No |
| run_tags | object |  | No |
| env | [ [EnvVar](#envvar) ] |  | No |
