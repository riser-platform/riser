package ui

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RenderViewWriter(t *testing.T) {
	tests := []struct {
		outputFormat string
		renderHuman  bool
		renderJson   bool
	}{
		{
			outputFormat: OutputFormatHuman,
			renderHuman:  true,
		},
		{
			outputFormat: OutputFormatJson,
			renderJson:   true,
		},
	}

	for _, tt := range tests {
		view := &FakeView{}
		SetOutputFormat(tt.outputFormat)
		err := RenderViewWriter(view, ioutil.Discard)

		assert.NoError(t, err)
		assert.Equal(t, tt.renderHuman, view.RenderHumanCalled)
		assert.Equal(t, tt.renderJson, view.RenderJsonCalled)
	}
}
