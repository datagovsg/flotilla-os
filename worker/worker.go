package worker

import (
	"fmt"
	"github.com/stitchfix/flotilla-os/config"
	"github.com/stitchfix/flotilla-os/execution/engine"
	flotillaLog "github.com/stitchfix/flotilla-os/log"
	"github.com/stitchfix/flotilla-os/queue"
	"github.com/stitchfix/flotilla-os/state"
)

//
// Worker defines a background worker process
//
type Worker interface {
	Run()
}

func NewWorker(
	workerType string,
	log flotillaLog.Logger,
	conf config.Config,
	ee engine.Engine,
	qm queue.Manager,
	sm state.Manager) (Worker, error) {
	switch workerType {
	case "submit":
		return &submitWorker{}, nil
	case "retry":
		return &retryWorker{
			sm:  sm,
			qm:  qm,
			log: log,
		}, nil
	case "status":
		return &statusWorker{}, nil
	case "reassign":
		return &reassignWorker{}, nil
	default:
		return nil, fmt.Errorf("No workerType %s exists", workerType)
	}
}