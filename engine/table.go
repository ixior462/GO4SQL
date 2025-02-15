package engine

import (
	"github.com/LissaGreense/GO4SQL/token"
	"hash/adler32"
)

// Table - Contain Columns that store values in engine
type Table struct {
	Columns []*Column
}

func (table *Table) isEqual(secondTable *Table) bool {
	if len(table.Columns) != len(secondTable.Columns) {
		return false
	}

	for i := range table.Columns {
		if table.Columns[i].Name != secondTable.Columns[i].Name {
			return false
		}
		if table.Columns[i].Type.Literal != secondTable.Columns[i].Type.Literal {
			return false
		}
		if table.Columns[i].Type.Type != secondTable.Columns[i].Type.Type {
			return false
		}
		if len(table.Columns[i].Values) != len(secondTable.Columns[i].Values) {
			return false
		}
		for j := range table.Columns[i].Values {
			if table.Columns[i].Values[j].ToString() != secondTable.Columns[i].Values[j].ToString() {
				return false
			}
		}
	}

	return true
}

// getDistinctTable - Takes input table, and returns new one without any duplicates
func (table *Table) getDistinctTable() *Table {
	distinctTable := getCopyOfTableWithoutRows(table)

	rowsCount := len(table.Columns[0].Values)

	checksumSet := map[uint32]struct{}{}

	for iRow := 0; iRow < rowsCount; iRow++ {

		mergedColumnValues := ""
		for iColumn := range table.Columns {
			fieldValue := table.Columns[iColumn].Values[iRow].ToString()
			if table.Columns[iColumn].Type.Literal == token.TEXT {
				fieldValue = "'" + fieldValue + "'"
			}
			mergedColumnValues += fieldValue
		}
		checksum := adler32.Checksum([]byte(mergedColumnValues))

		_, exist := checksumSet[checksum]
		if !exist {
			checksumSet[checksum] = struct{}{}
			for i, column := range distinctTable.Columns {
				column.Values = append(column.Values, table.Columns[i].Values[iRow])
			}
		}
	}

	return distinctTable
}

// ToString - Return string contain all values and Column names in Table
func (table *Table) ToString() string {
	columWidths := getColumWidths(table.Columns)
	bar := getBar(columWidths)
	result := bar + "\n"

	result += "|"
	for i := range table.Columns {
		result += " "
		for j := 0; j < columWidths[i]-len(table.Columns[i].Name); j++ {
			result += " "
		}
		result += table.Columns[i].Name
		result += " |"
	}
	result += "\n" + bar + "\n"

	if len(table.Columns) == 0 {
		return result
	}

	rowsCount := len(table.Columns[0].Values)

	for iRow := 0; iRow < rowsCount; iRow++ {
		result += "|"

		for iColumn := range table.Columns {
			result += " "

			printedValue := table.Columns[iColumn].Values[iRow].ToString()
			if table.Columns[iColumn].Type.Literal == token.TEXT &&
				table.Columns[iColumn].Values[iRow].GetType() != NullType {
				printedValue = "'" + printedValue + "'"
			}
			for i := 0; i < columWidths[iColumn]-len(printedValue); i++ {
				result += " "
			}

			result += printedValue + " |"
		}

		result += "\n"
	}

	return result + bar
}

func (table *Table) getTableCopyWithAddedPrefixToColumnNames(columnNamePrefix string) *Table {
	newTable := &Table{Columns: []*Column{}}

	for _, column := range table.Columns {
		newTable.Columns = append(newTable.Columns,
			&Column{
				Type:   column.Type,
				Values: column.Values,
				Name:   columnNamePrefix + column.Name,
			})
	}

	return newTable
}

func getBar(columWidths []int) string {
	bar := "+"

	for i := 0; i < len(columWidths); i++ {
		bar += "-"
		for j := 0; j < columWidths[i]; j++ {
			bar += "-"
		}
		bar += "-+"
	}

	return bar
}

func getColumWidths(columns []*Column) []int {
	widths := make([]int, 0)

	for iColumn := range columns {
		maxLength := len(columns[iColumn].Name)
		for iRow := range columns[iColumn].Values {
			valueLength := len(columns[iColumn].Values[iRow].ToString())
			if columns[iColumn].Type.Literal == token.TEXT {
				valueLength += 2 // double '
			}
			if valueLength > maxLength {
				maxLength = valueLength
			}
		}
		widths = append(widths, maxLength)
	}

	return widths
}
