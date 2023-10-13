package parsehtml

import (
	"errors"
	"strings"
)

type Table struct {
	ColumnName []string
	Row        [][]string
}

func (h *HTMLObjectWrapper) AllTableTag() []*Tag {
	return h.tables
}

func (t *Tag) ParseToTable() (*Table, error) {
	var result *Table = &Table{}

	// check if this tag is a table tag
	if t.Name != "table" {
		return nil, errors.New("this tag is not a table tag")
	}

	// get all column names
	var columnNames []string
	var theadTag *Tag = t.FirstTag("thead")
	if theadTag != nil {
		var trTag *Tag = theadTag.FirstTag("tr")
		if trTag != nil {
			for _, thTag := range trTag.AllTag("th") {
				columnNames = append(columnNames, strings.TrimSpace(thTag.Text[0]))
			}
		}
	}

	// get all rows
	var rows [][]string
	var tbodyTag *Tag = t.FirstTag("tbody")
	if tbodyTag != nil {
		for _, trTag := range tbodyTag.AllTag("tr") {
			var row []string
			for _, tdTag := range trTag.AllTag("td") {
				row = append(row, strings.TrimSpace(tdTag.Text[0]))
			}
			rows = append(rows, row)
		}
	}

	result.ColumnName = columnNames
	result.Row = rows
	return result, nil
}
