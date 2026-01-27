package dataframe

import (
	"fmt"
	"sort"
)

// FilterFunc defines a filter function for rows.
type FilterFunc func(row Row) bool

// SortOrder defines ascending or descending order.
type SortOrder int

const (
	// Ascending sort order
	Ascending SortOrder = iota
	// Descending sort order
	Descending
)

// ILoc selects rows and columns by integer position.
func (df *DataFrame) ILoc(rowStart, rowEnd, colStart, colEnd int) *DataFrame {
	if rowStart < 0 {
		rowStart = 0
	}
	if rowEnd > df.shape[0] {
		rowEnd = df.shape[0]
	}
	if colStart < 0 {
		colStart = 0
	}
	if colEnd > df.shape[1] {
		colEnd = df.shape[1]
	}

	cols := df.columns[colStart:colEnd]
	seriesMap := make(map[string]*Series)
	for _, col := range cols {
		s := df.data[col]
		seriesMap[col] = s.Slice(rowStart, rowEnd)
	}

	newIndex := df.index.Slice(rowStart, rowEnd)
	return &DataFrame{columns: append([]string{}, cols...), data: seriesMap, index: newIndex, shape: [2]int{rowEnd - rowStart, colEnd - colStart}}
}

// Loc selects rows and columns by labels.
func (df *DataFrame) Loc(rowLabels interface{}, colLabels interface{}) *DataFrame {
	// For simplicity: rowLabels can be []interface{} or nil; colLabels can be []string or nil
	var rowPositions []int
	switch v := rowLabels.(type) {
	case nil:
		for i := 0; i < df.shape[0]; i++ {
			rowPositions = append(rowPositions, i)
		}
	case []interface{}:
		for _, label := range v {
			pos, err := df.index.GetLoc(label)
			if err == nil {
				rowPositions = append(rowPositions, pos)
			}
		}
	default:
		if pos, err := df.index.GetLoc(v); err == nil {
			rowPositions = []int{pos}
		}
	}

	var cols []string
	switch v := colLabels.(type) {
	case nil:
		cols = append([]string{}, df.columns...)
	case []string:
		cols = v
	case string:
		cols = []string{v}
	default:
		cols = append([]string{}, df.columns...)
	}

	seriesMap := make(map[string]*Series)
	for _, col := range cols {
		s, ok := df.data[col]
		if !ok {
			continue
		}
		newData := make([]interface{}, len(rowPositions))
		newLabels := make([]interface{}, len(rowPositions))
		for i, pos := range rowPositions {
			newData[i] = s.data[pos]
			label, _ := df.index.Get(pos)
			newLabels[i] = label
		}
		seriesMap[col] = NewSeriesWithIndex(newData, col, NewIndex(newLabels, df.index.Name()))
	}

	return &DataFrame{
		columns: cols,
		data:    seriesMap,
		index:   NewIndex(extractLabels(df.index, rowPositions), df.index.Name()),
		shape:   [2]int{len(rowPositions), len(cols)},
	}
}

func extractLabels(idx *Index, positions []int) []interface{} {
	labels := make([]interface{}, len(positions))
	for i, pos := range positions {
		label, _ := idx.Get(pos)
		labels[i] = label
	}
	return labels
}

// Filter filters rows using the provided function.
func (df *DataFrame) Filter(fn FilterFunc) *DataFrame {
	var rows []int
	for i := 0; i < df.shape[0]; i++ {
		row, _ := df.Row(i)
		if fn(row) {
			rows = append(rows, i)
		}
	}
	return df.Loc(extractLabels(df.index, rows), nil)
}

// AddColumn adds a new column to the DataFrame.
func (df *DataFrame) AddColumn(name string, series *Series) *DataFrame {
	if series.Len() != df.shape[0] {
		return df
	}
	newDF := df.Copy()
	newDF.columns = append(newDF.columns, name)
	newDF.data[name] = series.Copy()
	newDF.shape[1] = len(newDF.columns)
	return newDF
}

// Drop removes columns from the DataFrame.
func (df *DataFrame) Drop(columns ...string) *DataFrame {
	toDrop := make(map[string]bool)
	for _, col := range columns {
		toDrop[col] = true
	}
	newCols := make([]string, 0, len(df.columns))
	newData := make(map[string]*Series)
	for _, col := range df.columns {
		if !toDrop[col] {
			newCols = append(newCols, col)
			newData[col] = df.data[col].Copy()
		}
	}
	return &DataFrame{columns: newCols, data: newData, index: df.index.Copy(), shape: [2]int{df.shape[0], len(newCols)}}
}

// Rename renames columns according to the mapping.
func (df *DataFrame) Rename(mapping map[string]string) *DataFrame {
	newCols := make([]string, len(df.columns))
	newData := make(map[string]*Series)
	for i, col := range df.columns {
		newCol := col
		if v, ok := mapping[col]; ok {
			newCol = v
		}
		newCols[i] = newCol
		newData[newCol] = df.data[col].Copy()
		newData[newCol].SetName(newCol)
	}
	return &DataFrame{columns: newCols, data: newData, index: df.index.Copy(), shape: [2]int{df.shape[0], len(newCols)}}
}

// SortBy sorts the DataFrame by a column.
func (df *DataFrame) SortBy(column string, order SortOrder) *DataFrame {
	s, ok := df.data[column]
	if !ok {
		return df
	}

	type indexedValue struct {
		index int
		value interface{}
	}
	indexed := make([]indexedValue, df.shape[0])
	for i := 0; i < df.shape[0]; i++ {
		indexed[i] = indexedValue{i, s.data[i]}
	}

	sort.Slice(indexed, func(i, j int) bool {
		vi, vj := indexed[i].value, indexed[j].value
		if vi == nil && vj == nil {
			return false
		}
		if vi == nil {
			return order == Descending
		}
		if vj == nil {
			return order == Ascending
		}
		fi, erri := toFloat64(vi)
		fj, errj := toFloat64(vj)
		if erri == nil && errj == nil {
			if order == Ascending {
				return fi < fj
			}
			return fi > fj
		}
		if order == Ascending {
			return fmt.Sprintf("%v", vi) < fmt.Sprintf("%v", vj)
		}
		return fmt.Sprintf("%v", vi) > fmt.Sprintf("%v", vj)
	})

	newDF := df.Copy()
	newIndexLabels := make([]interface{}, df.shape[0])
	for i, iv := range indexed {
		label, _ := df.index.Get(iv.index)
		newIndexLabels[i] = label
	}
	newDF.index = NewIndex(newIndexLabels, df.index.Name())
	for _, col := range df.columns {
		newData := make([]interface{}, df.shape[0])
		for i, iv := range indexed {
			newData[i] = df.data[col].data[iv.index]
		}
		newDF.data[col] = NewSeriesWithIndex(newData, col, newDF.index)
	}
	return newDF
}

// Describe returns a statistical summary of numeric columns.
func (df *DataFrame) Describe() *DataFrame {
	stats := []string{"count", "mean", "std", "min", "max"}
	colData := make(map[string][]interface{})
	for _, stat := range stats {
		colData[stat] = make([]interface{}, 0, len(df.columns))
	}

	var statIndex []interface{}
	for _, col := range df.columns {
		s := df.data[col]
		count := float64(s.Count())
		mean := s.Mean()
		std := s.Std()
		min := s.Min()
		max := s.Max()

		colData["count"] = append(colData["count"], count)
		colData["mean"] = append(colData["mean"], mean)
		colData["std"] = append(colData["std"], std)
		colData["min"] = append(colData["min"], min)
		colData["max"] = append(colData["max"], max)
		statIndex = append(statIndex, col)
	}

	dfStats, _ := New(colData)
	dfStats.index = NewIndex(statIndex, "column")
	return dfStats
}
