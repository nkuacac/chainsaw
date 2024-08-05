package table

import (
	"github.com/jedib0t/go-pretty/table"
)

func RenderWithName(name string, header []string, rows [][]interface{}) string {
	t := table.NewWriter()
	t.Style().Options = table.OptionsNoBordersAndSeparators
	if len(header) > 0 {
		new_header := convertToInterface(header...)
		new_header[0] = name
		t.AppendHeader(new_header)
	}
	for i, row := range rows {
		new_row := []interface{}{i}
		new_row = append(new_row, row...)
		t.AppendRow(new_row)
	}
	return t.Render()
}

func RenderCSV(rows [][]interface{}) string {
	t := table.NewWriter()
	for _, row := range rows {
		t.AppendRow(row)
	}
	return t.RenderCSV()
}

func convertToInterface(t ...string) []interface{} {
	s := make([]interface{}, len(t)+1)
	for i, v := range t {
		s[i+1] = v
	}
	return s
}
