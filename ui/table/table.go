package table

import (
	"github.com/alexeyco/simpletable"
)

type Table struct {
	internal *simpletable.Table
}

// Default creates a default table
func Default() *Table {
	internal := simpletable.New()
	internal.SetStyle(simpletable.StyleCompactLite)
	internal.Header = &simpletable.Header{Cells: []*simpletable.Cell{}}
	return &Table{internal}
}

func (t *Table) Header(values ...string) *Table {
	for _, v := range values {
		t.internal.Header.Cells = append(t.internal.Header.Cells, defaultCell(v))
	}
	return t
}

func (t *Table) AddRow(values ...string) *Table {
	cells := []*simpletable.Cell{}
	for _, v := range values {
		cells = append(cells, defaultCell(v))
	}
	t.internal.Body.Cells = append(t.internal.Body.Cells, cells)
	return t
}

func (t *Table) String() string {
	return t.internal.String()
}

func defaultCell(text string) *simpletable.Cell {
	return &simpletable.Cell{Align: simpletable.AlignLeft, Text: text}
}
