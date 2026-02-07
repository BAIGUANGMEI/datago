package tests

import (
	"testing"

	"github.com/BAIGUANGMEI/datago/dataframe"
)

func TestParallelApply(t *testing.T) {
	// Create test data
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = float64(i)
	}
	s := dataframe.NewSeries(data, "values")

	// Apply function that doubles values
	result := s.ParallelApply(func(v interface{}) interface{} {
		if f, ok := v.(float64); ok {
			return f * 2
		}
		return v
	})

	if result.Len() != s.Len() {
		t.Errorf("Expected length %d, got %d", s.Len(), result.Len())
	}

	// Check first and last values
	first, _ := result.Get(0)
	if first.(float64) != 0 {
		t.Errorf("Expected first value 0, got %v", first)
	}

	last, _ := result.Get(9999)
	if last.(float64) != 19998 {
		t.Errorf("Expected last value 19998, got %v", last)
	}
}

func TestParallelFilter(t *testing.T) {
	// Create test DataFrame
	data := map[string][]interface{}{
		"value": make([]interface{}, 1000),
	}
	for i := 0; i < 1000; i++ {
		data["value"][i] = i
	}
	df, _ := dataframe.New(data)

	// Filter values > 500
	result := df.ParallelFilter(func(row dataframe.Row) bool {
		val := row.Get("value")
		if v, ok := val.(int); ok {
			return v > 500
		}
		return false
	})

	if result.Shape()[0] != 499 {
		t.Errorf("Expected 499 rows, got %d", result.Shape()[0])
	}
}

func TestParallelTransform(t *testing.T) {
	// Create test DataFrame
	data := map[string][]interface{}{
		"a": {1.0, 2.0, 3.0, 4.0, 5.0},
		"b": {10.0, 20.0, 30.0, 40.0, 50.0},
	}
	df, _ := dataframe.New(data)

	// Transform: multiply each value by 2
	result := df.ParallelTransform(func(s *dataframe.Series) *dataframe.Series {
		return s.Mul(2.0)
	})

	if result.Shape() != df.Shape() {
		t.Errorf("Shape mismatch")
	}
}

func TestParallelSum(t *testing.T) {
	data := map[string][]interface{}{
		"a": {1.0, 2.0, 3.0, 4.0, 5.0},
		"b": {10.0, 20.0, 30.0, 40.0, 50.0},
	}
	df, _ := dataframe.New(data)

	sums := df.ParallelSum()

	if sums["a"] != 15.0 {
		t.Errorf("Expected sum of a = 15, got %v", sums["a"])
	}
	if sums["b"] != 150.0 {
		t.Errorf("Expected sum of b = 150, got %v", sums["b"])
	}
}

func TestParallelMean(t *testing.T) {
	data := map[string][]interface{}{
		"a": {1.0, 2.0, 3.0, 4.0, 5.0},
		"b": {10.0, 20.0, 30.0, 40.0, 50.0},
	}
	df, _ := dataframe.New(data)

	means := df.ParallelMean()

	if means["a"] != 3.0 {
		t.Errorf("Expected mean of a = 3, got %v", means["a"])
	}
	if means["b"] != 30.0 {
		t.Errorf("Expected mean of b = 30, got %v", means["b"])
	}
}

func TestParallelGroupByAgg(t *testing.T) {
	data := map[string][]interface{}{
		"group": {"A", "A", "B", "B", "A", "B"},
		"value": {10.0, 20.0, 30.0, 40.0, 50.0, 60.0},
	}
	df, _ := dataframe.New(data)

	gb, _ := df.GroupBy("group")

	aggFuncs := map[string][]dataframe.AggFunc{
		"value": {dataframe.AggSum, dataframe.AggMean},
	}

	result, err := gb.ParallelAgg(aggFuncs)
	if err != nil {
		t.Fatalf("ParallelAgg failed: %v", err)
	}

	if result.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows, got %d", result.Shape()[0])
	}
}

func TestChunkedApply(t *testing.T) {
	// Create large data
	data := make([]interface{}, 50000)
	for i := range data {
		data[i] = float64(i)
	}
	s := dataframe.NewSeries(data, "values")

	// Apply function in chunks
	result := s.ChunkedApply(func(chunk []interface{}) []interface{} {
		result := make([]interface{}, len(chunk))
		for i, v := range chunk {
			if f, ok := v.(float64); ok {
				result[i] = f * 2
			} else {
				result[i] = v
			}
		}
		return result
	}, 10000)

	if result.Len() != s.Len() {
		t.Errorf("Expected length %d, got %d", s.Len(), result.Len())
	}
}

func TestParallelChunkedApply(t *testing.T) {
	// Create large data
	data := make([]interface{}, 100000)
	for i := range data {
		data[i] = float64(i)
	}
	s := dataframe.NewSeries(data, "values")

	// Apply function in parallel chunks
	result := s.ParallelChunkedApply(func(chunk []interface{}) []interface{} {
		result := make([]interface{}, len(chunk))
		for i, v := range chunk {
			if f, ok := v.(float64); ok {
				result[i] = f + 1
			} else {
				result[i] = v
			}
		}
		return result
	}, 10000)

	if result.Len() != s.Len() {
		t.Errorf("Expected length %d, got %d", s.Len(), result.Len())
	}

	// Verify values
	first, _ := result.Get(0)
	if first.(float64) != 1.0 {
		t.Errorf("Expected first value 1, got %v", first)
	}

	last, _ := result.Get(99999)
	if last.(float64) != 100000.0 {
		t.Errorf("Expected last value 100000, got %v", last)
	}
}

func TestParallelMapSeries(t *testing.T) {
	series := make([]*dataframe.Series, 5)
	for i := 0; i < 5; i++ {
		data := make([]interface{}, 100)
		for j := range data {
			data[j] = float64(j * (i + 1))
		}
		series[i] = dataframe.NewSeries(data, "series_"+string(rune('0'+i)))
	}

	// Map function that doubles values
	results := dataframe.ParallelMapSeries(series, func(s *dataframe.Series) *dataframe.Series {
		return s.Mul(2.0)
	})

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	for i, s := range results {
		if s.Len() != 100 {
			t.Errorf("Series %d: expected length 100, got %d", i, s.Len())
		}
	}
}

func TestParallelWithCustomOptions(t *testing.T) {
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = float64(i)
	}
	s := dataframe.NewSeries(data, "values")

	opts := dataframe.ParallelOptions{
		NumWorkers: 4,
		ChunkSize:  500,
	}

	result := s.ParallelApply(func(v interface{}) interface{} {
		if f, ok := v.(float64); ok {
			return f * 3
		}
		return v
	}, opts)

	if result.Len() != s.Len() {
		t.Errorf("Expected length %d, got %d", s.Len(), result.Len())
	}
}

// Benchmark tests
func BenchmarkSeriesApply(b *testing.B) {
	data := make([]interface{}, 100000)
	for i := range data {
		data[i] = float64(i)
	}
	s := dataframe.NewSeries(data, "values")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Apply(func(v interface{}) interface{} {
			if f, ok := v.(float64); ok {
				return f * 2
			}
			return v
		})
	}
}

func BenchmarkSeriesParallelApply(b *testing.B) {
	data := make([]interface{}, 100000)
	for i := range data {
		data[i] = float64(i)
	}
	s := dataframe.NewSeries(data, "values")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ParallelApply(func(v interface{}) interface{} {
			if f, ok := v.(float64); ok {
				return f * 2
			}
			return v
		})
	}
}

func BenchmarkDataFrameFilter(b *testing.B) {
	data := map[string][]interface{}{
		"value": make([]interface{}, 100000),
	}
	for i := 0; i < 100000; i++ {
		data["value"][i] = i
	}
	df, _ := dataframe.New(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		df.Filter(func(row dataframe.Row) bool {
			val := row.Get("value")
			if v, ok := val.(int); ok {
				return v > 50000
			}
			return false
		})
	}
}

func BenchmarkDataFrameParallelFilter(b *testing.B) {
	data := map[string][]interface{}{
		"value": make([]interface{}, 100000),
	}
	for i := 0; i < 100000; i++ {
		data["value"][i] = i
	}
	df, _ := dataframe.New(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		df.ParallelFilter(func(row dataframe.Row) bool {
			val := row.Get("value")
			if v, ok := val.(int); ok {
				return v > 50000
			}
			return false
		})
	}
}
