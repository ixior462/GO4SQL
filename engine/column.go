package engine

import (
	"github.com/LissaGreense/GO4SQL/token"
)

// Column - part of the Table containing name of Column and values in it
type Column struct {
	Name   string
	Type   token.Token
	Values []ValueInterface
}

func extractColumnContent(columns []*Column, wantedColumnNames *[]string, tableName string) (*Table, error) {
	selectedTable := &Table{Columns: make([]*Column, 0)}
	mappedIndexes := make([]int, 0)
	for wantedColumnIndex := range *wantedColumnNames {
		for columnNameIndex := range columns {
			if columns[columnNameIndex].Name == (*wantedColumnNames)[wantedColumnIndex] {
				mappedIndexes = append(mappedIndexes, columnNameIndex)
				break
			}
			if columnNameIndex == len(columns)-1 {
				return nil, &ColumnDoesNotExistError{columnName: (*wantedColumnNames)[wantedColumnIndex], tableName: tableName}
			}
		}
	}

	for i := range mappedIndexes {
		selectedTable.Columns = append(selectedTable.Columns, &Column{
			Name:   columns[mappedIndexes[i]].Name,
			Type:   columns[mappedIndexes[i]].Type,
			Values: make([]ValueInterface, 0),
		})
	}
	if len(columns) == 0 {
		return selectedTable, nil
	}

	rowsCount := len(columns[0].Values)

	for iRow := 0; iRow < rowsCount; iRow++ {
		for iColumn := 0; iColumn < len(mappedIndexes); iColumn++ {
			selectedTable.Columns[iColumn].Values =
				append(selectedTable.Columns[iColumn].Values, columns[mappedIndexes[iColumn]].Values[iRow])
		}
	}
	return selectedTable, nil
}
