package dataframe

import (
	"fmt"
	"strings"
)

// DataFrame represents a 2-dimensional labeled data structure.
type DataFrame struct {
	columns []string
	data    map[string]*Series
	index   *Index
	shape   [2]int // [rows, cols]
}

// Row represents a single row of a DataFrame.
type Row struct {
	data map[string]interface{}
}

// Get returns the value for the given column name.
func (r Row) Get(column string) interface{} {
	return r.data[column]
}

// New creates a DataFrame from a map of column name to values.
func New(data map[string][]interface{}) (*DataFrame, error) {
	if len(data) == 0 {
		return &DataFrame{columns: []string{}, data: map[string]*Series{}, index: NewRangeIndex(0), shape: [2]int{0, 0}}, nil
	}

	columns := make([]string, 0, len(data))
	var rowCount int
	for col, values := range data {
		columns = append(columns, col)
		if rowCount == 0 {
			rowCount = len(values)
		} else if len(values) != rowCount {
			return nil, fmt.Errorf("column '%s' length %d does not match %d", col, len(values), rowCount)
		}
	}

	seriesMap := make(map[string]*Series)
	for _, col := range columns {
		seriesMap[col] = NewSeries(data[col], col)
	}

	return &DataFrame{
		columns: columns,
		data:    seriesMap,
		index:   NewRangeIndex(rowCount),
		shape:   [2]int{rowCount, len(columns)},
	}, nil
}

// FromRecords creates a DataFrame from records and columns.
func FromRecords(records [][]interface{}, columns []string) (*DataFrame, error) {
	if len(records) == 0 {
		return &DataFrame{columns: columns, data: map[string]*Series{}, index: NewRangeIndex(0), shape: [2]int{0, len(columns)}}, nil
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("columns cannot be empty")
	}

	colData := make(map[string][]interface{})
	for _, col := range columns {
		colData[col] = make([]interface{}, 0, len(records))
	}

	for i, row := range records {
		if len(row) != len(columns) {
			return nil, fmt.Errorf("row %d length %d does not match columns length %d", i, len(row), len(columns))
		}
		for j, col := range columns {
			colData[col] = append(colData[col], row[j])
		}
	}

	return New(colData)
}

// Columns returns the column names.
func (df *DataFrame) Columns() []string {
	return df.columns
}

// GetSeries returns the Series for a column name.
func (df *DataFrame) GetSeries(column string) (*Series, bool) {
	series, ok := df.data[column]
	if !ok {
		return nil, false
	}
	return series, true
}

// SetColumn sets or replaces a column with the provided Series.
func (df *DataFrame) SetColumn(name string, series *Series) error {
	if series.Len() != df.shape[0] {
		return fmt.Errorf("series length %d does not match dataframe rows %d", series.Len(), df.shape[0])
	}
	if _, ok := df.data[name]; !ok {
		df.columns = append(df.columns, name)
		df.shape[1] = len(df.columns)
	}
	series.SetName(name)
	df.data[name] = series
	return nil
}

// Index returns the index.
func (df *DataFrame) Index() *Index {
	return df.index
}

// Shape returns the (rows, cols).
func (df *DataFrame) Shape() [2]int {
	return df.shape
}

// Copy returns a shallow copy of the DataFrame.
func (df *DataFrame) Copy() *DataFrame {
	seriesMap := make(map[string]*Series)
	for _, col := range df.columns {
		seriesMap[col] = df.data[col].Copy()
	}
	cols := make([]string, len(df.columns))
	copy(cols, df.columns)
	return &DataFrame{columns: cols, data: seriesMap, index: df.index.Copy(), shape: df.shape}
}

// Head returns the first n rows.
func (df *DataFrame) Head(n int) *DataFrame {
	if n > df.shape[0] {
		n = df.shape[0]
	}
	return df.ILoc(0, n, 0, df.shape[1])
}

// Tail returns the last n rows.
func (df *DataFrame) Tail(n int) *DataFrame {
	if n > df.shape[0] {
		n = df.shape[0]
	}
	start := df.shape[0] - n
	return df.ILoc(start, df.shape[0], 0, df.shape[1])
}

// Select returns a DataFrame with the specified columns.
func (df *DataFrame) Select(columns ...string) *DataFrame {
	seriesMap := make(map[string]*Series)
	cols := make([]string, 0, len(columns))
	for _, col := range columns {
		if s, ok := df.data[col]; ok {
			seriesMap[col] = s.Copy()
			cols = append(cols, col)
		}
	}
	return &DataFrame{columns: cols, data: seriesMap, index: df.index.Copy(), shape: [2]int{df.shape[0], len(cols)}}
}

// At returns a cell value at row index label and column name.
func (df *DataFrame) At(rowLabel interface{}, column string) (interface{}, error) {
	rowPos, err := df.index.GetLoc(rowLabel)
	if err != nil {
		return nil, err
	}
	series, ok := df.data[column]
	if !ok {
		return nil, fmt.Errorf("column '%s' not found", column)
	}
	return series.Get(rowPos)
}

// Row returns a Row by position.
func (df *DataFrame) Row(pos int) (Row, error) {
	if pos < 0 || pos >= df.shape[0] {
		return Row{}, fmt.Errorf("row %d out of range", pos)
	}
	row := make(map[string]interface{})
	for _, col := range df.columns {
		v, _ := df.data[col].Get(pos)
		row[col] = v
	}
	return Row{data: row}, nil
}

// String returns a string representation of the DataFrame.
func (df *DataFrame) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("DataFrame: rows=%d, cols=%d\n", df.shape[0], df.shape[1]))
	if df.shape[1] == 0 {
		return sb.String()
	}

	// Header
	sb.WriteString("index\t")
	for _, col := range df.columns {
		sb.WriteString(col + "\t")
	}
	sb.WriteString("\n")

	maxShow := 5
	rows := df.shape[0]
	if rows <= maxShow*2 {
		for i := 0; i < rows; i++ {
			label, _ := df.index.Get(i)
			sb.WriteString(fmt.Sprintf("%v\t", label))
			for _, col := range df.columns {
				v, _ := df.data[col].Get(i)
				sb.WriteString(fmt.Sprintf("%v\t", v))
			}
			sb.WriteString("\n")
		}
	} else {
		for i := 0; i < maxShow; i++ {
			label, _ := df.index.Get(i)
			sb.WriteString(fmt.Sprintf("%v\t", label))
			for _, col := range df.columns {
				v, _ := df.data[col].Get(i)
				sb.WriteString(fmt.Sprintf("%v\t", v))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("...\n")
		for i := rows - maxShow; i < rows; i++ {
			label, _ := df.index.Get(i)
			sb.WriteString(fmt.Sprintf("%v\t", label))
			for _, col := range df.columns {
				v, _ := df.data[col].Get(i)
				sb.WriteString(fmt.Sprintf("%v\t", v))
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
