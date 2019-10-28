package table

import (
	"github.com/jedib0t/go-pretty/table"
)

const defaultMaxColumnLength = 50

type Table struct {
	internal table.Writer
}

// Default creates a default table
func Default() *Table {
	internal := table.NewWriter()
	internal.SetStyle(table.StyleLight)
	internal.Style().Options.DrawBorder = false
	internal.Style().Options.SeparateColumns = false
	return &Table{internal}
}

func (t *Table) Header(values ...string) *Table {
	t.internal.AppendHeader(t.createRow(values))
	columnLengths := []int{}
	for range values {
		columnLengths = append(columnLengths, defaultMaxColumnLength)
	}
	t.internal.SetAllowedColumnLengths(columnLengths)
	return t
}

func (t *Table) AddRow(values ...string) *Table {
	t.internal.AppendRow(t.createRow(values))
	return t
}

func (t *Table) String() string {
	return t.internal.Render()
}

func (t *Table) createRow(values []string) table.Row {
	row := table.Row{}
	for _, v := range values {
		row = append(row, v)
	}
	return row
}
