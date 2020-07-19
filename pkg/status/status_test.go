package status

import (
	"riser/pkg/util"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_GetRevisionStatus_PercentBasedRouting(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			Revisions: []model.DeploymentRevisionStatus{
				{
					Name: "rev0",
				},
				{
					Name: "rev1",
				},
				{
					Name: "rev2",
				},
			},
			Traffic: []model.DeploymentTrafficStatus{
				{
					RevisionName: "rev1",
					Percent:      util.PtrInt64(90),
				},
				{
					RevisionName: "rev2",
					Percent:      util.PtrInt64(10),
				},
			},
		},
	}

	result := GetRevisionStatus(deploymentStatus, true)

	assert.Len(t, result, 2)
	assert.Equal(t, "rev1", result[0].Name)
	assert.Equal(t, deploymentStatus.DeploymentStatusMutable.Traffic[0], result[0].Traffic)
	assert.Equal(t, "rev2", result[1].Name)
	assert.Equal(t, deploymentStatus.DeploymentStatusMutable.Traffic[1], result[1].Traffic)
}

// This may not be necessary as recent versions of KNative seems to report 100% in this case.
func Test_GetRevisionStatus_NoTraffic(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			Revisions: []model.DeploymentRevisionStatus{
				{
					Name: "rev0",
				},
				{
					Name: "rev1",
				},
			},
		},
	}

	result := GetRevisionStatus(deploymentStatus, true)

	assert.Len(t, result, 0)
}

func Test_GetRevisionStatus_LatestCreatedNoTraffic(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			LatestCreatedRevisionName: "rev2",
			Revisions: []model.DeploymentRevisionStatus{
				{
					Name:          "rev0",
					RiserRevision: 0,
				},
				{
					Name:          "rev1",
					RiserRevision: 1,
				},
				{
					Name:          "rev2",
					RiserRevision: 2,
				},
			},
			Traffic: []model.DeploymentTrafficStatus{
				{
					RevisionName: "rev0",
					Percent:      util.PtrInt64(100),
				},
			},
		},
	}

	result := GetRevisionStatus(deploymentStatus, true)

	assert.Len(t, result, 2)
	// Results should be sorted descending
	assert.Equal(t, "rev2", result[0].Name)
	assert.Equal(t, "rev0", result[1].Name)
	assert.Empty(t, result[0].Traffic)
}

func Test_GetRevisionStatus_AllRevisions(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			LatestCreatedRevisionName: "rev2",
			Revisions: []model.DeploymentRevisionStatus{
				{
					Name: "rev0",
				},
				{
					Name: "rev1",
				},
				{
					Name: "rev2",
				},
			},
			Traffic: []model.DeploymentTrafficStatus{
				{
					RevisionName: "rev0",
					Percent:      util.PtrInt64(100),
				},
			},
		},
	}

	result := GetRevisionStatus(deploymentStatus, false)

	assert.Len(t, result, 3)
}

func Test_GetRevisionStatus_SortsByRevision(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			Revisions: []model.DeploymentRevisionStatus{
				{
					RiserRevision: 0,
				},
				{
					RiserRevision: 2,
				},
				{
					RiserRevision: 1,
				},
			},
		},
	}

	result := GetRevisionStatus(deploymentStatus, false)

	assert.Len(t, result, 3)
	assert.EqualValues(t, 2, result[0].RiserRevision)
	assert.EqualValues(t, 1, result[1].RiserRevision)
	assert.EqualValues(t, 0, result[2].RiserRevision)
}

func Test_GetRevisionStatus_UnobservedStatus(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		RiserRevision: 1,
		DeploymentStatusMutable: model.DeploymentStatusMutable{
			ObservedRiserRevision:     0,
			LatestCreatedRevisionName: "rev2",
			Revisions: []model.DeploymentRevisionStatus{
				{
					Name:          "rev0",
					RiserRevision: 0,
				},
			},
			Traffic: []model.DeploymentTrafficStatus{
				{
					RevisionName: "rev0",
					Percent:      util.PtrInt64(100),
				},
			},
		},
	}

	result := GetRevisionStatus(deploymentStatus, true)

	assert.Len(t, result, 2)
	assert.Equal(t, deploymentStatus.RiserRevision, result[0].RiserRevision)
	assert.Equal(t, model.RevisionStatusWaiting, result[0].RevisionStatus)
	assert.Equal(t, "This revision has not yet been observed", result[0].RevisionStatusReason)
	assert.Equal(t, "rev0", result[1].Name)
	assert.Empty(t, result[0].Traffic)
}
