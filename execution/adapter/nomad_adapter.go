package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/datagovsg/flotilla-os/config"
	"github.com/datagovsg/flotilla-os/state"

	nomad "github.com/hashicorp/nomad/api"
	jobspec "github.com/hashicorp/nomad/jobspec"
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
// Assume only docker tasks used in this nomad job, and only 1 task per nomad job
// Assume only 1 allocation per nomad job
//
func (a *nomadAdapter) AdaptTask(job nomad.Job) state.Run {
	val, _ := json.Marshal(job)
	submitTime := gjson.GetBytes(val, "SubmitTime").Time()
	exitCode := gjson.GetBytes(val, "TaskGroups.0.Tasks.0.KillSignal").Int()
	status := gjson.GetBytes(val, "Status").Str

	jobID := job.ID
	// should deal with error fetching allocations here
	resp, _, _ := a.nc.Jobs().Allocations(*jobID, true, nil)
	allocation := resp[0]
	desiredStatus := allocation.DesiredStatus
	// env := gjson.GetBytes(val, "TaskGroups.0.tasks.0.Env")

	run := state.Run{
		Status:    status,
		StartedAt: &submitTime,
		ExitCode:  &exitCode,
		// use allocation of job as a proxy for DesiredStatus
		// Env:       &env,
	}

	if desiredStatus != "" && desiredStatus != nomad.AllocDesiredStatusRun {
		run.Status = state.StatusStopped
	}

	if a.needsRetried(run, job) {
		run.Status = state.StatusNeedsRetry
	}

	return run
}

func (a *nomadAdapter) needsRetried(run state.Run, job nomad.Job) bool {
	//
	// This is a -strong- indication of abnormal exit, not internal to the run
	//
	// if run.Status == state.StatusStopped && run.ExitCode == nil {
	// 	containerReason := "?"
	// 	if len(task.Containers) == 1 {
	// 		container := task.Containers[0]
	// 		if container != nil && container.Reason != nil {
	// 			containerReason = *container.Reason
	// 		}
	// 	}

	// 	for _, retriable := range a.retriable {
	// 		// Container's stopped reason contains a retriable error
	// 		if strings.Contains(containerReason, retriable) {
	// 			return true
	// 		}
	// 	}
	// }
	return false
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

	jpath := definition.Template
	var jobfile io.Reader

	// Get the pwd and open the jobfile
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Nooo 1")
		fmt.Println(err)
		return NomadRunInput{}
	}
	file, err := os.Open(pwd + "/jobspecs/" + jpath)
	defer file.Close()
	if err != nil {
		fmt.Println("Nooo 2")
		fmt.Println(err)
		return NomadRunInput{}
	}
	jobfile = file

	// parse the jobfile from HCL to a nomad.Job struct
	job, err := jobspec.Parse(jobfile)
	if err != nil {
		fmt.Println("Nooo 3")
		fmt.Println(err)
		return NomadRunInput{}
	}

	// TODO inject the image from definition into thee job

	// inject the envvar from definition into the job
	task := job.TaskGroups[0].Tasks[0]
	envList := run.Env
	taskMap := map[string]string{}
	for _, envVar := range *envList {
		taskMap[envVar.Name] = envVar.Value
	}
	task.Env = taskMap

	rtn := NomadRunInput{
		Job: job,
	}

	return rtn
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
