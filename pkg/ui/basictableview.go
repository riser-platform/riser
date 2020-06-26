package ui

import (
	"fmt"
	"io"
	"riser/pkg/ui/table"
)

// BasicTableView provides provides a View for basic table layouts.
// It assumes a simple 1:1 relationship between header and simple row values,
// allowing a simple JSON structure to be represented without any additional types
type BasicTableView struct {
	header []string
	rows   [][]interface{}
}

func NewBasicTableView() *BasicTableView {
	return &BasicTableView{}
}

func (view *BasicTableView) Header(header ...string) {
	view.header = header
}

// AddRow adds a row to the logical table. Each row must have the exact number of items as in the Header
func (view *BasicTableView) AddRow(values ...interface{}) {
	view.rows = append(view.rows, values)
}

func (view *BasicTableView) RenderHuman(writer io.Writer) error {
	table := table.Default().Header(view.header...)

	for _, row := range view.rows {
		table.AddRow(toString(row)...)
	}

	_, err := writer.Write([]byte(table.String() + "\n"))
	return err
}

func (view *BasicTableView) RenderJson(writer io.Writer) error {
	data := []map[string]interface{}{}
	for _, obj := range view.rows {
		jsonObj := map[string]interface{}{}
		for valIdx, val := range obj {
			jsonObj[view.header[valIdx]] = val
		}

		data = append(data, jsonObj)
	}

	return RenderJson(data, writer)
}

func toString(values []interface{}) []string {
	arr := make([]string, len(values))
	for i, val := range values {
		arr[i] = fmt.Sprintf("%v", val)
	}

	return arr
}
