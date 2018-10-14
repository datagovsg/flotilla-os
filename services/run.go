package services

import (
	"github.com/stitchfix/flotilla-os/config"
	"github.com/stitchfix/flotilla-os/state"
)

type RunService interface {
	UpdateTags(runID string, userTags state.UserTagMap, taskID string) error
}

type runService struct {
	sm state.Manager
}

func NewRunService(conf config.Config, sm state.Manager) (RunService, error) {
	return &runService{sm}, nil
}

func (rs *runService) UpdateTags(runID string, userTags state.UserTagMap, taskID string) error {
	return rs.sm.UpdateRunTags(runID, userTags, taskID)
}
