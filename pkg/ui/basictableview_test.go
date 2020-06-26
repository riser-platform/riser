package ui

import (
	"bytes"
	"riser/pkg/ui/table"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RenderHuman(t *testing.T) {
	view := NewBasicTableView()
	view.Header("H1", "H2")
	view.AddRow("h1v1", "h1v2")
	view.AddRow("h2v1", "h2v2")

	table := table.Default().Header("H1", "H2")
	table.AddRow("h1v1", "h1v2")
	table.AddRow("h2v1", "h2v2")

	var b bytes.Buffer

	err := view.RenderHuman(&b)

	assert.NoError(t, err)

	assert.Regexp(t, table.String(), b.String())
}

func Test_RenderJson(t *testing.T) {
	view := NewBasicTableView()
	view.Header("H1", "H2")
	view.AddRow("h1v1", "h1v2")
	view.AddRow("h2v1", "h2v2")

	var b bytes.Buffer

	err := view.RenderJson(&b)

	assert.NoError(t, err)

	assert.Regexp(t, `[{"H1":"h1v1","H2":"h1v2"},{"H1":"h2v1","H2":"h2v2"}]`, b.String())
}
