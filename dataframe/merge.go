package dataframe

import (
	"fmt"
)

// JoinType defines the type of join operation
type JoinType int

const (
	// InnerJoin returns only matching rows from both DataFrames
	InnerJoin JoinType = iota
	// LeftJoin returns all rows from left DataFrame and matching rows from right
	LeftJoin
	// RightJoin returns all rows from right DataFrame and matching rows from left
	RightJoin
	// OuterJoin returns all rows from both DataFrames
	OuterJoin
)

// String returns the string representation of JoinType
func (j JoinType) String() string {
	switch j {
	case InnerJoin:
		return "inner"
	case LeftJoin:
		return "left"
	case RightJoin:
		return "right"
	case OuterJoin:
		return "outer"
	default:
		return "unknown"
	}
}

// MergeOptions defines options for merge operations
type MergeOptions struct {
	How         JoinType // type of join
	On          []string // columns to join on (same name in both DataFrames)
	LeftOn      []string // columns to join on from left DataFrame
	RightOn     []string // columns to join on from right DataFrame
	Suffixes    [2]string // suffixes to use for overlapping columns
	Indicator   bool      // add _merge column indicating source
}

// DefaultMergeOptions returns default merge options
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		How:      InnerJoin,
		Suffixes: [2]string{"_x", "_y"},
	}
}

// Merge merges two DataFrames based on common columns or specified keys
func Merge(left, right *DataFrame, opts MergeOptions) (*DataFrame, error) {
	if left == nil || right == nil {
		return nil, fmt.Errorf("both DataFrames must be non-nil")
	}

	// Determine join keys
	leftKeys, rightKeys, err := resolveJoinKeys(left, right, opts)
	if err != nil {
		return nil, err
	}

	// Build index for right DataFrame
	rightIndex := buildJoinIndex(right, rightKeys)

	// Perform join based on type
	switch opts.How {
	case InnerJoin:
		return innerJoin(left, right, leftKeys, rightKeys, rightIndex, opts)
	case LeftJoin:
		return leftJoin(left, right, leftKeys, rightKeys, rightIndex, opts)
	case RightJoin:
		return rightJoin(left, right, leftKeys, rightKeys, rightIndex, opts)
	case OuterJoin:
		return outerJoin(left, right, leftKeys, rightKeys, rightIndex, opts)
	default:
		return nil, fmt.Errorf("unknown join type: %v", opts.How)
	}
}

// resolveJoinKeys determines the columns to join on
func resolveJoinKeys(left, right *DataFrame, opts MergeOptions) ([]string, []string, error) {
	var leftKeys, rightKeys []string

	if len(opts.On) > 0 {
		// Same column names in both DataFrames
		for _, col := range opts.On {
			if _, ok := left.data[col]; !ok {
				return nil, nil, fmt.Errorf("column '%s' not found in left DataFrame", col)
			}
			if _, ok := right.data[col]; !ok {
				return nil, nil, fmt.Errorf("column '%s' not found in right DataFrame", col)
			}
		}
		leftKeys = opts.On
		rightKeys = opts.On
	} else if len(opts.LeftOn) > 0 && len(opts.RightOn) > 0 {
		// Different column names
		if len(opts.LeftOn) != len(opts.RightOn) {
			return nil, nil, fmt.Errorf("LeftOn and RightOn must have same length")
		}
		for _, col := range opts.LeftOn {
			if _, ok := left.data[col]; !ok {
				return nil, nil, fmt.Errorf("column '%s' not found in left DataFrame", col)
			}
		}
		for _, col := range opts.RightOn {
			if _, ok := right.data[col]; !ok {
				return nil, nil, fmt.Errorf("column '%s' not found in right DataFrame", col)
			}
		}
		leftKeys = opts.LeftOn
		rightKeys = opts.RightOn
	} else {
		// Auto-detect common columns
		commonCols := findCommonColumns(left, right)
		if len(commonCols) == 0 {
			return nil, nil, fmt.Errorf("no common columns found and no join keys specified")
		}
		leftKeys = commonCols
		rightKeys = commonCols
	}

	return leftKeys, rightKeys, nil
}

// findCommonColumns finds columns present in both DataFrames
func findCommonColumns(left, right *DataFrame) []string {
	var common []string
	for _, col := range left.columns {
		if _, ok := right.data[col]; ok {
			common = append(common, col)
		}
	}
	return common
}

// buildJoinIndex builds a hash index for join operations
func buildJoinIndex(df *DataFrame, keys []string) map[string][]int {
	index := make(map[string][]int)
	for i := 0; i < df.shape[0]; i++ {
		key := buildRowKey(df, keys, i)
		index[key] = append(index[key], i)
	}
	return index
}

// buildRowKey creates a unique string key for a row based on specified columns
func buildRowKey(df *DataFrame, keys []string, rowIdx int) string {
	key := ""
	for i, col := range keys {
		s := df.data[col]
		val, _ := s.Get(rowIdx)
		if i > 0 {
			key += "\x00"
		}
		key += fmt.Sprintf("%v", val)
	}
	return key
}

// innerJoin performs an inner join
func innerJoin(left, right *DataFrame, leftKeys, rightKeys []string, rightIndex map[string][]int, opts MergeOptions) (*DataFrame, error) {
	resultCols, colMapping := prepareResultColumns(left, right, leftKeys, rightKeys, opts)
	resultData := initResultData(resultCols)
	var indicators []interface{}

	for i := 0; i < left.shape[0]; i++ {
		leftKey := buildRowKey(left, leftKeys, i)
		if rightRows, ok := rightIndex[leftKey]; ok {
			for _, rightRow := range rightRows {
				appendJoinedRow(resultData, colMapping, left, right, i, rightRow, leftKeys, rightKeys, opts)
				if opts.Indicator {
					indicators = append(indicators, "both")
				}
			}
		}
	}

	return buildJoinResult(resultCols, resultData, indicators, opts)
}

// leftJoin performs a left join
func leftJoin(left, right *DataFrame, leftKeys, rightKeys []string, rightIndex map[string][]int, opts MergeOptions) (*DataFrame, error) {
	resultCols, colMapping := prepareResultColumns(left, right, leftKeys, rightKeys, opts)
	resultData := initResultData(resultCols)
	var indicators []interface{}

	for i := 0; i < left.shape[0]; i++ {
		leftKey := buildRowKey(left, leftKeys, i)
		if rightRows, ok := rightIndex[leftKey]; ok {
			for _, rightRow := range rightRows {
				appendJoinedRow(resultData, colMapping, left, right, i, rightRow, leftKeys, rightKeys, opts)
				if opts.Indicator {
					indicators = append(indicators, "both")
				}
			}
		} else {
			// No match - include left row with nulls for right
			appendLeftOnlyRow(resultData, colMapping, left, right, i, leftKeys, rightKeys, opts)
			if opts.Indicator {
				indicators = append(indicators, "left_only")
			}
		}
	}

	return buildJoinResult(resultCols, resultData, indicators, opts)
}

// rightJoin performs a right join
func rightJoin(left, right *DataFrame, leftKeys, rightKeys []string, rightIndex map[string][]int, opts MergeOptions) (*DataFrame, error) {
	// Build left index
	leftIndex := buildJoinIndex(left, leftKeys)

	resultCols, colMapping := prepareResultColumns(left, right, leftKeys, rightKeys, opts)
	resultData := initResultData(resultCols)
	var indicators []interface{}

	for i := 0; i < right.shape[0]; i++ {
		rightKey := buildRowKey(right, rightKeys, i)
		if leftRows, ok := leftIndex[rightKey]; ok {
			for _, leftRow := range leftRows {
				appendJoinedRow(resultData, colMapping, left, right, leftRow, i, leftKeys, rightKeys, opts)
				if opts.Indicator {
					indicators = append(indicators, "both")
				}
			}
		} else {
			// No match - include right row with nulls for left
			appendRightOnlyRow(resultData, colMapping, left, right, i, leftKeys, rightKeys, opts)
			if opts.Indicator {
				indicators = append(indicators, "right_only")
			}
		}
	}

	return buildJoinResult(resultCols, resultData, indicators, opts)
}

// outerJoin performs a full outer join
func outerJoin(left, right *DataFrame, leftKeys, rightKeys []string, rightIndex map[string][]int, opts MergeOptions) (*DataFrame, error) {
	resultCols, colMapping := prepareResultColumns(left, right, leftKeys, rightKeys, opts)
	resultData := initResultData(resultCols)
	var indicators []interface{}

	// Track which right rows have been matched
	matchedRight := make(map[int]bool)

	// Process all left rows
	for i := 0; i < left.shape[0]; i++ {
		leftKey := buildRowKey(left, leftKeys, i)
		if rightRows, ok := rightIndex[leftKey]; ok {
			for _, rightRow := range rightRows {
				appendJoinedRow(resultData, colMapping, left, right, i, rightRow, leftKeys, rightKeys, opts)
				matchedRight[rightRow] = true
				if opts.Indicator {
					indicators = append(indicators, "both")
				}
			}
		} else {
			appendLeftOnlyRow(resultData, colMapping, left, right, i, leftKeys, rightKeys, opts)
			if opts.Indicator {
				indicators = append(indicators, "left_only")
			}
		}
	}

	// Add unmatched right rows
	for i := 0; i < right.shape[0]; i++ {
		if !matchedRight[i] {
			appendRightOnlyRow(resultData, colMapping, left, right, i, leftKeys, rightKeys, opts)
			if opts.Indicator {
				indicators = append(indicators, "right_only")
			}
		}
	}

	return buildJoinResult(resultCols, resultData, indicators, opts)
}

// columnMapping stores information about how to map columns in the result
type columnMapping struct {
	source    string // "left", "right", or "key"
	srcCol    string // original column name
	isKey     bool   // whether this is a join key column
	keyIndex  int    // index in keys array (if isKey)
}

// prepareResultColumns determines the columns in the result DataFrame
func prepareResultColumns(left, right *DataFrame, leftKeys, rightKeys []string, opts MergeOptions) ([]string, map[string]columnMapping) {
	var resultCols []string
	colMapping := make(map[string]columnMapping)

	// Track right columns that are join keys
	rightKeySet := make(map[string]bool)
	for _, k := range rightKeys {
		rightKeySet[k] = true
	}

	// Track left columns that are join keys
	leftKeySet := make(map[string]bool)
	for _, k := range leftKeys {
		leftKeySet[k] = true
	}

	// Add left columns
	for _, col := range left.columns {
		resultCol := col
		_, inRight := right.data[col]
		isLeftKey := leftKeySet[col]
		
		if inRight && !isLeftKey {
			// Overlapping column, not a key - add suffix
			resultCol = col + opts.Suffixes[0]
		}
		
		resultCols = append(resultCols, resultCol)
		if isLeftKey {
			for i, k := range leftKeys {
				if k == col {
					colMapping[resultCol] = columnMapping{source: "key", srcCol: col, isKey: true, keyIndex: i}
					break
				}
			}
		} else {
			colMapping[resultCol] = columnMapping{source: "left", srcCol: col}
		}
	}

	// Add right columns (excluding join keys with same name)
	for _, col := range right.columns {
		// Skip if it's a key column with same name as left key
		isRightKey := rightKeySet[col]
		if isRightKey {
			// Check if there's a corresponding left key with same name
			for i, rk := range rightKeys {
				if rk == col && leftKeys[i] == col {
					// Same column name used as key on both sides - already included
					continue
				}
			}
		}

		resultCol := col
		_, inLeft := left.data[col]
		
		if inLeft {
			if isRightKey {
				// It's a key but with different name - skip as we use left key
				continue
			}
			// Overlapping column - add suffix
			resultCol = col + opts.Suffixes[1]
		}

		// Check if already added (for key columns)
		alreadyAdded := false
		for _, rc := range resultCols {
			if rc == resultCol {
				alreadyAdded = true
				break
			}
		}
		if alreadyAdded {
			continue
		}

		resultCols = append(resultCols, resultCol)
		colMapping[resultCol] = columnMapping{source: "right", srcCol: col}
	}

	return resultCols, colMapping
}

// initResultData initializes the result data structure
func initResultData(cols []string) map[string][]interface{} {
	data := make(map[string][]interface{})
	for _, col := range cols {
		data[col] = []interface{}{}
	}
	return data
}

// appendJoinedRow adds a row from both DataFrames
func appendJoinedRow(resultData map[string][]interface{}, colMapping map[string]columnMapping, 
	left, right *DataFrame, leftRow, rightRow int, leftKeys, rightKeys []string, opts MergeOptions) {
	
	for resultCol, mapping := range colMapping {
		var val interface{}
		if mapping.isKey {
			// Use left key value for key columns
			s := left.data[leftKeys[mapping.keyIndex]]
			val, _ = s.Get(leftRow)
		} else if mapping.source == "left" {
			s := left.data[mapping.srcCol]
			val, _ = s.Get(leftRow)
		} else {
			s := right.data[mapping.srcCol]
			val, _ = s.Get(rightRow)
		}
		resultData[resultCol] = append(resultData[resultCol], val)
	}
}

// appendLeftOnlyRow adds a row from left DataFrame with nulls for right
func appendLeftOnlyRow(resultData map[string][]interface{}, colMapping map[string]columnMapping,
	left, right *DataFrame, leftRow int, leftKeys, rightKeys []string, opts MergeOptions) {
	
	for resultCol, mapping := range colMapping {
		var val interface{}
		if mapping.isKey {
			s := left.data[leftKeys[mapping.keyIndex]]
			val, _ = s.Get(leftRow)
		} else if mapping.source == "left" {
			s := left.data[mapping.srcCol]
			val, _ = s.Get(leftRow)
		} else {
			val = nil
		}
		resultData[resultCol] = append(resultData[resultCol], val)
	}
}

// appendRightOnlyRow adds a row from right DataFrame with nulls for left
func appendRightOnlyRow(resultData map[string][]interface{}, colMapping map[string]columnMapping,
	left, right *DataFrame, rightRow int, leftKeys, rightKeys []string, opts MergeOptions) {
	
	for resultCol, mapping := range colMapping {
		var val interface{}
		if mapping.isKey {
			// Use right key value for key columns
			s := right.data[rightKeys[mapping.keyIndex]]
			val, _ = s.Get(rightRow)
		} else if mapping.source == "left" {
			val = nil
		} else {
			s := right.data[mapping.srcCol]
			val, _ = s.Get(rightRow)
		}
		resultData[resultCol] = append(resultData[resultCol], val)
	}
}

// buildJoinResult builds the final DataFrame from join results
func buildJoinResult(cols []string, data map[string][]interface{}, indicators []interface{}, opts MergeOptions) (*DataFrame, error) {
	if opts.Indicator {
		cols = append(cols, "_merge")
		data["_merge"] = indicators
	}

	seriesMap := make(map[string]*Series)
	rowCount := 0
	for _, col := range cols {
		seriesMap[col] = NewSeries(data[col], col)
		if rowCount == 0 {
			rowCount = len(data[col])
		}
	}

	return &DataFrame{
		columns: cols,
		data:    seriesMap,
		index:   NewRangeIndex(rowCount),
		shape:   [2]int{rowCount, len(cols)},
	}, nil
}

// Join is a convenience method for joining DataFrames
func (df *DataFrame) Join(other *DataFrame, on []string, how JoinType) (*DataFrame, error) {
	opts := DefaultMergeOptions()
	opts.On = on
	opts.How = how
	return Merge(df, other, opts)
}

// MergeOn is a convenience method for merging on specified columns
func (df *DataFrame) MergeOn(other *DataFrame, leftOn, rightOn []string, how JoinType) (*DataFrame, error) {
	opts := DefaultMergeOptions()
	opts.LeftOn = leftOn
	opts.RightOn = rightOn
	opts.How = how
	return Merge(df, other, opts)
}
