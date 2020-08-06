package cmd

import (
	"bytes"
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

func Test_formatDockerTag(t *testing.T) {
	tests := []struct {
		dockerImage string
		expected    string
	}{
		{"foo:v1", "v1"},
		{"foo", style.Warn("Unknown")},
	}

	for _, tt := range tests {
		result := formatDockerTag(tt.dockerImage)
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

	for _, tt := range tests {
		result := formatTraffic(&model.DeploymentTrafficStatus{Percent: tt.percent})
		assert.Equal(t, tt.expected, result)
	}
}
