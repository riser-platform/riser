package status

import (
	"github.com/riser-platform/riser-server/api/v1/model"
)

func GetLatestReadyRevision(deploymentStatus *model.DeploymentStatus) *model.DeploymentRevisionStatus {
	for _, revision := range deploymentStatus.Revisions {
		if revision.Name == deploymentStatus.LatestReadyRevisionName {
			return &revision
		}
	}

	return nil
}
