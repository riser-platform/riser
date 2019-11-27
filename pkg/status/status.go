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
		for _, traffic := range deploymentStatus.Traffic {
			if traffic.RevisionName == revision.Name && hasActiveRoute(&traffic) {
				activeStatuses = append(activeStatuses, RevisionStatusWithTraffic{revision, traffic})
			} else if revision.Name == deploymentStatus.LatestCreatedRevisionName {
				activeStatuses = append(activeStatuses, RevisionStatusWithTraffic{DeploymentRevisionStatus: revision})
			}
		}
	}

	return activeStatuses
}

func hasActiveRoute(traffic *model.DeploymentTrafficStatus) bool {
	// TODO: This is not strictly correct, although it may be fine in practice. Need to test once we implement rollout. In KNative,
	// it's possible to not define percent and just say "latest=true". However  it appears that while percent is optional it always returns some value
	// in the status, so we may be able to just remove the latest check alltogether.
	return traffic.Percent != nil || (traffic.Latest != nil && *traffic.Latest)
}
