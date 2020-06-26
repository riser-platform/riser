package cmd

import (
	"bytes"
	"fmt"
	"riser/pkg/ui/style"
	"riser/pkg/util"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_RenderHuman_NoDeployments(t *testing.T) {
	view := &statusView{
		appName: "myapp",
		status:  &model.AppStatus{},
	}

	var b bytes.Buffer

	err := view.RenderHuman(&b)

	assert.NoError(t, err)
	assert.Equal(t, "There are no deployments for the app \"myapp\". Use \"riser deploy\" to make your first deployment.\n", b.String())
}

func Test_deploymentObserved(t *testing.T) {
	newStatus := func(rev int64, observed int64) model.DeploymentStatus {
		return model.DeploymentStatus{
			RiserRevision: rev,
			DeploymentStatusMutable: model.DeploymentStatusMutable{
				ObservedRiserRevision: observed,
			},
			// For test readability
			DeploymentName: fmt.Sprintf("rev: %d, observed: %d", rev, observed),
		}
	}

	tests := []struct {
		status           model.DeploymentStatus
		expectedObserved bool
	}{
		{
			newStatus(1, 1),
			true,
		},
		{
			newStatus(2, 1),
			false,
		},
		{
			newStatus(1, 2),
			true,
		},
	}

	view := statusView{}

	for _, tt := range tests {
		result := view.deploymentObserved(tt.status)
		assert.Equal(t, tt.expectedObserved, result, tt.status.DeploymentName)
	}
}

func Test_formatDockerTag(t *testing.T) {
	tests := []struct {
		dockerImage string
		expected    string
	}{
		{"foo:v1", "v1"},
		{"foo", style.Warn("Unknown")},
	}

	view := statusView{}

	for _, tt := range tests {
		result := view.formatDockerTag(tt.dockerImage)
		assert.Equal(t, tt.expected, result)
	}
}

func Test_formatDeploymentName(t *testing.T) {
	newStatus := func(rev int64, observed int64) model.DeploymentStatus {
		return model.DeploymentStatus{
			RiserRevision: rev,
			DeploymentStatusMutable: model.DeploymentStatusMutable{
				ObservedRiserRevision: observed,
			},
			DeploymentName: "mydeployment",
		}
	}

	tests := []struct {
		status   model.DeploymentStatus
		expected string
	}{
		{
			newStatus(1, 1),
			"mydeployment",
		},
		{
			newStatus(2, 1),
			style.Emphasis("*") + "mydeployment",
		},
	}

	view := statusView{}

	for _, tt := range tests {
		result := view.formatDeploymentName(tt.status)
		assert.Equal(t, tt.expected, result)
	}
}

func Test_formatTraffic(t *testing.T) {
	tests := []struct {
		percent  *int64
		expected string
	}{
		{
			util.PtrInt64(100),
			"100%",
		},
		{
			nil,
			"0%",
		},
	}

	view := statusView{}

	for _, tt := range tests {
		result := view.formatTraffic(&model.DeploymentTrafficStatus{Percent: tt.percent})
		assert.Equal(t, tt.expected, result)
	}
}
