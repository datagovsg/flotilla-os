package cluster

import (
	"github.com/datagovsg/flotilla-os/config"
	"github.com/datagovsg/flotilla-os/state"
)

// fill in with do-nothing functions until we think about what to do
type NomadClusterClient struct {
}

//
// Name is the name of the client
//
func (ncc *NomadClusterClient) Name() string {
	return "nomad"
}

func (ncc *NomadClusterClient) Initialize(conf config.Config) error {
	return nil
}

//
// CanBeRun determines whether a task formed from the specified definition
// can be run on clusterName
//
func (ncc *NomadClusterClient) CanBeRun(clusterName string, definition state.Definition) (bool, error) {
	return true, nil
}

//
// ListClusters gets a list of cluster names
//
func (ncc *NomadClusterClient) ListClusters() ([]string, error) {
	return []string{}, nil
}
