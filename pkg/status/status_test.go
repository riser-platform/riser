package status

import (
	"riser/pkg/util"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_GetActiveRevisions_PercentBasedRouting(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			Revisions: []model.DeploymentRevisionStatus{
				model.DeploymentRevisionStatus{
					Name: "rev0",
				},
				model.DeploymentRevisionStatus{
					Name: "rev1",
				},
				model.DeploymentRevisionStatus{
					Name: "rev2",
				},
			},
			Traffic: []model.DeploymentTrafficStatus{
				model.DeploymentTrafficStatus{
					RevisionName: "rev1",
					Percent:      util.PtrInt64(90),
				},
				model.DeploymentTrafficStatus{
					RevisionName: "rev2",
					Percent:      util.PtrInt64(10),
				},
			},
		},
	}

	result := GetActiveRevisions(deploymentStatus)

	assert.Len(t, result, 2)
	assert.Equal(t, "rev1", result[0].Name)
	assert.Equal(t, deploymentStatus.DeploymentStatusMutable.Traffic[0], result[0].Traffic)
	assert.Equal(t, "rev2", result[1].Name)
	assert.Equal(t, deploymentStatus.DeploymentStatusMutable.Traffic[1], result[1].Traffic)
}

// This may not be necessary as recent versions of KNative seems to report 100% in this case.
func Test_GetActiveRevisions_NoTraffic(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			Revisions: []model.DeploymentRevisionStatus{
				model.DeploymentRevisionStatus{
					Name: "rev0",
				},
				model.DeploymentRevisionStatus{
					Name: "rev1",
				},
			},
		},
	}

	result := GetActiveRevisions(deploymentStatus)

	assert.Len(t, result, 0)
}

func Test_GetActiveRevisions_LatestCreatedNoTraffic(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			LatestCreatedRevisionName: "rev2",
			Revisions: []model.DeploymentRevisionStatus{
				model.DeploymentRevisionStatus{
					Name: "rev0",
				},
				model.DeploymentRevisionStatus{
					Name: "rev1",
				},
				model.DeploymentRevisionStatus{
					Name: "rev2",
				},
			},
			Traffic: []model.DeploymentTrafficStatus{
				model.DeploymentTrafficStatus{
					RevisionName: "rev0",
					Percent:      util.PtrInt64(100),
				},
			},
		},
	}

	result := GetActiveRevisions(deploymentStatus)

	assert.Len(t, result, 2)
	assert.Equal(t, "rev0", result[0].Name)
	assert.Equal(t, "rev2", result[1].Name)
	assert.Empty(t, result[1].Traffic)
}
