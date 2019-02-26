package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/datagovsg/flotilla-os/config"
	"github.com/datagovsg/flotilla-os/execution/adapter"
	"github.com/datagovsg/flotilla-os/queue"
	"github.com/datagovsg/flotilla-os/state"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/command"
	"github.com/hashicorp/nomad/helper"
	"github.com/pkg/errors"
)

//
// NomadExecutionEngine submits runs to ecs
//
type NomadExecutionEngine struct {
	nomadClient api.Client
	sqsClient   sqsClient
	adapter     adapter.NomadAdapter
	qm          queue.Manager
	statusQurl  string
}

type sqsClient interface {
	GetQueueAttributes(input *sqs.GetQueueAttributesInput) (*sqs.GetQueueAttributesOutput, error)
	SetQueueAttributes(input *sqs.SetQueueAttributesInput) (*sqs.SetQueueAttributesOutput, error)
}

//
// Initialize configures the NomadExecutionEngine and initializes internal clients
//
func (ee *NomadExecutionEngine) Initialize(conf config.Config) error {
	// Get the HTTP client
	nomadClient, err := c.Meta.Client()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}
}

//
// PollStatus pops status updates from the status queue using the QueueManager
//
func (ee *NomadExecutionEngine) PollStatus() (RunReceipt, error) {
	var (
		receipt RunReceipt
		update  ecsUpdate
		err     error
	)

	rawReceipt, err := ee.qm.ReceiveStatus(ee.statusQurl)
	if err != nil {
		return receipt, errors.Wrapf(err, "problem getting status from [%s]", ee.statusQurl)
	}

	//
	// If we receive an update that is empty, don't try to deserialize it
	//
	if rawReceipt.StatusUpdate != nil {
		err = json.Unmarshal([]byte(*rawReceipt.StatusUpdate), &update)
		if err != nil {
			return receipt, errors.Wrapf(err, "unable to parse status update with json [%s]", *rawReceipt.StatusUpdate)
		}
		adapted := ee.adapter.AdaptTask(update.Detail)
		receipt.Run = &adapted
	}

	receipt.Done = rawReceipt.Done
	return receipt, nil
}

//
// PollRuns receives -at most- one run per queue that is pending execution
//
func (ee *NomadExecutionEngine) PollRuns() ([]RunReceipt, error) {
	queues, err := ee.qm.List()
	if err != nil {
		return nil, errors.Wrap(err, "problem listing queues to poll")
	}

	var runs []RunReceipt
	for _, qurl := range queues {
		//
		// Get new queued Run
		//
		runReceipt, err := ee.qm.ReceiveRun(qurl)

		if err != nil {
			return runs, errors.Wrapf(err, "problem receiving run from queue url [%s]", qurl)
		}

		if runReceipt.Run == nil {
			continue
		}

		runs = append(runs, RunReceipt{runReceipt})
	}
	return runs, nil
}

//
// Enqueue pushes a run onto the queue using the QueueManager
//
func (ee *NomadExecutionEngine) Enqueue(run state.Run) error {
	// Get qurl
	qurl, err := ee.qm.QurlFor(run.ClusterName, true)
	if err != nil {
		return errors.Wrapf(err, "problem getting queue url for [%s]", run.ClusterName)
	}

	// Queue run
	if err = ee.qm.Enqueue(qurl, run); err != nil {
		return errors.Wrapf(err, "problem enqueing run [%s] to queue [%s]", run.RunID, qurl)
	}
	return nil
}

//
// Execute takes a pre-configured run and definition and submits them for execution
// to Nomad
//
func (ee *NomadExecutionEngine) Execute(definition state.Definition, run state.Run) (state.Run, bool, error) {
	// Definition to Job
	var executed state.Run
	job := ee.adapter.AdaptRun(definition, run)

	// Parse the Vault token
	if vaultToken == "" {
		// Check the environment variable
		vaultToken = os.Getenv("VAULT_TOKEN")
	}

	if vaultToken != "" {
		job.VaultToken = helper.StringToPtr(vaultToken)
	}

	if output {
		req := api.RegisterJobRequest{Job: job}
		buf, err := json.MarshalIndent(req, "", "    ")
		if err != nil {
			fmt.Println(fmt.Sprintf("Error converting job: %s", err))
			return 1
		}

		fmt.Println(string(buf))
		return 0
	}

	// Parse the check-index
	checkIndex, enforce, err := parseCheckIndex(checkIndexStr)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error parsing check-index value %q: %v", checkIndexStr, err))
		return 1
	}

	// Set the register options
	opts := &api.RegisterOptions{}
	if enforce {
		opts.EnforceIndex = true
		opts.ModifyIndex = checkIndex
	}
	if override {
		opts.PolicyOverride = true
	}

	// Submit the job
	resp, _, err := client.Jobs().RegisterOpts(job, opts, nil)
	if err != nil {
		if strings.Contains(err.Error(), api.RegisterEnforceIndexErrPrefix) {
			// Format the error specially if the error is due to index
			// enforcement
			matches := enforceIndexRegex.FindStringSubmatch(err.Error())
			if len(matches) == 2 {
				fmt.Println(matches[1]) // The matched group
				fmt.Println("Job not updated")
				return 1
			}
		}

		fmt.Println(fmt.Sprintf("Error submitting job: %s", err))
		return 1
	}

	if len(result.Failures) != 0 {
		msg := make([]string, len(result.Failures))
		for i, failure := range result.Failures {
			msg[i] = *failure.Reason
		}
		// Retry these, they are very rare;
		// our upfront validation catches the obvious image and cluster resources
		// IMPORTANT - log these messages
		return executed, true, errors.Errorf("ERRORS: %s", strings.Join(msg, "\n"))
	}

	return ee.adapter.AdaptTask(*result.Tasks[0]), false, nil
}

//
// Terminate takes a valid run and stops it
//
func (ee *NomadExecutionEngine) Terminate(run state.Run) error {
	// Check if the job exists
	jobs, _, err := client.Jobs().PrefixList(jobID)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error deregistering job: %s", err))
		return 1
	}
	if len(jobs) == 0 {
		fmt.Println(fmt.Sprintf("No job(s) with prefix or id %q found", jobID))
		return 1
	}
	if len(jobs) > 1 && strings.TrimSpace(jobID) != jobs[0].ID {
		fmt.Println(fmt.Sprintf("Prefix matched multiple jobs\n\n%s", createStatusListOutput(jobs)))
		return 1
	}
	// Prefix lookup matched a single job
	job, _, err := client.Jobs().Info(jobs[0].ID, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error deregistering job: %s", err))
		return 1
	}

	// Invoke the stop
	evalID, _, err := client.Jobs().Deregister(*job.ID, purge, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error deregistering job: %s", err))
		return 1
	}

	return err

}

//
// The concept of defining a pre-defined task and deregistering a defined task does not apply to Nomad
// These 2 functions are functionally used and are preserved to match the Engine interface
//
func (ee *NomadExecutionEngine) Define(definition state.Definition) (state.Definition, error) {
	return definition, err
}

func (ee *NomadExecutionEngine) Deregister(definition state.Definition) error {
	return err
}
