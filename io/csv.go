package io

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/datago/dataframe"
)

// CSVOptions defines options for reading CSV files.
type CSVOptions struct {
	Separator rune
	HasHeader bool
	SkipRows  int
	UseCols   []string
	DTypes    map[string]dataframe.DType
}

// CSVWriteOptions defines options for writing CSV files.
type CSVWriteOptions struct {
	Separator     rune
	IncludeHeader *bool
	IncludeIndex  bool
	IndexName     string
}

// ReadCSV reads a CSV file and returns a DataFrame.
func ReadCSV(path string, opts CSVOptions) (*dataframe.DataFrame, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	if opts.Separator != 0 {
		reader.Comma = opts.Separator
	}
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
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

// WriteCSV writes a DataFrame to a CSV file.
func WriteCSV(path string, df *dataframe.DataFrame, opts CSVWriteOptions) error {
	if df == nil {
		return fmt.Errorf("dataframe is nil")
	}

	includeHeader := true
	if opts.IncludeHeader != nil {
		includeHeader = *opts.IncludeHeader
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	if opts.Separator != 0 {
		writer.Comma = opts.Separator
	}
	defer writer.Flush()

	cols := df.Columns()
	rows := df.Shape()[0]

	if includeHeader {
		header := make([]string, 0, len(cols)+1)
		if opts.IncludeIndex {
			indexName := opts.IndexName
			if indexName == "" {
				indexName = "index"
			}
			header = append(header, indexName)
		}
		header = append(header, cols...)
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	for r := 0; r < rows; r++ {
		record := make([]string, 0, len(cols)+1)
		if opts.IncludeIndex {
			label, err := df.Index().Get(r)
			if err != nil {
				return err
			}
			record = append(record, fmt.Sprintf("%v", label))
		}
		for _, col := range cols {
			series, ok := df.GetSeries(col)
			if !ok {
				return fmt.Errorf("column '%s' not found", col)
			}
			value, err := series.Get(r)
			if err != nil {
				return err
			}
			if value == nil {
				record = append(record, "")
			} else {
				record = append(record, fmt.Sprintf("%v", value))
			}
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return writer.Error()
}

// WriteSeriesCSV writes a Series to a CSV file.
func WriteSeriesCSV(path string, s *dataframe.Series, opts CSVWriteOptions) error {
	if s == nil {
		return fmt.Errorf("series is nil")
	}

	includeHeader := true
	if opts.IncludeHeader != nil {
		includeHeader = *opts.IncludeHeader
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	if opts.Separator != 0 {
		writer.Comma = opts.Separator
	}
	defer writer.Flush()

	if includeHeader {
		header := make([]string, 0, 2)
		if opts.IncludeIndex {
			indexName := opts.IndexName
			if indexName == "" {
				indexName = "index"
			}
			header = append(header, indexName)
		}
		seriesName := s.Name()
		if seriesName == "" {
			seriesName = "value"
		}
		header = append(header, seriesName)
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	for i := 0; i < s.Len(); i++ {
		record := make([]string, 0, 2)
		if opts.IncludeIndex {
			label, err := s.Index().Get(i)
			if err != nil {
				return err
			}
			record = append(record, fmt.Sprintf("%v", label))
		}
		value, err := s.Get(i)
		if err != nil {
			return err
		}
		if value == nil {
			record = append(record, "")
		} else {
			record = append(record, fmt.Sprintf("%v", value))
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return writer.Error()
}
