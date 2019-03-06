package engine

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/datagovsg/flotilla-os/config"
	"github.com/datagovsg/flotilla-os/queue"
	"github.com/datagovsg/flotilla-os/state"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"
)

//
// Engine defines the execution engine interface.
//
type Engine interface {
	Initialize(conf config.Config) error
	// v0
	Execute(definition state.Definition, run state.Run) (state.Run, bool, error)

	// v1 - once runs contain a copy of relevant definition info
	// Execute(run state.Run) error

	Define(definition state.Definition) (state.Definition, error)

	Deregister(definition state.Definition) error

	Terminate(run state.Run) error

	Enqueue(run state.Run) error

	PollRuns() ([]RunReceipt, error)

	PollStatus() (RunReceipt, error)
}

type sqsClient interface {
	GetQueueAttributes(input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error)
	SetQueueAttributes(input *sqs.SetQueueAttributesInput) (*sqs.SetQueueAttributesOutput, error)
}

type cloudwatchServiceClient interface {
	PutRule(input *cloudwatchevents.PutRuleInput) (*cloudwatchevents.PutRuleOutput, error)
	PutTargets(input *cloudwatchevents.PutTargetsInput) (*cloudwatchevents.PutTargetsOutput, error)
	ListRuleNamesByTarget(input *cloudwatchevents.ListRuleNamesByTargetInput) (*cloudwatchevents.ListRuleNamesByTargetOutput, error)
}

type nomadUpdate struct {
	Detail nomad.Job `json:"detail"`
}

type ecsUpdate struct {
	Detail ecs.Task `json:"detail"`
}

type RunReceipt struct {
	queue.RunReceipt
}

//
// NewExecutionEngine initializes and returns a new Engine
//
func NewExecutionEngine(conf config.Config, qm queue.Manager) (Engine, error) {
	name := "nomad"
	if conf.IsSet("execution_engine") {
		name = conf.GetString("execution_engine")
	}

	switch name {
	case "ecs":
		eng := &ECSExecutionEngine{qm: qm}
		if err := eng.Initialize(conf); err != nil {
			return nil, errors.Wrap(err, "problem initializing ECSExecutionEngine")
		}
		return eng, nil
	case "nomad":
		eng := &NomadExecutionEngine{qm: qm}
		if err := eng.Initialize(conf); err != nil {
			return nil, errors.Wrap(err, "problem initializing NomadExecutionEngine")
		}
		return eng, nil
	default:
		return nil, fmt.Errorf("no Engine named [%s] was found", name)
	}
}
