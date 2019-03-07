package adapter

import (
	"encoding/json"

	"github.com/datagovsg/flotilla-os/config"
	"github.com/datagovsg/flotilla-os/state"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/tidwall/gjson"
)

//
// NomadAdapter translates back and forth from Nomad api objects to our representation
//
type NomadAdapter interface {
	// AdaptTask converts from an nomad job to a Run
	AdaptTask(task nomad.Job) state.Run
	// AdaptRun translates Definition and Run into the required arguments to run a Nomad job.
	AdaptRun(definition state.Definition, run state.Run) NomadRunInput

	// AdaptTaskDef converts from an nomad jobspec to a Definition
	// AdaptTaskDef(task api.Job) state.Definition
	// AdaptDefinition translates from Definition to a nomad Job struct
	// AdaptDefinition(definition state.Definition) api.Job
}

// type NomadClient interface {
// 	// ParseHCL is used to convert the HCL repesentation of a Job to JSON server side.
// 	ParseHCL(jobHCL string, canonicalize bool) (*api.Job, error)
// 	Register(job *api.Job, q *api.WriteOptions) (*api.JobRegisterResponse, *api.WriteMeta, error)
// 	List(q *api.QueryOptions) ([]*api.JobListStub, *api.QueryMeta, error)
// }

type NomadRunInput struct {
	Options nomad.RegisterOptions
	Job     *nomad.Job
}

type nomadAdapter struct {
	nc        nomad.Client
	conf      config.Config
	retriable []string
}

//
// NewNomadAdapter configures and returns a nomad adapter for translating
// from Nomad api specific objects to our representation
//
func NewNomadAdapter(conf config.Config, nc nomad.Client) (NomadAdapter, error) {
	adapter := nomadAdapter{
		conf: conf,
		nc:   nc,
		retriable: []string{
			"CannotCreateContainerError",
			"CannotStartContainerError",
			"CannotPullContainerError",
		},
	}
	return &adapter, nil
}

//
// AdaptTask converts from an nomad job to a generic run
//
func (a *nomadAdapter) AdaptTask(job nomad.Job) state.Run {
	val, _ := json.Marshal(job)
	submitTime := gjson.GetBytes(val, "SubmitTime").Time()
	exitCode := gjson.GetBytes(val, "TaskGroups.0.Tasks.0.KillSignal").Int()
	// env := gjson.GetBytes(val, "TaskGroups.0.tasks.0.Env")

	run := state.Run{
		// TaskArn:      "", // not applicable
		RunID:        gjson.GetBytes(val, "Name").Str,
		DefinitionID: gjson.GetBytes(val, "Name").Str,
		// Alias:        "", // not applicable
		Image: gjson.GetBytes(val, "TaskGroups.0.Tasks.0.Config.Image").Str,
		// ClusterName:  "", // not applicable
		ExitCode:  &exitCode,
		Status:    gjson.GetBytes(val, "Status").Str,
		StartedAt: &submitTime,
		// FinishedAt: gjson.GetBytes(val, ""), // not applicable
		// InstanceID: gjson.GetBytes(val, ""), // not applicable
		InstanceDNSName: gjson.GetBytes(val, "TaskGroups[0].Tasks[0].Config.dns_servers").Str,
		// GroupName: gjson.GetBytes(val, ""), // not applicable
		// User: gjson.GetBytes(val, TaskGroups.0.Tasks.0.User), // not applicable
		TaskType: gjson.GetBytes(val, "Type").Str,
		// Env:      &env,
	}

	// describe job (allocations?)
	// type Run struct {
	// 	TaskArn         string     `json:"task_arn"` // not applicable
	// 	RunID           string     `json:"id"` // allocation?
	// 	DefinitionID    string     `json:"id"` // job name
	// 	Alias           string     `json:"alias"` // user given name
	// 	Image           string     `json:"taskgroups[0].tasks[0].config.image"` // docker container name
	// 	ClusterName     string     `json:"cluster"` // not applicable
	// 	ExitCode        *int64     `json:"taskgroups[0].tasks[0].killsignal"` // kill signal?
	// 	Status          string     `json:"status"` // enum("running", "complete", "failed", "lost", "queued", "starting")
	// 	StartedAt       *time.Time `json:"submittime"` // submittime?
	// 	FinishedAt      *time.Time `json:"finished_at"` // not sure if applicable
	// 	InstanceID      string     `json:"-"` // allocation?
	// 	InstanceDNSName string     `json:"taskgroups[0].tasks[0].config.dns_servers"` // dns_servers?
	// 	GroupName       string     `json:"group_name"` // not applicable
	// 	User            string     `json:"taskgroups[0].tasks[0].user"` // not applicable
	// 	TaskType        string     `json:"type"` // enum("service", "batch", "parameterized")
	// 	Env             *EnvList   `json:"taskgroups[0].tasks[0].env"`
	// }

	// needsRetried(run, task)

	return run
}

func (a *nomadAdapter) AdaptRun(definition state.Definition, run state.Run) NomadRunInput {

	// Set the register options
	// if enforce {
	// 	opts.EnforceIndex = true
	// 	opts.ModifyIndex = checkIndex
	// }
	// if override {
	// 	opts.PolicyOverride = true
	// }
	return NomadRunInput{}
}

//
// AdaptTaskDef converts from an nomad jobspec to a generic definition
//
// func (a *nomadAdapter) AdaptTaskDef(taskDef ecs.TaskDefinition) state.Definition {
// 	definition := state.Definition{}

// 	// type Definition struct {
// 	// 	Alias         string     `json:"alias"`          // User given name when defining task in Flotilla
// 	// 	Memory        *int64     `json:"memory"`         // Memory of the Docker image
// 	// 	User          string     `json:"user,omitempty"` // Name of the user running in Docker image
// 	// 	Arn           string     `json:"arn,omitempty"`
// 	// 	DefinitionID  string     `json:"definition_id"` // `name` in Nomad
// 	// 	Image         string     `json:"image"`
// 	// 	GroupName     string     `json:"group_name"`     // `dockerLabels` in ECS
// 	// 	ContainerName string     `json:"container_name"` // same as definition_id for ECS
// 	// 	Command       string     `json:"command,omitempty"`
// 	// 	TaskType      string     `json:"-"`
// 	// 	Env           *EnvList   `json:"env"`
// 	// 	Ports         *PortsList `json:"ports,omitempty"`
// 	// 	Tags          *Tags      `json:"tags,omitempty"`
// 	// }

// 	return definition
// }

//
// AdaptDefinition translates from definition to the ecs arguments for registering a nomad job
//
// func (a *nomadAdapter) AdaptDefinition(definition state.Definition) api.Job {
// 	job := api.Job{}

// 	// Author the job file according to the job specification
// 	// Plan and review the changes with a Nomad server
// 	// Submit the job file to a Nomad server
// 	// (Optional) Review job status and logs

// 	return job
// }
