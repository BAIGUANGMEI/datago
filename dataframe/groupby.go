package dataframe

import (
	"fmt"
	"sort"
	"sync"
)

// GroupBy represents a grouped DataFrame for aggregation operations
type GroupBy struct {
	df       *DataFrame
	byKeys   []string                    // column names to group by
	groups   map[string][]int            // group key -> row indices
	keyOrder []string                    // maintain order of groups
	mu       sync.RWMutex
}

// GroupByResult represents the result of a groupby aggregation
type GroupByResult struct {
	Keys   [][]interface{}          // group key values
	Values map[string][]interface{} // aggregated values per column
}

// AggFunc defines an aggregation function type
type AggFunc func(*Series) interface{}

// Predefined aggregation functions
var (
	AggSum = func(s *Series) interface{} {
		return s.Sum()
	}
	AggMean = func(s *Series) interface{} {
		return s.Mean()
	}
	AggMin = func(s *Series) interface{} {
		return s.Min()
	}
	AggMax = func(s *Series) interface{} {
		return s.Max()
	}
	AggCount = func(s *Series) interface{} {
		return s.Count()
	}
	AggStd = func(s *Series) interface{} {
		return s.Std()
	}
	AggVar = func(s *Series) interface{} {
		return s.Var()
	}
	AggFirst = func(s *Series) interface{} {
		if s.Len() > 0 {
			v, _ := s.Get(0)
			return v
		}
		return nil
	}
	AggLast = func(s *Series) interface{} {
		if s.Len() > 0 {
			v, _ := s.Get(s.Len() - 1)
			return v
		}
		return nil
	}
)

// GroupBy groups the DataFrame by the specified columns
func (df *DataFrame) GroupBy(columns ...string) (*GroupBy, error) {
	// Validate columns exist
	for _, col := range columns {
		if _, ok := df.data[col]; !ok {
			return nil, fmt.Errorf("column '%s' not found", col)
		}
	}

	gb := &GroupBy{
		df:       df,
		byKeys:   columns,
		groups:   make(map[string][]int),
		keyOrder: make([]string, 0),
	}

	// Build groups
	for i := 0; i < df.shape[0]; i++ {
		key := gb.buildGroupKey(i)
		if _, exists := gb.groups[key]; !exists {
			gb.keyOrder = append(gb.keyOrder, key)
		}
		gb.groups[key] = append(gb.groups[key], i)
	}

	return gb, nil
}

// buildGroupKey creates a unique string key for a row based on grouping columns
func (gb *GroupBy) buildGroupKey(rowIdx int) string {
	key := ""
	for i, col := range gb.byKeys {
		s := gb.df.data[col]
		val, _ := s.Get(rowIdx)
		if i > 0 {
			key += "\x00" // null separator
		}
		key += fmt.Sprintf("%v", val)
	}
	return key
}

// getGroupKeyValues extracts the actual values for a group key
func (gb *GroupBy) getGroupKeyValues(rowIdx int) []interface{} {
	values := make([]interface{}, len(gb.byKeys))
	for i, col := range gb.byKeys {
		s := gb.df.data[col]
		val, _ := s.Get(rowIdx)
		values[i] = val
	}
	return values
}

// NGroups returns the number of groups
func (gb *GroupBy) NGroups() int {
	return len(gb.groups)
}

// Groups returns the group keys and their row indices
func (gb *GroupBy) Groups() map[string][]int {
	return gb.groups
}

// Size returns a Series with the size of each group
func (gb *GroupBy) Size() *DataFrame {
	keyData := make(map[string][]interface{})
	for _, col := range gb.byKeys {
		keyData[col] = make([]interface{}, 0, len(gb.keyOrder))
	}
	sizes := make([]interface{}, 0, len(gb.keyOrder))

	for _, groupKey := range gb.keyOrder {
		indices := gb.groups[groupKey]
		if len(indices) > 0 {
			keyVals := gb.getGroupKeyValues(indices[0])
			for i, col := range gb.byKeys {
				keyData[col] = append(keyData[col], keyVals[i])
			}
			sizes = append(sizes, len(indices))
		}
	}

	// Build result DataFrame
	data := make(map[string][]interface{})
	for col, vals := range keyData {
		data[col] = vals
	}
	data["size"] = sizes

	result, _ := New(data)
	return result
}

// Agg applies multiple aggregation functions to specified columns
func (gb *GroupBy) Agg(aggFuncs map[string][]AggFunc) (*DataFrame, error) {
	// Validate columns
	for col := range aggFuncs {
		if _, ok := gb.df.data[col]; !ok {
			return nil, fmt.Errorf("column '%s' not found", col)
		}
	}

	// Prepare result data
	keyData := make(map[string][]interface{})
	for _, col := range gb.byKeys {
		keyData[col] = make([]interface{}, 0, len(gb.keyOrder))
	}

	aggData := make(map[string][]interface{})
	for col, funcs := range aggFuncs {
		for i := range funcs {
			aggCol := fmt.Sprintf("%s_%d", col, i)
			aggData[aggCol] = make([]interface{}, 0, len(gb.keyOrder))
		}
	}

	// Apply aggregations
	for _, groupKey := range gb.keyOrder {
		indices := gb.groups[groupKey]
		if len(indices) == 0 {
			continue
		}

		// Add key values
		keyVals := gb.getGroupKeyValues(indices[0])
		for i, col := range gb.byKeys {
			keyData[col] = append(keyData[col], keyVals[i])
		}

		// Apply aggregation functions
		for col, funcs := range aggFuncs {
			groupSeries := gb.getGroupSeries(col, indices)
			for i, fn := range funcs {
				aggCol := fmt.Sprintf("%s_%d", col, i)
				aggData[aggCol] = append(aggData[aggCol], fn(groupSeries))
			}
		}
	}

	// Build result DataFrame
	data := make(map[string][]interface{})
	for col, vals := range keyData {
		data[col] = vals
	}
	for col, vals := range aggData {
		data[col] = vals
	}

	return New(data)
}

// Sum computes sum for all numeric columns
func (gb *GroupBy) Sum(columns ...string) *DataFrame {
	return gb.applyAgg(AggSum, "sum", columns...)
}

// Mean computes mean for all numeric columns
func (gb *GroupBy) Mean(columns ...string) *DataFrame {
	return gb.applyAgg(AggMean, "mean", columns...)
}

// Min computes minimum for all numeric columns
func (gb *GroupBy) Min(columns ...string) *DataFrame {
	return gb.applyAgg(AggMin, "min", columns...)
}

// Max computes maximum for all numeric columns
func (gb *GroupBy) Max(columns ...string) *DataFrame {
	return gb.applyAgg(AggMax, "max", columns...)
}

// Count computes count for all columns
func (gb *GroupBy) Count(columns ...string) *DataFrame {
	return gb.applyAgg(AggCount, "count", columns...)
}

// Std computes standard deviation for all numeric columns
func (gb *GroupBy) Std(columns ...string) *DataFrame {
	return gb.applyAgg(AggStd, "std", columns...)
}

// First returns first value in each group
func (gb *GroupBy) First(columns ...string) *DataFrame {
	return gb.applyAgg(AggFirst, "first", columns...)
}

// Last returns last value in each group
func (gb *GroupBy) Last(columns ...string) *DataFrame {
	return gb.applyAgg(AggLast, "last", columns...)
}

// applyAgg applies a single aggregation function to columns
func (gb *GroupBy) applyAgg(aggFunc AggFunc, suffix string, columns ...string) *DataFrame {
	// If no columns specified, use all non-key columns
	if len(columns) == 0 {
		for _, col := range gb.df.columns {
			isKey := false
			for _, key := range gb.byKeys {
				if col == key {
					isKey = true
					break
				}
			}
			if !isKey {
				columns = append(columns, col)
			}
		}
	}

	// Prepare result data
	keyData := make(map[string][]interface{})
	for _, col := range gb.byKeys {
		keyData[col] = make([]interface{}, 0, len(gb.keyOrder))
	}

	aggData := make(map[string][]interface{})
	for _, col := range columns {
		aggData[col+"_"+suffix] = make([]interface{}, 0, len(gb.keyOrder))
	}

	// Apply aggregation
	for _, groupKey := range gb.keyOrder {
		indices := gb.groups[groupKey]
		if len(indices) == 0 {
			continue
		}

		// Add key values
		keyVals := gb.getGroupKeyValues(indices[0])
		for i, col := range gb.byKeys {
			keyData[col] = append(keyData[col], keyVals[i])
		}

		// Apply aggregation
		for _, col := range columns {
			if _, ok := gb.df.data[col]; !ok {
				continue
			}
			groupSeries := gb.getGroupSeries(col, indices)
			aggData[col+"_"+suffix] = append(aggData[col+"_"+suffix], aggFunc(groupSeries))
		}
	}

	// Build result DataFrame
	data := make(map[string][]interface{})
	for col, vals := range keyData {
		data[col] = vals
	}
	for col, vals := range aggData {
		data[col] = vals
	}

	result, _ := New(data)
	return result
}

// getGroupSeries extracts a Series for a specific group
func (gb *GroupBy) getGroupSeries(col string, indices []int) *Series {
	s := gb.df.data[col]
	groupData := make([]interface{}, len(indices))
	for i, idx := range indices {
		groupData[i], _ = s.Get(idx)
	}
	return NewSeries(groupData, col)
}

// Apply applies a custom function to each group
func (gb *GroupBy) Apply(fn func(*DataFrame) *DataFrame) *DataFrame {
	var results []*DataFrame

	for _, groupKey := range gb.keyOrder {
		indices := gb.groups[groupKey]
		if len(indices) == 0 {
			continue
		}

		// Create group DataFrame
		groupDF := gb.getGroupDataFrame(indices)

		// Apply function
		result := fn(groupDF)
		if result != nil && result.shape[0] > 0 {
			results = append(results, result)
		}
	}

	// Concatenate results
	if len(results) == 0 {
		return &DataFrame{columns: []string{}, data: map[string]*Series{}, index: NewRangeIndex(0), shape: [2]int{0, 0}}
	}

	return Concat(results...)
}

// getGroupDataFrame extracts a DataFrame for a specific group
func (gb *GroupBy) getGroupDataFrame(indices []int) *DataFrame {
	seriesMap := make(map[string]*Series)
	for _, col := range gb.df.columns {
		s := gb.df.data[col]
		groupData := make([]interface{}, len(indices))
		for i, idx := range indices {
			groupData[i], _ = s.Get(idx)
		}
		seriesMap[col] = NewSeries(groupData, col)
	}

	cols := make([]string, len(gb.df.columns))
	copy(cols, gb.df.columns)

	return &DataFrame{
		columns: cols,
		data:    seriesMap,
		index:   NewRangeIndex(len(indices)),
		shape:   [2]int{len(indices), len(cols)},
	}
}

// Filter filters groups based on a predicate
func (gb *GroupBy) Filter(predicate func(*DataFrame) bool) *DataFrame {
	var allIndices []int

	for _, groupKey := range gb.keyOrder {
		indices := gb.groups[groupKey]
		if len(indices) == 0 {
			continue
		}

		groupDF := gb.getGroupDataFrame(indices)
		if predicate(groupDF) {
			allIndices = append(allIndices, indices...)
		}
	}

	if len(allIndices) == 0 {
		return &DataFrame{columns: gb.df.columns, data: map[string]*Series{}, index: NewRangeIndex(0), shape: [2]int{0, len(gb.df.columns)}}
	}

	// Sort indices to maintain order
	sort.Ints(allIndices)

	// Build result DataFrame
	seriesMap := make(map[string]*Series)
	for _, col := range gb.df.columns {
		s := gb.df.data[col]
		newData := make([]interface{}, len(allIndices))
		for i, idx := range allIndices {
			newData[i], _ = s.Get(idx)
		}
		seriesMap[col] = NewSeries(newData, col)
	}

	cols := make([]string, len(gb.df.columns))
	copy(cols, gb.df.columns)

	return &DataFrame{
		columns: cols,
		data:    seriesMap,
		index:   NewRangeIndex(len(allIndices)),
		shape:   [2]int{len(allIndices), len(cols)},
	}
}

// Transform applies a function to each group and returns result with original index
func (gb *GroupBy) Transform(col string, fn func(*Series) *Series) (*Series, error) {
	if _, ok := gb.df.data[col]; !ok {
		return nil, fmt.Errorf("column '%s' not found", col)
	}

	result := make([]interface{}, gb.df.shape[0])

	for _, indices := range gb.groups {
		if len(indices) == 0 {
			continue
		}

		groupSeries := gb.getGroupSeries(col, indices)
		transformed := fn(groupSeries)

		for i, idx := range indices {
			if i < transformed.Len() {
				val, _ := transformed.Get(i)
				result[idx] = val
			}
		}
	}

	return NewSeries(result, col+"_transformed"), nil
}

// Concat concatenates multiple DataFrames vertically
func Concat(dfs ...*DataFrame) *DataFrame {
	if len(dfs) == 0 {
		return &DataFrame{columns: []string{}, data: map[string]*Series{}, index: NewRangeIndex(0), shape: [2]int{0, 0}}
	}

	// Use first DataFrame's columns as reference
	cols := make([]string, len(dfs[0].columns))
	copy(cols, dfs[0].columns)

	// Collect all data
	colData := make(map[string][]interface{})
	for _, col := range cols {
		colData[col] = []interface{}{}
	}

	totalRows := 0
	for _, df := range dfs {
		for _, col := range cols {
			if s, ok := df.data[col]; ok {
				colData[col] = append(colData[col], s.data...)
			} else {
				// Fill with nil if column doesn't exist
				for i := 0; i < df.shape[0]; i++ {
					colData[col] = append(colData[col], nil)
				}
			}
		}
		totalRows += df.shape[0]
	}

	// Build result DataFrame
	seriesMap := make(map[string]*Series)
	for _, col := range cols {
		seriesMap[col] = NewSeries(colData[col], col)
	}

	return &DataFrame{
		columns: cols,
		data:    seriesMap,
		index:   NewRangeIndex(totalRows),
		shape:   [2]int{totalRows, len(cols)},
	}
}
