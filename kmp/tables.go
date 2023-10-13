package kmp

import (
	"fmt"
	"strings"
)

type Table struct {
	Index       int
	Start       int
	End         int
	Atrributes  map[string]string
	ColumnsName []string
	RowData     [][]string
}

type TableTagError struct {
	errorMessage string
}

func (t *TableTagError) Error() string {
	return fmt.Sprintf("TableTagError: %s", t.errorMessage)
}

// search for all table tag in the html
// return a list of indexes
func SearchAllTableTag(html string) []int {
	return SearchAllOccurrences("<table", html)
}

func SearchAllTableTagEnd(html string) []int {
	return SearchAllOccurrences("</table", html)
}

func ProcessTableTag(html string) ([]Table, error) {
	var result []Table
	var tableTag []int = SearchAllTableTag(html)
	var tableTagEnd []int = SearchAllTableTagEnd(html)

	if len(tableTag) == len(tableTagEnd) && len(tableTag) == 1 {
		return append(result, Table{
			Index: 0,
			Start: tableTag[0],
			End:   tableTagEnd[0],
		}), nil
	}

	if len(tableTag) != len(tableTagEnd) {
		return nil, &TableTagError{errorMessage: "Asymmetric table tag - <table> and </table> are not equal"}
	}

	if len(tableTag) == 0 {
		return nil, &TableTagError{errorMessage: "No table tag found"}
	}

	for i := range tableTag {
		newTable := Table{
			Index: i,
			Start: tableTag[i],
			End:   tableTagEnd[i],
		}
		result = append(result, newTable)
	}

	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].End > result[j].Start {
				temp := result[i].End
				result[i].End = result[j].End
				result[j].End = temp
				continue
			}
			break
		}
	}

	for i := range result {
		ParseTable(html, &result[i])
	}
	return result, nil
}

func ParseTable(html string, tb *Table) {
	var tableContent string = html[tb.Start:tb.End]
	ExtractTableAttributes(tableContent, tb)

	endOfHeader := strings.Index(tableContent, "</tr>")
	if endOfHeader == -1 {
		return
	}
	var tableHeader string = tableContent[:endOfHeader]
	var tableBody string = tableContent[endOfHeader:]

	startOfHeader := strings.Index(tableHeader, "<tr")
	endOfHeader = strings.Index(tableHeader, "</tr>")
	startOfBody := strings.Index(tableBody, "<tr")
	endOfBody := strings.Index(tableBody, "</tr>")
	if startOfHeader == -1 || endOfHeader == -1 || startOfBody == -1 || endOfBody == -1 {
		return
	}
	var tableHeaderContent string = tableHeader[startOfHeader:endOfHeader]
	var tableBodyContent string = tableBody[startOfBody:endOfBody]

	var headerRow []string = ParseRow(tableHeaderContent)
	var bodyRow [][]string = ParseRows(tableBodyContent)

	tb.ColumnsName = headerRow
	tb.RowData = bodyRow
}

func ParseRow(row string) []string {
	var result []string
	startOfRowContent := strings.Index(row, "<tr")
	endOfRowContent := strings.Index(row, "</tr>")
	var rowContent string = row[startOfRowContent:endOfRowContent]
	var rowContentTag []int = SearchAllOccurrences("<td", rowContent)
	var rowContentTagEnd []int = SearchAllOccurrences("</td>", rowContent)

	if len(rowContentTag) != len(rowContentTagEnd) {
		return nil
	}

	for i := range rowContentTag {
		result = append(result, rowContent[rowContentTag[i]:rowContentTagEnd[i]])
	}
	return result
}

func ParseRows(rows string) [][]string {
	var result [][]string
	var rowsContent []string = strings.Split(rows, "</tr>")
	for i := range rowsContent {
		result = append(result, ParseRow(rowsContent[i]))
	}
	return result
}

func ExtractTableAttributes(tableContent string, tb *Table) {
	var tableAttributes string = tableContent[:strings.Index(tableContent, ">")]
	var tableAttributesList []string = strings.Split(tableAttributes, " ")
	tb.Atrributes = make(map[string]string)
	if len(tableAttributesList) == 1 {
		return
	}
	for i := range tableAttributesList {
		var attr string = tableAttributesList[i]
		var equalSignIndex int = strings.Index(attr, "=")
		if equalSignIndex != -1 {
			var attrName string = attr[1:equalSignIndex]
			var attrValue string = attr[equalSignIndex+1 : len(attr)-1]
			tb.Atrributes[attrName] = attrValue
		}
	}
}
