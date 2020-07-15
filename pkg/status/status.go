package status

import (
	"sort"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type RevisionStatusWithTraffic struct {
	model.DeploymentRevisionStatus
	Traffic model.DeploymentTrafficStatus
}

/*
GetRevisionStatus returns all revisions with traffic status. If activeOnly is true, it will only return active revisions.
An active revision is:
- A revision that has some % of traffic going to it
- The latest revision which may or may not yet be receiving traffic
*/
func GetRevisionStatus(deploymentStatus *model.DeploymentStatus, activeOnly bool) []RevisionStatusWithTraffic {
	activeStatuses := []RevisionStatusWithTraffic{}

	for _, revision := range deploymentStatus.Revisions {
		hasTrafficStatus := false
		for _, traffic := range deploymentStatus.Traffic {
			if traffic.RevisionName == revision.Name && hasTraffic(&traffic) {
				hasTrafficStatus = true
				activeStatuses = append(activeStatuses, RevisionStatusWithTraffic{revision, traffic})
			}
		}

		if !hasTrafficStatus && (!activeOnly || revision.Name == deploymentStatus.LatestCreatedRevisionName) {
			activeStatuses = append(activeStatuses, RevisionStatusWithTraffic{DeploymentRevisionStatus: revision})
		}
	}

	sort.Slice(activeStatuses, func(i, j int) bool {
		return activeStatuses[i].RiserRevision > activeStatuses[j].RiserRevision
	})

	return activeStatuses
}

func hasTraffic(traffic *model.DeploymentTrafficStatus) bool {
	return traffic.Percent != nil && *traffic.Percent > 0
}
