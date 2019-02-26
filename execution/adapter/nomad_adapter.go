package adapter

import (
	"fmt"
	"github.com/datagovsg/flotilla-os/config"
	"github.com/datagovsg/flotilla-os/state"
	"github.com/hashicorp/nomad/api"
	"strings"
)

//
// NomadAdapter translates back and forth from Nomad api objects to our representation
//
type NomadAdapter interface {
	// AdaptTask converts from an nomad job to a generic run
	AdaptTask(task api.Job) state.Run
	// AdaptTaskDef converts from an nomad jobspec to a generic definition
	AdaptTaskDef(task api.Job) state.Definition

	// AdaptDefinition translates from definition to a nomad Job struct (required for executing a job)
	AdaptDefinition(definition state.Definition) api.Job
}

type NomadClient interface {
	// ParseHCL is used to convert the HCL repesentation of a Job to JSON server side.
	ParseHCL(jobHCL string, canonicalize bool) (*api.Job, error)
	Register(job *api.Job, q *api.WriteOptions) (*api.JobRegisterResponse, *api.WriteMeta, error)
	List(q *api.QueryOptions) ([]*api.JobListStub, *api.QueryMeta, error)
}

type nomadAdapter struct {
	nc        NomadClient
	conf      config.Config
	retriable []string
}

//
// NewNomadAdapter configures and returns a nomad adapter for translating
// from Nomad api specific objects to our representation
//
func NewNomadAdapter(conf config.Config, nc NomadClient) (NomadAdapter, error) {
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
func (a *nomadAdapter) AdaptTask(task api.Job) state.Run {
	run := state.Run{}

	// describe job (allocations?)
	// type Run struct {
	// 	TaskArn         string     `json:"task_arn"`
	// 	RunID           string     `json:"run_id"`
	// 	DefinitionID    string     `json:"definition_id"`
	// 	Alias           string     `json:"alias"`
	// 	Image           string     `json:"image"`
	// 	ClusterName     string     `json:"cluster"`
	// 	ExitCode        *int64     `json:"exit_code,omitempty"`
	// 	Status          string     `json:"status"`
	// 	StartedAt       *time.Time `json:"started_at,omitempty"`
	// 	FinishedAt      *time.Time `json:"finished_at,omitempty"`
	// 	InstanceID      string     `json:"-"`
	// 	InstanceDNSName string     `json:"-"`
	// 	GroupName       string     `json:"group_name"`
	// 	User            string     `json:"user,omitempty"`
	// 	TaskType        string     `json:"-"`
	// 	Env             *EnvList   `json:"env,omitempty"`
	// }

	// needsRetried(run, task)

	return run
}

//
// AdaptTaskDef converts from an nomad jobspec to a generic definition
//
func (a *nomadAdapter) AdaptTaskDef(taskDef ecs.TaskDefinition) state.Definition {

	// type Definition struct {
	// 	Alias         string     `json:"alias"`          // User given name when defining task in Flotilla
	// 	Memory        *int64     `json:"memory"`         // Memory of the Docker image
	// 	User          string     `json:"user,omitempty"` // Name of the user running in Docker image
	// 	Arn           string     `json:"arn,omitempty"`
	// 	DefinitionID  string     `json:"definition_id"` // `name` in Nomad
	// 	Image         string     `json:"image"`
	// 	GroupName     string     `json:"group_name"`     // `dockerLabels` in ECS
	// 	ContainerName string     `json:"container_name"` // same as definition_id for ECS
	// 	Command       string     `json:"command,omitempty"`
	// 	TaskType      string     `json:"-"`
	// 	Env           *EnvList   `json:"env"`
	// 	Ports         *PortsList `json:"ports,omitempty"`
	// 	Tags          *Tags      `json:"tags,omitempty"`
	// }

}

//
// AdaptDefinition translates from definition to the ecs arguments for registering a nomad job
//
func (a *nomadAdapter) AdaptDefinition(definition state.Definition) api.Job {

	// Author the job file according to the job specification
	// Plan and review the changes with a Nomad server
	// Submit the job file to a Nomad server
	// (Optional) Review job status and logs
}
