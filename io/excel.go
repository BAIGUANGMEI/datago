package io

import (
	"fmt"

	"github.com/BAIGUANGMEI/datago/dataframe"
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

// ExcelWriteOptions defines options for writing Excel files.
type ExcelWriteOptions struct {
	Sheet         string
	IncludeHeader *bool
	IncludeIndex  bool
	IndexName     string
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

// WriteExcel writes a DataFrame to an Excel file.
func WriteExcel(path string, df *dataframe.DataFrame, opts ExcelWriteOptions) error {
	if df == nil {
		return fmt.Errorf("dataframe is nil")
	}

	includeHeader := true
	if opts.IncludeHeader != nil {
		includeHeader = *opts.IncludeHeader
	}

	sheet := opts.Sheet
	if sheet == "" {
		sheet = "Sheet1"
	}

	f := excelize.NewFile()
	if sheet != "Sheet1" {
		if err := f.SetSheetName("Sheet1", sheet); err != nil {
			return err
		}
	}

	rows := df.Shape()[0]
	cols := df.Columns()

	rowOffset := 1
	colOffset := 1
	if includeHeader {
		if opts.IncludeIndex {
			indexName := opts.IndexName
			if indexName == "" {
				indexName = "index"
			}
			cell, _ := excelize.CoordinatesToCellName(colOffset, rowOffset)
			if err := f.SetCellValue(sheet, cell, indexName); err != nil {
				return err
			}
			colOffset++
		}

		for i, col := range cols {
			cell, _ := excelize.CoordinatesToCellName(colOffset+i, rowOffset)
			if err := f.SetCellValue(sheet, cell, col); err != nil {
				return err
			}
		}
		rowOffset++
	}

	for r := 0; r < rows; r++ {
		colStart := 1
		if opts.IncludeIndex {
			label, err := df.Index().Get(r)
			if err != nil {
				return err
			}
			cell, _ := excelize.CoordinatesToCellName(colStart, rowOffset+r)
			if err := f.SetCellValue(sheet, cell, label); err != nil {
				return err
			}
			colStart++
		}

		for c, col := range cols {
			series, ok := df.GetSeries(col)
			if !ok {
				return fmt.Errorf("column '%s' not found", col)
			}
			value, err := series.Get(r)
			if err != nil {
				return err
			}
			cell, _ := excelize.CoordinatesToCellName(colStart+c, rowOffset+r)
			if value == nil {
				value = ""
			}
			if err := f.SetCellValue(sheet, cell, value); err != nil {
				return err
			}
		}
	}

	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}

// WriteSeriesExcel writes a Series to an Excel file.
func WriteSeriesExcel(path string, s *dataframe.Series, opts ExcelWriteOptions) error {
	if s == nil {
		return fmt.Errorf("series is nil")
	}

	includeHeader := true
	if opts.IncludeHeader != nil {
		includeHeader = *opts.IncludeHeader
	}

	sheet := opts.Sheet
	if sheet == "" {
		sheet = "Sheet1"
	}

	f := excelize.NewFile()
	if sheet != "Sheet1" {
		if err := f.SetSheetName("Sheet1", sheet); err != nil {
			return err
		}
	}

	rowOffset := 1
	colOffset := 1
	if includeHeader {
		if opts.IncludeIndex {
			indexName := opts.IndexName
			if indexName == "" {
				indexName = "index"
			}
			cell, _ := excelize.CoordinatesToCellName(colOffset, rowOffset)
			if err := f.SetCellValue(sheet, cell, indexName); err != nil {
				return err
			}
			colOffset++
		}

		seriesName := s.Name()
		if seriesName == "" {
			seriesName = "value"
		}
		cell, _ := excelize.CoordinatesToCellName(colOffset, rowOffset)
		if err := f.SetCellValue(sheet, cell, seriesName); err != nil {
			return err
		}
		rowOffset++
	}

	for i := 0; i < s.Len(); i++ {
		colStart := 1
		if opts.IncludeIndex {
			label, err := s.Index().Get(i)
			if err != nil {
				return err
			}
			cell, _ := excelize.CoordinatesToCellName(colStart, rowOffset+i)
			if err := f.SetCellValue(sheet, cell, label); err != nil {
				return err
			}
			colStart++
		}

		value, err := s.Get(i)
		if err != nil {
			return err
		}
		cell, _ := excelize.CoordinatesToCellName(colStart, rowOffset+i)
		if value == nil {
			value = ""
		}
		if err := f.SetCellValue(sheet, cell, value); err != nil {
			return err
		}
	}

	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}
