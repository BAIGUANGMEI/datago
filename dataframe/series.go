package dataframe

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Series represents a one-dimensional labeled array
type Series struct {
	name  string        // Series name
	data  []interface{} // Data values
	dtype DType         // Data type
	index *Index        // Row index
}

// NewSeries creates a new Series from data
func NewSeries(data []interface{}, name string) *Series {
	dtype := InferDTypeFromSlice(data)
	return &Series{
		name:  name,
		data:  data,
		dtype: dtype,
		index: NewRangeIndex(len(data)),
	}
}

// NewSeriesWithIndex creates a new Series with custom index
func NewSeriesWithIndex(data []interface{}, name string, index *Index) *Series {
	if index == nil {
		index = NewRangeIndex(len(data))
	}
	dtype := InferDTypeFromSlice(data)
	return &Series{
		name:  name,
		data:  data,
		dtype: dtype,
		index: index,
	}
}

// NewSeriesFromInts creates a Series from int slice
func NewSeriesFromInts(data []int, name string) *Series {
	values := make([]interface{}, len(data))
	for i, v := range data {
		values[i] = int64(v)
	}
	return &Series{
		name:  name,
		data:  values,
		dtype: DTypeInt64,
		index: NewRangeIndex(len(data)),
	}
}

// NewSeriesFromInt64s creates a Series from int64 slice
func NewSeriesFromInt64s(data []int64, name string) *Series {
	values := make([]interface{}, len(data))
	for i, v := range data {
		values[i] = v
	}
	return &Series{
		name:  name,
		data:  values,
		dtype: DTypeInt64,
		index: NewRangeIndex(len(data)),
	}
}

// NewSeriesFromFloat64s creates a Series from float64 slice
func NewSeriesFromFloat64s(data []float64, name string) *Series {
	values := make([]interface{}, len(data))
	for i, v := range data {
		values[i] = v
	}
	return &Series{
		name:  name,
		data:  values,
		dtype: DTypeFloat64,
		index: NewRangeIndex(len(data)),
	}
}

// NewSeriesFromStrings creates a Series from string slice
func NewSeriesFromStrings(data []string, name string) *Series {
	values := make([]interface{}, len(data))
	for i, v := range data {
		values[i] = v
	}
	return &Series{
		name:  name,
		data:  values,
		dtype: DTypeString,
		index: NewRangeIndex(len(data)),
	}
}

// NewSeriesFromBools creates a Series from bool slice
func NewSeriesFromBools(data []bool, name string) *Series {
	values := make([]interface{}, len(data))
	for i, v := range data {
		values[i] = v
	}
	return &Series{
		name:  name,
		data:  values,
		dtype: DTypeBool,
		index: NewRangeIndex(len(data)),
	}
}

// Name returns the name of the Series
func (s *Series) Name() string {
	return s.name
}

// SetName sets the name of the Series
func (s *Series) SetName(name string) *Series {
	s.name = name
	return s
}

// DType returns the data type of the Series
func (s *Series) DType() DType {
	return s.dtype
}

// Index returns the index of the Series
func (s *Series) Index() *Index {
	return s.index
}

// SetIndex sets the index of the Series
func (s *Series) SetIndex(index *Index) *Series {
	s.index = index
	return s
}

// Len returns the length of the Series
func (s *Series) Len() int {
	return len(s.data)
}

// Values returns all values in the Series
func (s *Series) Values() []interface{} {
	return s.data
}

// Get returns the value at the specified position
func (s *Series) Get(pos int) (interface{}, error) {
	if pos < 0 || pos >= len(s.data) {
		return nil, fmt.Errorf("index %d out of range [0, %d)", pos, len(s.data))
	}
	return s.data[pos], nil
}

// At returns the value at the specified label
func (s *Series) At(label interface{}) (interface{}, error) {
	pos, err := s.index.GetLoc(label)
	if err != nil {
		return nil, err
	}
	return s.data[pos], nil
}

// Set sets the value at the specified position
func (s *Series) Set(pos int, value interface{}) error {
	if pos < 0 || pos >= len(s.data) {
		return fmt.Errorf("index %d out of range [0, %d)", pos, len(s.data))
	}
	s.data[pos] = value
	return nil
}

// Copy creates a copy of the Series
func (s *Series) Copy() *Series {
	newData := make([]interface{}, len(s.data))
	copy(newData, s.data)
	return &Series{
		name:  s.name,
		data:  newData,
		dtype: s.dtype,
		index: s.index.Copy(),
	}
}

// Head returns the first n elements
func (s *Series) Head(n int) *Series {
	if n > len(s.data) {
		n = len(s.data)
	}
	return &Series{
		name:  s.name,
		data:  s.data[:n],
		dtype: s.dtype,
		index: s.index.Slice(0, n),
	}
}

// Tail returns the last n elements
func (s *Series) Tail(n int) *Series {
	if n > len(s.data) {
		n = len(s.data)
	}
	start := len(s.data) - n
	return &Series{
		name:  s.name,
		data:  s.data[start:],
		dtype: s.dtype,
		index: s.index.Slice(start, len(s.data)),
	}
}

// Slice returns a slice of the Series
func (s *Series) Slice(start, end int) *Series {
	if start < 0 {
		start = 0
	}
	if end > len(s.data) {
		end = len(s.data)
	}
	return &Series{
		name:  s.name,
		data:  s.data[start:end],
		dtype: s.dtype,
		index: s.index.Slice(start, end),
	}
}

// ============ Statistical Methods ============

// Sum returns the sum of all numeric values
func (s *Series) Sum() float64 {
	var sum float64
	for _, v := range s.data {
		if v == nil || IsNA(v) {
			continue
		}
		f, err := toFloat64(v)
		if err == nil {
			sum += f
		}
	}
	return sum
}

// Mean returns the mean of all numeric values
func (s *Series) Mean() float64 {
	count := 0
	var sum float64
	for _, v := range s.data {
		if v == nil || IsNA(v) {
			continue
		}
		f, err := toFloat64(v)
		if err == nil {
			sum += f
			count++
		}
	}
	if count == 0 {
		return math.NaN()
	}
	return sum / float64(count)
}

// Median returns the median of all numeric values
func (s *Series) Median() float64 {
	var values []float64
	for _, v := range s.data {
		if v == nil || IsNA(v) {
			continue
		}
		f, err := toFloat64(v)
		if err == nil {
			values = append(values, f)
		}
	}
	if len(values) == 0 {
		return math.NaN()
	}
	sort.Float64s(values)
	n := len(values)
	if n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2
	}
	return values[n/2]
}

// Std returns the standard deviation
func (s *Series) Std() float64 {
	return math.Sqrt(s.Var())
}

// Var returns the variance
func (s *Series) Var() float64 {
	mean := s.Mean()
	if math.IsNaN(mean) {
		return math.NaN()
	}
	var sumSq float64
	count := 0
	for _, v := range s.data {
		if v == nil || IsNA(v) {
			continue
		}
		f, err := toFloat64(v)
		if err == nil {
			diff := f - mean
			sumSq += diff * diff
			count++
		}
	}
	if count <= 1 {
		return math.NaN()
	}
	return sumSq / float64(count-1) // Sample variance (n-1)
}

// Min returns the minimum value
func (s *Series) Min() interface{} {
	var minVal float64 = math.MaxFloat64
	found := false
	for _, v := range s.data {
		if v == nil || IsNA(v) {
			continue
		}
		f, err := toFloat64(v)
		if err == nil {
			if f < minVal {
				minVal = f
				found = true
			}
		}
	}
	if !found {
		return nil
	}
	return minVal
}

// Max returns the maximum value
func (s *Series) Max() interface{} {
	var maxVal float64 = -math.MaxFloat64
	found := false
	for _, v := range s.data {
		if v == nil || IsNA(v) {
			continue
		}
		f, err := toFloat64(v)
		if err == nil {
			if f > maxVal {
				maxVal = f
				found = true
			}
		}
	}
	if !found {
		return nil
	}
	return maxVal
}

// Count returns the number of non-NA values
func (s *Series) Count() int {
	count := 0
	for _, v := range s.data {
		if v != nil && !IsNA(v) {
			count++
		}
	}
	return count
}

// Unique returns a Series with unique values
func (s *Series) Unique() *Series {
	seen := make(map[interface{}]bool)
	var unique []interface{}
	for _, v := range s.data {
		key := fmt.Sprintf("%v", v)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, v)
		}
	}
	return NewSeries(unique, s.name)
}

// NUnique returns the number of unique values
func (s *Series) NUnique() int {
	return s.Unique().Len()
}

// ValueCounts returns a Series with counts of unique values
func (s *Series) ValueCounts() *Series {
	counts := make(map[string]int)
	for _, v := range s.data {
		key := fmt.Sprintf("%v", v)
		counts[key]++
	}

	var labels []interface{}
	var values []interface{}
	for k, v := range counts {
		labels = append(labels, k)
		values = append(values, v)
	}

	result := NewSeries(values, "count")
	result.index = NewIndex(labels, s.name)
	return result
}

// ============ Data Manipulation Methods ============

// Apply applies a function to each element
func (s *Series) Apply(fn func(interface{}) interface{}) *Series {
	newData := make([]interface{}, len(s.data))
	for i, v := range s.data {
		newData[i] = fn(v)
	}
	return &Series{
		name:  s.name,
		data:  newData,
		dtype: DTypeObject,
		index: s.index.Copy(),
	}
}

// Map applies a mapping to each element
func (s *Series) Map(mapping map[interface{}]interface{}) *Series {
	newData := make([]interface{}, len(s.data))
	for i, v := range s.data {
		if mapped, ok := mapping[v]; ok {
			newData[i] = mapped
		} else {
			newData[i] = v
		}
	}
	return &Series{
		name:  s.name,
		data:  newData,
		dtype: DTypeObject,
		index: s.index.Copy(),
	}
}

// FillNA fills NA values with the specified value
func (s *Series) FillNA(value interface{}) *Series {
	newData := make([]interface{}, len(s.data))
	for i, v := range s.data {
		if v == nil || IsNA(v) {
			newData[i] = value
		} else {
			newData[i] = v
		}
	}
	return &Series{
		name:  s.name,
		data:  newData,
		dtype: s.dtype,
		index: s.index.Copy(),
	}
}

// DropNA removes NA values
func (s *Series) DropNA() *Series {
	var newData []interface{}
	var newLabels []interface{}
	for i, v := range s.data {
		if v != nil && !IsNA(v) {
			newData = append(newData, v)
			label, _ := s.index.Get(i)
			newLabels = append(newLabels, label)
		}
	}
	return &Series{
		name:  s.name,
		data:  newData,
		dtype: s.dtype,
		index: NewIndex(newLabels, s.index.Name()),
	}
}

// IsNA returns a boolean Series indicating NA values
func (s *Series) IsNA() *Series {
	newData := make([]interface{}, len(s.data))
	for i, v := range s.data {
		newData[i] = v == nil || IsNA(v)
	}
	return &Series{
		name:  s.name + "_isna",
		data:  newData,
		dtype: DTypeBool,
		index: s.index.Copy(),
	}
}

// NotNA returns a boolean Series indicating non-NA values
func (s *Series) NotNA() *Series {
	newData := make([]interface{}, len(s.data))
	for i, v := range s.data {
		newData[i] = v != nil && !IsNA(v)
	}
	return &Series{
		name:  s.name + "_notna",
		data:  newData,
		dtype: DTypeBool,
		index: s.index.Copy(),
	}
}

// AsType converts the Series to the specified data type
func (s *Series) AsType(dtype DType) (*Series, error) {
	newData := make([]interface{}, len(s.data))
	for i, v := range s.data {
		converted, err := ConvertToType(v, dtype)
		if err != nil {
			return nil, fmt.Errorf("error converting element %d: %w", i, err)
		}
		newData[i] = converted
	}
	return &Series{
		name:  s.name,
		data:  newData,
		dtype: dtype,
		index: s.index.Copy(),
	}, nil
}

// SortValues sorts the Series by values
func (s *Series) SortValues(ascending bool) *Series {
	type indexedValue struct {
		index int
		value interface{}
	}

	indexed := make([]indexedValue, len(s.data))
	for i, v := range s.data {
		indexed[i] = indexedValue{i, v}
	}

	sort.Slice(indexed, func(i, j int) bool {
		vi, vj := indexed[i].value, indexed[j].value
		// Handle nil values
		if vi == nil && vj == nil {
			return false
		}
		if vi == nil {
			return !ascending
		}
		if vj == nil {
			return ascending
		}
		// Compare values
		fi, erri := toFloat64(vi)
		fj, errj := toFloat64(vj)
		if erri == nil && errj == nil {
			if ascending {
				return fi < fj
			}
			return fi > fj
		}
		// Fall back to string comparison
		if ascending {
			return fmt.Sprintf("%v", vi) < fmt.Sprintf("%v", vj)
		}
		return fmt.Sprintf("%v", vi) > fmt.Sprintf("%v", vj)
	})

	newData := make([]interface{}, len(s.data))
	newLabels := make([]interface{}, len(s.data))
	for i, iv := range indexed {
		newData[i] = iv.value
		label, _ := s.index.Get(iv.index)
		newLabels[i] = label
	}

	return &Series{
		name:  s.name,
		data:  newData,
		dtype: s.dtype,
		index: NewIndex(newLabels, s.index.Name()),
	}
}

// ============ String Representation ============

// String returns the string representation of the Series
func (s *Series) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Series: %s (dtype: %s, length: %d)\n", s.name, s.dtype, len(s.data)))

	maxShow := 10
	if len(s.data) <= maxShow*2 {
		for i, v := range s.data {
			label, _ := s.index.Get(i)
			sb.WriteString(fmt.Sprintf("%v    %v\n", label, v))
		}
	} else {
		// Show first and last elements
		for i := 0; i < maxShow; i++ {
			label, _ := s.index.Get(i)
			sb.WriteString(fmt.Sprintf("%v    %v\n", label, s.data[i]))
		}
		sb.WriteString("...\n")
		for i := len(s.data) - maxShow; i < len(s.data); i++ {
			label, _ := s.index.Get(i)
			sb.WriteString(fmt.Sprintf("%v    %v\n", label, s.data[i]))
		}
	}
	return sb.String()
}

// ============ Arithmetic Operations ============

// Add adds a value or Series to this Series
func (s *Series) Add(other interface{}) *Series {
	return s.arithmeticOp(other, func(a, b float64) float64 { return a + b })
}

// Sub subtracts a value or Series from this Series
func (s *Series) Sub(other interface{}) *Series {
	return s.arithmeticOp(other, func(a, b float64) float64 { return a - b })
}

// Mul multiplies this Series by a value or Series
func (s *Series) Mul(other interface{}) *Series {
	return s.arithmeticOp(other, func(a, b float64) float64 { return a * b })
}

// Div divides this Series by a value or Series
func (s *Series) Div(other interface{}) *Series {
	return s.arithmeticOp(other, func(a, b float64) float64 {
		if b == 0 {
			return math.NaN()
		}
		return a / b
	})
}

func (s *Series) arithmeticOp(other interface{}, op func(float64, float64) float64) *Series {
	newData := make([]interface{}, len(s.data))

	switch v := other.(type) {
	case *Series:
		for i := 0; i < len(s.data); i++ {
			if i >= len(v.data) {
				newData[i] = nil
				continue
			}
			a, erra := toFloat64(s.data[i])
			b, errb := toFloat64(v.data[i])
			if erra != nil || errb != nil {
				newData[i] = nil
			} else {
				newData[i] = op(a, b)
			}
		}
	default:
		bVal, err := toFloat64(other)
		if err != nil {
			for i := range newData {
				newData[i] = nil
			}
		} else {
			for i, val := range s.data {
				a, err := toFloat64(val)
				if err != nil {
					newData[i] = nil
				} else {
					newData[i] = op(a, bVal)
				}
			}
		}
	}

	return &Series{
		name:  s.name,
		data:  newData,
		dtype: DTypeFloat64,
		index: s.index.Copy(),
	}
}
