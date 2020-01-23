package status

import (
	"github.com/riser-platform/riser-server/api/v1/model"
)

type RevisionStatusWithTraffic struct {
	model.DeploymentRevisionStatus
	Traffic model.DeploymentTrafficStatus
}

/*
GetActiveRevisions returns all active revisions. An active revision is:
- A revision that has some % of traffic going to it
- The latest revision which may or may not yet be receiving traffic
*/
func GetActiveRevisions(deploymentStatus *model.DeploymentStatus) []RevisionStatusWithTraffic {
	activeStatuses := []RevisionStatusWithTraffic{}
	for _, revision := range deploymentStatus.Revisions {
		hasTrafficStatus := false
		for _, traffic := range deploymentStatus.Traffic {
			if traffic.RevisionName == revision.Name && isActive(&traffic) {
				hasTrafficStatus = true
				activeStatuses = append(activeStatuses, RevisionStatusWithTraffic{revision, traffic})
			}
		}

		if !hasTrafficStatus && revision.Name == deploymentStatus.LatestCreatedRevisionName {
			activeStatuses = append(activeStatuses, RevisionStatusWithTraffic{DeploymentRevisionStatus: revision})
		}
	}

	return activeStatuses
}

func isActive(traffic *model.DeploymentTrafficStatus) bool {
	return traffic.Percent != nil && *traffic.Percent > 0
}
