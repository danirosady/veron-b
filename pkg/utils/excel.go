package utils

import (
	"github.com/xuri/excelize/v2"
)

func ExportToExcel(filename string, sheetName string, headers []string, data [][]interface{}) error {
	f := excelize.NewFile()
	defer f.Close()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	f.SetActiveSheet(index)

	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#CCCCCC"},
			Pattern: 1,
		},
	})
	if err != nil {
		return err
	}

	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, styleHeader)
	}

	for rowIdx, row := range data {
		for colIdx, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}
