package dataframe

import (
	"runtime"
	"sync"
)

// ParallelOptions defines options for parallel operations
type ParallelOptions struct {
	NumWorkers int  // number of goroutines to use (0 = auto)
	ChunkSize  int  // minimum chunk size per worker
}

// DefaultParallelOptions returns default parallel options
func DefaultParallelOptions() ParallelOptions {
	return ParallelOptions{
		NumWorkers: 0,  // auto-detect
		ChunkSize:  1000,
	}
}

// getNumWorkers returns the number of workers to use
func getNumWorkers(opts ParallelOptions, dataSize int) int {
	if opts.NumWorkers > 0 {
		return opts.NumWorkers
	}
	// Use number of CPUs, but limit based on data size
	numCPU := runtime.NumCPU()
	minChunk := opts.ChunkSize
	if minChunk <= 0 {
		minChunk = 1000
	}
	maxWorkers := (dataSize + minChunk - 1) / minChunk
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	if numCPU < maxWorkers {
		return numCPU
	}
	return maxWorkers
}

// ParallelApply applies a function to each element of a Series in parallel
func (s *Series) ParallelApply(fn func(interface{}) interface{}, opts ...ParallelOptions) *Series {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	n := s.Len()
	if n == 0 {
		return NewSeries([]interface{}{}, s.name)
	}

	numWorkers := getNumWorkers(opt, n)
	if numWorkers <= 1 {
		return s.Apply(fn)
	}

	result := make([]interface{}, n)
	chunkSize := (n + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		if start >= n {
			wg.Done()
			continue
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				result[i] = fn(s.data[i])
			}
		}(start, end)
	}

	wg.Wait()

	return &Series{
		name:  s.name,
		data:  result,
		dtype: DTypeObject,
		index: s.index.Copy(),
	}
}

// ParallelFilter filters the DataFrame using parallel processing
func (df *DataFrame) ParallelFilter(fn FilterFunc, opts ...ParallelOptions) *DataFrame {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	n := df.shape[0]
	if n == 0 {
		return df.Copy()
	}

	numWorkers := getNumWorkers(opt, n)
	if numWorkers <= 1 {
		return df.Filter(fn)
	}

	// Each worker collects matching indices
	type result struct {
		indices []int
	}

	chunkSize := (n + numWorkers - 1) / numWorkers
	results := make([]result, numWorkers)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		if start >= n {
			wg.Done()
			continue
		}

		go func(w, start, end int) {
			defer wg.Done()
			var indices []int
			for i := start; i < end; i++ {
				row, _ := df.Row(i)
				if fn(row) {
					indices = append(indices, i)
				}
			}
			results[w].indices = indices
		}(w, start, end)
	}

	wg.Wait()

	// Collect all matching indices
	var allIndices []int
	for _, r := range results {
		allIndices = append(allIndices, r.indices...)
	}

	if len(allIndices) == 0 {
		return &DataFrame{
			columns: df.columns,
			data:    map[string]*Series{},
			index:   NewRangeIndex(0),
			shape:   [2]int{0, len(df.columns)},
		}
	}

	// Build result DataFrame
	seriesMap := make(map[string]*Series)
	for _, col := range df.columns {
		s := df.data[col]
		newData := make([]interface{}, len(allIndices))
		for i, idx := range allIndices {
			newData[i], _ = s.Get(idx)
		}
		seriesMap[col] = NewSeries(newData, col)
	}

	cols := make([]string, len(df.columns))
	copy(cols, df.columns)

	return &DataFrame{
		columns: cols,
		data:    seriesMap,
		index:   NewRangeIndex(len(allIndices)),
		shape:   [2]int{len(allIndices), len(cols)},
	}
}

// ParallelTransform applies a transformation function to each column in parallel
func (df *DataFrame) ParallelTransform(fn func(*Series) *Series, opts ...ParallelOptions) *DataFrame {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	numCols := len(df.columns)
	if numCols == 0 {
		return df.Copy()
	}

	numWorkers := getNumWorkers(opt, numCols)
	if numWorkers > numCols {
		numWorkers = numCols
	}

	resultSeries := make(map[string]*Series)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Create channel for columns to process
	colChan := make(chan string, numCols)
	for _, col := range df.columns {
		colChan <- col
	}
	close(colChan)

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for col := range colChan {
				s := df.data[col]
				transformed := fn(s)
				mu.Lock()
				resultSeries[col] = transformed
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	cols := make([]string, len(df.columns))
	copy(cols, df.columns)

	return &DataFrame{
		columns: cols,
		data:    resultSeries,
		index:   df.index.Copy(),
		shape:   df.shape,
	}
}

// ParallelSum computes sum for all numeric columns in parallel
func (df *DataFrame) ParallelSum(opts ...ParallelOptions) map[string]float64 {
	return df.parallelAggFloat64(func(s *Series) float64 { return s.Sum() }, opts...)
}

// ParallelMean computes mean for all numeric columns in parallel
func (df *DataFrame) ParallelMean(opts ...ParallelOptions) map[string]float64 {
	return df.parallelAggFloat64(func(s *Series) float64 { return s.Mean() }, opts...)
}

// ParallelMin computes minimum for all numeric columns in parallel
func (df *DataFrame) ParallelMin(opts ...ParallelOptions) map[string]interface{} {
	return df.parallelAggInterface(func(s *Series) interface{} { return s.Min() }, opts...)
}

// ParallelMax computes maximum for all numeric columns in parallel
func (df *DataFrame) ParallelMax(opts ...ParallelOptions) map[string]interface{} {
	return df.parallelAggInterface(func(s *Series) interface{} { return s.Max() }, opts...)
}

// parallelAggFloat64 applies an aggregation function to all columns in parallel
func (df *DataFrame) parallelAggFloat64(fn func(*Series) float64, opts ...ParallelOptions) map[string]float64 {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	numCols := len(df.columns)
	if numCols == 0 {
		return make(map[string]float64)
	}

	numWorkers := getNumWorkers(opt, numCols)
	if numWorkers > numCols {
		numWorkers = numCols
	}

	result := make(map[string]float64)
	var mu sync.Mutex
	var wg sync.WaitGroup

	colChan := make(chan string, numCols)
	for _, col := range df.columns {
		colChan <- col
	}
	close(colChan)

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for col := range colChan {
				s := df.data[col]
				val := fn(s)
				mu.Lock()
				result[col] = val
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return result
}

// parallelAggInterface applies an aggregation function returning interface{} to all columns in parallel
func (df *DataFrame) parallelAggInterface(fn func(*Series) interface{}, opts ...ParallelOptions) map[string]interface{} {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	numCols := len(df.columns)
	if numCols == 0 {
		return make(map[string]interface{})
	}

	numWorkers := getNumWorkers(opt, numCols)
	if numWorkers > numCols {
		numWorkers = numCols
	}

	result := make(map[string]interface{})
	var mu sync.Mutex
	var wg sync.WaitGroup

	colChan := make(chan string, numCols)
	for _, col := range df.columns {
		colChan <- col
	}
	close(colChan)

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for col := range colChan {
				s := df.data[col]
				val := fn(s)
				mu.Lock()
				result[col] = val
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return result
}

// ParallelGroupByAgg performs parallel aggregation on grouped data
func (gb *GroupBy) ParallelAgg(aggFuncs map[string][]AggFunc, opts ...ParallelOptions) (*DataFrame, error) {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	// Validate columns
	for col := range aggFuncs {
		if _, ok := gb.df.data[col]; !ok {
			return nil, nil
		}
	}

	numGroups := len(gb.keyOrder)
	if numGroups == 0 {
		return New(map[string][]interface{}{})
	}

	numWorkers := getNumWorkers(opt, numGroups)
	if numWorkers > numGroups {
		numWorkers = numGroups
	}

	// Process groups in parallel
	type groupResult struct {
		keyVals []interface{}
		aggVals map[string][]interface{}
	}

	results := make([]groupResult, numGroups)
	chunkSize := (numGroups + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > numGroups {
			end = numGroups
		}
		if start >= numGroups {
			wg.Done()
			continue
		}

		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				groupKey := gb.keyOrder[i]
				indices := gb.groups[groupKey]
				if len(indices) == 0 {
					continue
				}

				keyVals := gb.getGroupKeyValues(indices[0])
				aggVals := make(map[string][]interface{})

				for col, funcs := range aggFuncs {
					groupSeries := gb.getGroupSeries(col, indices)
					aggVals[col] = make([]interface{}, len(funcs))
					for j, fn := range funcs {
						aggVals[col][j] = fn(groupSeries)
					}
				}

				results[i] = groupResult{keyVals: keyVals, aggVals: aggVals}
			}
		}(start, end)
	}

	wg.Wait()

	// Collect results
	keyData := make(map[string][]interface{})
	for _, col := range gb.byKeys {
		keyData[col] = make([]interface{}, 0, numGroups)
	}

	aggData := make(map[string][]interface{})
	for col, funcs := range aggFuncs {
		for i := range funcs {
			aggCol := col + "_" + string(rune('0'+i))
			aggData[aggCol] = make([]interface{}, 0, numGroups)
		}
	}

	for _, r := range results {
		if r.keyVals == nil {
			continue
		}
		for i, col := range gb.byKeys {
			keyData[col] = append(keyData[col], r.keyVals[i])
		}
		for col, vals := range r.aggVals {
			for i, val := range vals {
				aggCol := col + "_" + string(rune('0'+i))
				aggData[aggCol] = append(aggData[aggCol], val)
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

// ParallelMap applies a mapping function to multiple Series in parallel
func ParallelMapSeries(series []*Series, fn func(*Series) *Series, opts ...ParallelOptions) []*Series {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	n := len(series)
	if n == 0 {
		return series
	}

	numWorkers := getNumWorkers(opt, n)
	if numWorkers > n {
		numWorkers = n
	}

	results := make([]*Series, n)
	var wg sync.WaitGroup

	seriesChan := make(chan int, n)
	for i := range series {
		seriesChan <- i
	}
	close(seriesChan)

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for i := range seriesChan {
				results[i] = fn(series[i])
			}
		}()
	}

	wg.Wait()
	return results
}

// ParallelReadCSV reads multiple CSV files in parallel and concatenates them
func ParallelReadCSV(paths []string, readFunc func(string) (*DataFrame, error), opts ...ParallelOptions) (*DataFrame, error) {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	n := len(paths)
	if n == 0 {
		return New(map[string][]interface{}{})
	}

	numWorkers := getNumWorkers(opt, n)
	if numWorkers > n {
		numWorkers = n
	}

	results := make([]*DataFrame, n)
	errors := make([]error, n)
	var wg sync.WaitGroup

	pathChan := make(chan int, n)
	for i := range paths {
		pathChan <- i
	}
	close(pathChan)

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for i := range pathChan {
				df, err := readFunc(paths[i])
				results[i] = df
				errors[i] = err
			}
		}()
	}

	wg.Wait()

	// Check for errors
	for i, err := range errors {
		if err != nil {
			return nil, err
		}
		if results[i] == nil {
			continue
		}
	}

	// Filter out nil results
	var validResults []*DataFrame
	for _, df := range results {
		if df != nil {
			validResults = append(validResults, df)
		}
	}

	if len(validResults) == 0 {
		return New(map[string][]interface{}{})
	}

	return Concat(validResults...), nil
}

// ChunkedApply applies a function to a Series in chunks for memory efficiency
func (s *Series) ChunkedApply(fn func([]interface{}) []interface{}, chunkSize int) *Series {
	if chunkSize <= 0 {
		chunkSize = 10000
	}

	n := s.Len()
	if n == 0 {
		return NewSeries([]interface{}{}, s.name)
	}

	result := make([]interface{}, 0, n)

	for start := 0; start < n; start += chunkSize {
		end := start + chunkSize
		if end > n {
			end = n
		}
		chunk := s.data[start:end]
		processed := fn(chunk)
		result = append(result, processed...)
	}

	return &Series{
		name:  s.name,
		data:  result,
		dtype: DTypeObject,
		index: s.index.Copy(),
	}
}

// ParallelChunkedApply applies a function to chunks in parallel
func (s *Series) ParallelChunkedApply(fn func([]interface{}) []interface{}, chunkSize int, opts ...ParallelOptions) *Series {
	opt := DefaultParallelOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	if chunkSize <= 0 {
		chunkSize = 10000
	}

	n := s.Len()
	if n == 0 {
		return NewSeries([]interface{}{}, s.name)
	}

	numChunks := (n + chunkSize - 1) / chunkSize
	numWorkers := getNumWorkers(opt, numChunks)
	if numWorkers > numChunks {
		numWorkers = numChunks
	}

	// Process chunks
	type chunkResult struct {
		index int
		data  []interface{}
	}

	resultChan := make(chan chunkResult, numChunks)
	chunkChan := make(chan int, numChunks)

	for i := 0; i < numChunks; i++ {
		chunkChan <- i
	}
	close(chunkChan)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for chunkIdx := range chunkChan {
				start := chunkIdx * chunkSize
				end := start + chunkSize
				if end > n {
					end = n
				}
				chunk := s.data[start:end]
				processed := fn(chunk)
				resultChan <- chunkResult{index: chunkIdx, data: processed}
			}
		}()
	}

	// Collect results in a separate goroutine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Assemble results in order
	chunkResults := make([][]interface{}, numChunks)
	for r := range resultChan {
		chunkResults[r.index] = r.data
	}

	// Concatenate
	result := make([]interface{}, 0, n)
	for _, chunk := range chunkResults {
		result = append(result, chunk...)
	}

	return &Series{
		name:  s.name,
		data:  result,
		dtype: DTypeObject,
		index: s.index.Copy(),
	}
}
