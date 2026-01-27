package dataframe

import (
	"fmt"
)

// Index represents the row/column index of a DataFrame or Series
type Index struct {
	labels []interface{} // Index labels
	name   string        // Index name
}

// NewIndex creates a new Index from labels
func NewIndex(labels []interface{}, name string) *Index {
	return &Index{
		labels: labels,
		name:   name,
	}
}

// NewRangeIndex creates a new integer range index
func NewRangeIndex(size int) *Index {
	labels := make([]interface{}, size)
	for i := 0; i < size; i++ {
		labels[i] = i
	}
	return &Index{
		labels: labels,
		name:   "",
	}
}

// Len returns the length of the index
func (idx *Index) Len() int {
	return len(idx.labels)
}

// Name returns the name of the index
func (idx *Index) Name() string {
	return idx.name
}

// SetName sets the name of the index
func (idx *Index) SetName(name string) {
	idx.name = name
}

// Labels returns all labels in the index
func (idx *Index) Labels() []interface{} {
	return idx.labels
}

// Get returns the label at the specified position
func (idx *Index) Get(pos int) (interface{}, error) {
	if pos < 0 || pos >= len(idx.labels) {
		return nil, fmt.Errorf("index %d out of range [0, %d)", pos, len(idx.labels))
	}
	return idx.labels[pos], nil
}

// GetLoc returns the position of the specified label
func (idx *Index) GetLoc(label interface{}) (int, error) {
	for i, l := range idx.labels {
		if l == label {
			return i, nil
		}
	}
	return -1, fmt.Errorf("label %v not found in index", label)
}

// Contains checks if the index contains the specified label
func (idx *Index) Contains(label interface{}) bool {
	_, err := idx.GetLoc(label)
	return err == nil
}

// Slice returns a new index with elements from start to end
func (idx *Index) Slice(start, end int) *Index {
	if start < 0 {
		start = 0
	}
	if end > len(idx.labels) {
		end = len(idx.labels)
	}
	newLabels := make([]interface{}, end-start)
	copy(newLabels, idx.labels[start:end])
	return &Index{
		labels: newLabels,
		name:   idx.name,
	}
}

// Append adds a new label to the index
func (idx *Index) Append(label interface{}) *Index {
	newLabels := make([]interface{}, len(idx.labels)+1)
	copy(newLabels, idx.labels)
	newLabels[len(idx.labels)] = label
	return &Index{
		labels: newLabels,
		name:   idx.name,
	}
}

// Copy creates a copy of the index
func (idx *Index) Copy() *Index {
	newLabels := make([]interface{}, len(idx.labels))
	copy(newLabels, idx.labels)
	return &Index{
		labels: newLabels,
		name:   idx.name,
	}
}

// Reset resets the index to default integer range
func (idx *Index) Reset() *Index {
	return NewRangeIndex(len(idx.labels))
}

// ToStringSlice converts index labels to string slice
func (idx *Index) ToStringSlice() []string {
	result := make([]string, len(idx.labels))
	for i, label := range idx.labels {
		result[i] = fmt.Sprintf("%v", label)
	}
	return result
}

// Equals checks if two indexes are equal
func (idx *Index) Equals(other *Index) bool {
	if idx.Len() != other.Len() {
		return false
	}
	for i, label := range idx.labels {
		if label != other.labels[i] {
			return false
		}
	}
	return true
}

// Union returns the union of two indexes
func (idx *Index) Union(other *Index) *Index {
	seen := make(map[interface{}]bool)
	var labels []interface{}

	for _, label := range idx.labels {
		if !seen[label] {
			seen[label] = true
			labels = append(labels, label)
		}
	}
	for _, label := range other.labels {
		if !seen[label] {
			seen[label] = true
			labels = append(labels, label)
		}
	}

	return &Index{
		labels: labels,
		name:   idx.name,
	}
}

// Intersection returns the intersection of two indexes
func (idx *Index) Intersection(other *Index) *Index {
	otherSet := make(map[interface{}]bool)
	for _, label := range other.labels {
		otherSet[label] = true
	}

	var labels []interface{}
	seen := make(map[interface{}]bool)
	for _, label := range idx.labels {
		if otherSet[label] && !seen[label] {
			seen[label] = true
			labels = append(labels, label)
		}
	}

	return &Index{
		labels: labels,
		name:   idx.name,
	}
}

// Difference returns the difference of two indexes (elements in idx but not in other)
func (idx *Index) Difference(other *Index) *Index {
	otherSet := make(map[interface{}]bool)
	for _, label := range other.labels {
		otherSet[label] = true
	}

	var labels []interface{}
	for _, label := range idx.labels {
		if !otherSet[label] {
			labels = append(labels, label)
		}
	}

	return &Index{
		labels: labels,
		name:   idx.name,
	}
}
