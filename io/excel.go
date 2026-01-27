package io

import (
	"fmt"

	"github.com/datago/dataframe"
	"github.com/xuri/excelize/v2"
)

// ExcelOptions defines options for reading Excel files.
type ExcelOptions struct {
	Sheet     string
	HasHeader bool
	SkipRows  int
	UseCols   []string
	DTypes    map[string]dataframe.DType
}

// ReadExcel reads an Excel file and returns a DataFrame.
func ReadExcel(path string, opts ExcelOptions) (*dataframe.DataFrame, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	sheet := opts.Sheet
	if sheet == "" {
		sheet = f.GetSheetName(0)
		if sheet == "" {
			return nil, fmt.Errorf("no sheet found in excel file")
		}
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return dataframe.New(map[string][]interface{}{})
	}

	startRow := opts.SkipRows
	if startRow >= len(rows) {
		return dataframe.New(map[string][]interface{}{})
	}

	var columns []string
	dataStart := startRow
	if opts.HasHeader && startRow < len(rows) {
		columns = make([]string, len(rows[startRow]))
		for i, col := range rows[startRow] {
			if col == "" {
				columns[i] = fmt.Sprintf("col_%d", i)
			} else {
				columns[i] = col
			}
		}
		dataStart = startRow + 1
	} else {
		if len(rows[startRow]) == 0 {
			return dataframe.New(map[string][]interface{}{})
		}
		columns = make([]string, len(rows[startRow]))
		for i := range columns {
			columns[i] = fmt.Sprintf("col_%d", i)
		}
	}

	// Filter columns if UseCols is provided
	useCols := make(map[string]bool)
	if len(opts.UseCols) > 0 {
		for _, c := range opts.UseCols {
			useCols[c] = true
		}
	}

	colData := make(map[string][]interface{})
	colIndex := make([]int, 0, len(columns))
	selectedCols := make([]string, 0, len(columns))
	for i, col := range columns {
		if len(useCols) == 0 || useCols[col] {
			colData[col] = []interface{}{}
			colIndex = append(colIndex, i)
			selectedCols = append(selectedCols, col)
		}
	}

	for i := dataStart; i < len(rows); i++ {
		row := rows[i]
		for j, colIdx := range colIndex {
			col := selectedCols[j]
			if colIdx < len(row) {
				colData[col] = append(colData[col], row[colIdx])
			} else {
				colData[col] = append(colData[col], nil)
			}
		}
	}

	df, err := dataframe.New(colData)
	if err != nil {
		return nil, err
	}

	// Apply dtypes if provided
	for col, dtype := range opts.DTypes {
		if s, ok := df.GetSeries(col); ok {
			converted, err := s.AsType(dtype)
			if err == nil {
				_ = df.SetColumn(col, converted)
			}
		}
	}

	return df, nil
}
