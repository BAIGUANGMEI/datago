package tests

import (
	"testing"

	"github.com/BAIGUANGMEI/datago/dataframe"
)

func TestInnerJoin(t *testing.T) {
	// Create left DataFrame
	leftData := map[string][]interface{}{
		"id":   {1, 2, 3, 4},
		"name": {"Alice", "Bob", "Charlie", "David"},
	}
	left, err := dataframe.New(leftData)
	if err != nil {
		t.Fatalf("Failed to create left DataFrame: %v", err)
	}

	// Create right DataFrame
	rightData := map[string][]interface{}{
		"id":    {2, 3, 5},
		"score": {85.0, 90.0, 95.0},
	}
	right, err := dataframe.New(rightData)
	if err != nil {
		t.Fatalf("Failed to create right DataFrame: %v", err)
	}

	// Inner join
	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How: dataframe.InnerJoin,
		On:  []string{"id"},
	})
	if err != nil {
		t.Fatalf("Inner join failed: %v", err)
	}

	// Should have 2 matching rows (id 2 and 3)
	if result.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows, got %d", result.Shape()[0])
	}
}

func TestLeftJoin(t *testing.T) {
	leftData := map[string][]interface{}{
		"id":   {1, 2, 3},
		"name": {"Alice", "Bob", "Charlie"},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"id":    {2, 3, 4},
		"score": {85.0, 90.0, 95.0},
	}
	right, _ := dataframe.New(rightData)

	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How: dataframe.LeftJoin,
		On:  []string{"id"},
	})
	if err != nil {
		t.Fatalf("Left join failed: %v", err)
	}

	// Should have 3 rows (all from left)
	if result.Shape()[0] != 3 {
		t.Errorf("Expected 3 rows, got %d", result.Shape()[0])
	}
}

func TestRightJoin(t *testing.T) {
	leftData := map[string][]interface{}{
		"id":   {1, 2, 3},
		"name": {"Alice", "Bob", "Charlie"},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"id":    {2, 3, 4},
		"score": {85.0, 90.0, 95.0},
	}
	right, _ := dataframe.New(rightData)

	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How: dataframe.RightJoin,
		On:  []string{"id"},
	})
	if err != nil {
		t.Fatalf("Right join failed: %v", err)
	}

	// Should have 3 rows (all from right)
	if result.Shape()[0] != 3 {
		t.Errorf("Expected 3 rows, got %d", result.Shape()[0])
	}
}

func TestOuterJoin(t *testing.T) {
	leftData := map[string][]interface{}{
		"id":   {1, 2, 3},
		"name": {"Alice", "Bob", "Charlie"},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"id":    {2, 3, 4},
		"score": {85.0, 90.0, 95.0},
	}
	right, _ := dataframe.New(rightData)

	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How: dataframe.OuterJoin,
		On:  []string{"id"},
	})
	if err != nil {
		t.Fatalf("Outer join failed: %v", err)
	}

	// Should have 4 rows (1, 2, 3, 4)
	if result.Shape()[0] != 4 {
		t.Errorf("Expected 4 rows, got %d", result.Shape()[0])
	}
}

func TestMergeWithDifferentColumnNames(t *testing.T) {
	leftData := map[string][]interface{}{
		"left_id": {1, 2, 3},
		"name":    {"Alice", "Bob", "Charlie"},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"right_id": {2, 3, 4},
		"score":    {85.0, 90.0, 95.0},
	}
	right, _ := dataframe.New(rightData)

	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How:     dataframe.InnerJoin,
		LeftOn:  []string{"left_id"},
		RightOn: []string{"right_id"},
	})
	if err != nil {
		t.Fatalf("Merge with different column names failed: %v", err)
	}

	if result.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows, got %d", result.Shape()[0])
	}
}

func TestMergeWithIndicator(t *testing.T) {
	leftData := map[string][]interface{}{
		"id": {1, 2, 3},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"id": {2, 3, 4},
	}
	right, _ := dataframe.New(rightData)

	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How:       dataframe.OuterJoin,
		On:        []string{"id"},
		Indicator: true,
	})
	if err != nil {
		t.Fatalf("Merge with indicator failed: %v", err)
	}

	// Check that _merge column exists
	_, ok := result.GetSeries("_merge")
	if !ok {
		t.Error("Expected _merge column")
	}
}

func TestMergeWithSuffixes(t *testing.T) {
	leftData := map[string][]interface{}{
		"id":    {1, 2},
		"value": {10, 20},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"id":    {1, 2},
		"value": {100, 200},
	}
	right, _ := dataframe.New(rightData)

	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How:      dataframe.InnerJoin,
		On:       []string{"id"},
		Suffixes: [2]string{"_left", "_right"},
	})
	if err != nil {
		t.Fatalf("Merge with suffixes failed: %v", err)
	}

	// Check that suffixed columns exist
	_, ok1 := result.GetSeries("value_left")
	_, ok2 := result.GetSeries("value_right")
	if !ok1 || !ok2 {
		t.Error("Expected suffixed columns value_left and value_right")
	}
}

func TestDataFrameJoin(t *testing.T) {
	leftData := map[string][]interface{}{
		"id":   {1, 2, 3},
		"name": {"Alice", "Bob", "Charlie"},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"id":    {2, 3},
		"score": {85.0, 90.0},
	}
	right, _ := dataframe.New(rightData)

	result, err := left.Join(right, []string{"id"}, dataframe.InnerJoin)
	if err != nil {
		t.Fatalf("DataFrame.Join failed: %v", err)
	}

	if result.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows, got %d", result.Shape()[0])
	}
}

func TestMergeMultipleKeys(t *testing.T) {
	leftData := map[string][]interface{}{
		"year":    {2020, 2020, 2021, 2021},
		"quarter": {1, 2, 1, 2},
		"sales":   {100, 150, 200, 250},
	}
	left, _ := dataframe.New(leftData)

	rightData := map[string][]interface{}{
		"year":    {2020, 2021},
		"quarter": {1, 2},
		"target":  {120, 280},
	}
	right, _ := dataframe.New(rightData)

	result, err := dataframe.Merge(left, right, dataframe.MergeOptions{
		How: dataframe.InnerJoin,
		On:  []string{"year", "quarter"},
	})
	if err != nil {
		t.Fatalf("Merge with multiple keys failed: %v", err)
	}

	if result.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows, got %d", result.Shape()[0])
	}
}
