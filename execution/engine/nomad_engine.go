package engine

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/sqs"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"

	"github.com/datagovsg/flotilla-os/config"
	"github.com/datagovsg/flotilla-os/execution/adapter"
	"github.com/datagovsg/flotilla-os/queue"
	"github.com/datagovsg/flotilla-os/state"
)

//
// NomadExecutionEngine submits runs to ecs
//
type NomadExecutionEngine struct {
	nomadClient nomad.Client
	sqsClient   sqsClient
	cwClient    cloudwatchServiceClient
	adapter     adapter.NomadAdapter
	qm          queue.Manager
	statusQurl  string
}

// NewNomadClient is used to create a new client to interact with Nomad. The
// client implements the NomadClient interface.
func NewNomadClient() (nomad.Client, error) {

	config := nomad.DefaultConfig()

	c, err := nomad.NewClient(config)
	if err != nil {
		return nomad.Client{}, fmt.Errorf("problem initializing NomadClient")
	}

	return *c, nil
}

//
// Initialize configures the NomadExecutionEngine and initializes internal clients
//
func (ne *NomadExecutionEngine) Initialize(conf config.Config) error {
	// Get the HTTP client
	nomadClient, err := NewNomadClient()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error initializing client: %s", err))
	}
	ne.nomadClient = nomadClient

	flotillaMode := conf.GetString("flotilla_mode")

	var adpt adapter.NomadAdapter

	//
	// When mode is not test, setup and initialize all aws clients
	// - this isn't ideal; is there another way?
	//
	if flotillaMode != "test" {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(conf.GetString("aws_default_region"))}))

		ne.cwClient = cloudwatchevents.New(sess)
		ne.sqsClient = sqs.New(sess)
		adpt, err = adapter.NewNomadAdapter(conf, nomadClient)
		if err != nil {
			return errors.Wrap(err, "problem initializing NomadAdapter")
		}
	}

	ne.adapter = adpt

	if ne.qm == nil {
		return errors.Errorf("no queue.Manager implementation; ECSExecutionEngine needs a queue.Manager")
	}

	//
	// Calling QurlFor creates the status queue if it does not exist
	// - this is necessary for the next step of creating an nomad
	//   job update rule in cloudwatch which routes task updates
	//   to the status queue
	//
	statusQueue := conf.GetString("queue.status")
	ne.statusQurl, err = ne.qm.QurlFor(statusQueue, false)
	if err != nil {
		return errors.Wrapf(err, "problem getting queue url for status queue with name [%s]", statusQueue)
	}

	statusRule := conf.GetString("queue.status_rule")
	return ne.createOrUpdateEventRule(statusRule, statusQueue)
}

//
// This function is used to set up the CloudWatch routing rule that will redirect
// events with {source: ["nomad.script"]} to a SQS queue
//
func (ne *NomadExecutionEngine) createOrUpdateEventRule(statusRule string, statusQueue string) error {
	createUpdate, err := ne.cwClient.PutRule(&cloudwatchevents.PutRuleInput{
		Description:  aws.String("Routes ecs task status events to flotilla status queues"),
		Name:         &statusRule,
		EventPattern: aws.String(`{"source":["nomad.script"],"detail-type":["Nomad Job State Change"]}`),
	})

	if err != nil {
		return errors.Wrap(err, "problem creating ecs status event routing rule")
	}

	// Route status events to the status queue
	targetArn, err := ne.getTargetArn(ne.statusQurl)
	if err != nil {
		return errors.Wrapf(err, "problem getting target arn for [%s]", ne.statusQurl)
	}

	names, err := ne.cwClient.ListRuleNamesByTarget(&cloudwatchevents.ListRuleNamesByTargetInput{
		TargetArn: &targetArn,
	})
	if err != nil {
		return errors.Wrapf(err, "problem listing rules for target [%s]", targetArn)
	}

	targetExists := len(names.RuleNames) > 0 && *names.RuleNames[0] == statusRule
	if !targetExists {
		res, err := ne.cwClient.PutTargets(&cloudwatchevents.PutTargetsInput{
			Rule: &statusRule,
			Targets: []*cloudwatchevents.Target{
				{
					Arn: &targetArn,
					Id:  &statusQueue,
				},
			},
		})

		if err != nil {
			return errors.Wrapf(
				err, "problem adding [%s] as queue target for status rule [%s]", targetArn, statusRule)
		}

		if *res.FailedEntryCount > 0 {
			failed := res.FailedEntries[0]
			return errors.Errorf("error adding routing rule for ecs status messages [%s]", *failed.ErrorMessage)
		}
	}
	// Finally, add permissions to target queue
	return ne.setTargetPermissions(*createUpdate.RuleArn, targetArn)
}

func (ne *NomadExecutionEngine) getTargetArn(qurl string) (string, error) {
	var arn string
	res, err := ne.sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: &qurl,
		AttributeNames: []*string{
			aws.String("QueueArn"),
		},
	})
	if err != nil {
		return arn, errors.Wrapf(err, "problem getting attribute QueueArn for sqs queue with url [%s]", qurl)
	}
	if res.Attributes["QueueArn"] != nil {
		return *res.Attributes["QueueArn"], nil
	}
	return arn, errors.Errorf("couldn't get queue arn")
}

func (ne *NomadExecutionEngine) setTargetPermissions(sourceArn string, targetArn string) error {
	policyDoc := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Id":"flotilla-task-status-updates-to-sqs",
		"Statement": [{
			"Sid": "flotilla-task-status-updates-to-sqs-sid",
			"Effect": "Allow",
			"Principal": {
				"AWS": "*"
			},
			"Action": "sqs:SendMessage",
			"Resource": "%s",
			"Condition": {
				"ArnEquals": {
					"aws:SourceArn": "%s"
				}
			}
		}]
	}`, targetArn, sourceArn)

	// Check first
	res, err := ne.sqsClient.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: &ne.statusQurl,
	})
	if err != nil {
		return errors.Wrapf(err, "problem getting queue attributes for sqs queue [%s]", ne.statusQurl)
	}

	if policy, ok := res.Attributes["Policy"]; ok && *policy == policyDoc {
		return nil
	}

	if _, err := ne.sqsClient.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		Attributes: map[string]*string{
			"Policy": &policyDoc,
		},
		QueueUrl: &ne.statusQurl,
	}); err != nil {
		return errors.Wrapf(
			err, "problem setting permissions allowing [%s] to send events to [%s]", sourceArn, targetArn)
	}
	return nil
}

//
// PollStatus pops status updates from the status queue using the QueueManager
// Called by /workers/status_worker.go
//
func (ne *NomadExecutionEngine) PollStatus() (RunReceipt, error) {
	var (
		receipt RunReceipt
		update  nomadUpdate
		err     error
	)

	rawReceipt, err := ne.qm.ReceiveStatus(ne.statusQurl)
	if err != nil {
		return receipt, errors.Wrapf(err, "problem getting status from [%s]", ne.statusQurl)
	}

	//
	// If we receive an update that is empty, don't try to deserialize it
	//
	if rawReceipt.StatusUpdate != nil {
		err = json.Unmarshal([]byte(*rawReceipt.StatusUpdate), &update)
		if err != nil {
			return receipt, errors.Wrapf(err, "unable to parse status update with json [%s]", *rawReceipt.StatusUpdate)
		}
		adapted := ne.adapter.AdaptTask(update.Detail)
		receipt.Run = &adapted
	}

	receipt.Done = rawReceipt.Done
	return receipt, nil
}

//
// PollRuns receives -at most- one run per queue that is pending execution
// Called by /workers/submit_worker.go
//
func (ne *NomadExecutionEngine) PollRuns() ([]RunReceipt, error) {
	queues, err := ne.qm.List()
	if err != nil {
		return nil, errors.Wrap(err, "problem listing queues to poll")
	}

	var runs []RunReceipt
	for _, qurl := range queues {
		//
		// Get new queued Run
		//
		runReceipt, err := ne.qm.ReceiveRun(qurl)

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
// Called by /services/execution.go
//
func (ne *NomadExecutionEngine) Enqueue(run state.Run) error {
	// Get qurl
	qurl, err := ne.qm.QurlFor(run.ClusterName, true)
	if err != nil {
		return errors.Wrapf(err, "problem getting queue url for [%s]", run.ClusterName)
	}

	// Queue run
	if err = ne.qm.Enqueue(qurl, run); err != nil {
		return errors.Wrapf(err, "problem enqueing run [%s] to queue [%s]", run.RunID, qurl)
	}
	return nil
}

//
// Execute takes a pre-configured run and definition and submits them for execution to Nomad
// Called by /workers/submit_worker.go
//
func (ne *NomadExecutionEngine) Execute(definition state.Definition, run state.Run) (state.Run, bool, error) {

	var executed state.Run

	nomadRunInput := ne.adapter.AdaptRun(definition, run)
	job := nomadRunInput.Job
	opts := &nomadRunInput.Options

	// enforceIndexRegex is a regular expression which extracts the enforcement error
	var enforceIndexRegex = regexp.MustCompile(`\((Enforcing job modify index.*)\)`)

	// Submit the job
	resp, _, err := ne.nomadClient.Jobs().RegisterOpts(job, opts, nil)
	if err != nil {
		retryable := false
		if strings.Contains(err.Error(), nomad.RegisterEnforceIndexErrPrefix) {
			// Format the error specially if the error is due to index
			// enforcement
			matches := enforceIndexRegex.FindStringSubmatch(err.Error())
			if len(matches) == 2 {
				fmt.Println(matches[1]) // The matched group
				fmt.Println("Job not updated")
				retryable = false
			}
		}
		fmt.Println(fmt.Sprintf("Error submitting job: %s", err))
		return executed, retryable, errors.Wrapf(err, "problem executing run [%s]", run.RunID)
	}

	fmt.Println("--- SUBMITTED JOB ---")

	// Print any warnings if there are any
	if resp.Warnings != "" {
		return executed, true, errors.Errorf("[bold][yellow]Job Warnings:\n%s[reset]\n", resp.Warnings)
	}

	// resp = JobRegisterResponse is used to respond to a job registration
	// type JobRegisterResponse struct {
	// 	EvalID          string
	// 	EvalCreateIndex uint64
	// 	JobModifyIndex  uint64
	// 	Warnings string
	// 	QueryMeta
	// }

	return state.Run{}, false, nil
}

//
// Terminate takes a valid run and stops it
// called by /services/execution.go
//
func (ne *NomadExecutionEngine) Terminate(run state.Run) error {
	// Check if the job exists
	jobID := run.JobName
	fmt.Println(jobID)
	jobs, _, err := ne.nomadClient.Jobs().PrefixList(jobID)
	if err != nil {
		return errors.Wrapf(err, "Error deregistering Nomad job: [%s]", run.RunID)
	}
	if len(jobs) == 0 {
		return errors.Wrapf(err, "No job(s) with prefix or id %q found", jobID)
	}
	if len(jobs) > 1 && strings.TrimSpace(jobID) != jobs[0].ID {
		return errors.Wrapf(err, "Prefix matched multiple jobs\n\n%s", jobs)
	}
	// Prefix lookup matched a single job
	job, _, err := ne.nomadClient.Jobs().Info(jobs[0].ID, nil)
	if err != nil {
		return errors.Wrapf(err, "Error deregistering Nomad job: [%s]", run.RunID)
	}
	// Invoke the stop
	// https://www.nomadproject.io/docs/commands/job/stop.html#stop-options
	purge := false
	evalID, _, err := ne.nomadClient.Jobs().Deregister(*job.ID, purge, nil)
	fmt.Println("--- evalID start ---")
	fmt.Println(evalID)
	fmt.Println("---  evalID end  ---")
	if err != nil {
		return errors.Wrapf(err, "Error deregistering Nomad job: [%s]", run.RunID)
	}

	return nil

}

func (ne *NomadExecutionEngine) ConstructRun(
	definition state.Definition, clusterName string, env *state.EnvList, ownerID string, reservedEnv map[string]func(run state.Run) string) (state.Run, error) {

	var (
		run state.Run
		err error
	)

	runID, err := state.NewRunID()
	if err != nil {
		return run, err
	}

	job := ne.adapter.AdaptDefinition(definition)

	run = state.Run{
		RunID:        runID,
		GroupName:    definition.GroupName,
		DefinitionID: definition.DefinitionID,
		Alias:        definition.Alias,
		JobName:      *job.ID,
		Status:       state.StatusQueued,
		User:         ownerID,
	}
	runEnv := ne.constructEnviron(run, env, reservedEnv)
	run.Env = &runEnv
	return run, nil
}

//
// add reserved environment variables to the run
//
func (ne *NomadExecutionEngine) constructEnviron(run state.Run, env *state.EnvList, reservedEnv map[string]func(run state.Run) string) state.EnvList {
	var runEnv []state.EnvVar
	var envVarNames []string
	i := 0
	if env != nil {
		for _, e := range *env {
			runEnv = append(runEnv, e)
			envVarNames = append(envVarNames, e.Name)
			i++
		}
	}
	for k, f := range reservedEnv {
		if !stringInSlice(k, envVarNames) {
			envVar := state.EnvVar{
				Name:  k,
				Value: f(run),
			}
			runEnv = append(runEnv, envVar)
			envVarNames = append(envVarNames, k)
			i++
		}
	}
	return state.EnvList(runEnv)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//
// The concept of defining a pre-defined task and deregistering a defined task does not apply to Nomad
// These 2 functions are functionally used and are preserved to match the Engine interface
//
func (ne *NomadExecutionEngine) Define(definition state.Definition) (state.Definition, error) {

	// TODO set this to be editable in UI rather than being hardcoded here
	jpath := "test.nomad"
	definition.Template = jpath
	// TODO set the envlist to non-nil in UI so it doesn't error out
	emptyEnvList := state.EnvList{}
	definition.Env = &emptyEnvList
	return definition, nil
}

func (ne *NomadExecutionEngine) Deregister(definition state.Definition) error {
	return nil
}
