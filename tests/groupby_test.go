package tests

import (
	"testing"

	"github.com/datago/dataframe"
)

func TestGroupBy(t *testing.T) {
	// Create test data
	data := map[string][]interface{}{
		"category": {"A", "A", "B", "B", "A"},
		"value":    {10.0, 20.0, 30.0, 40.0, 50.0},
		"count":    {1, 2, 3, 4, 5},
	}
	df, err := dataframe.New(data)
	if err != nil {
		t.Fatalf("Failed to create DataFrame: %v", err)
	}

	// Test GroupBy creation
	gb, err := df.GroupBy("category")
	if err != nil {
		t.Fatalf("Failed to create GroupBy: %v", err)
	}

	// Test NGroups
	if gb.NGroups() != 2 {
		t.Errorf("Expected 2 groups, got %d", gb.NGroups())
	}

	// Test Size
	sizeDF := gb.Size()
	if sizeDF.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows in size DataFrame, got %d", sizeDF.Shape()[0])
	}

	// Test Sum
	sumDF := gb.Sum("value")
	if sumDF == nil {
		t.Fatal("Sum returned nil")
	}
	if sumDF.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows in sum DataFrame, got %d", sumDF.Shape()[0])
	}

	// Test Mean
	meanDF := gb.Mean("value")
	if meanDF == nil {
		t.Fatal("Mean returned nil")
	}

	// Test Count
	countDF := gb.Count("value")
	if countDF == nil {
		t.Fatal("Count returned nil")
	}
}

func TestGroupByMultipleColumns(t *testing.T) {
	data := map[string][]interface{}{
		"region":   {"East", "East", "West", "West", "East"},
		"category": {"A", "B", "A", "B", "A"},
		"sales":    {100.0, 200.0, 150.0, 250.0, 120.0},
	}
	df, err := dataframe.New(data)
	if err != nil {
		t.Fatalf("Failed to create DataFrame: %v", err)
	}

	gb, err := df.GroupBy("region", "category")
	if err != nil {
		t.Fatalf("Failed to create GroupBy: %v", err)
	}

	if gb.NGroups() != 4 {
		t.Errorf("Expected 4 groups, got %d", gb.NGroups())
	}

	sumDF := gb.Sum("sales")
	if sumDF.Shape()[0] != 4 {
		t.Errorf("Expected 4 rows, got %d", sumDF.Shape()[0])
	}
}

func TestGroupByAgg(t *testing.T) {
	data := map[string][]interface{}{
		"group": {"A", "A", "B", "B"},
		"value": {10.0, 20.0, 30.0, 40.0},
	}
	df, err := dataframe.New(data)
	if err != nil {
		t.Fatalf("Failed to create DataFrame: %v", err)
	}

	gb, err := df.GroupBy("group")
	if err != nil {
		t.Fatalf("Failed to create GroupBy: %v", err)
	}

	aggFuncs := map[string][]dataframe.AggFunc{
		"value": {dataframe.AggSum, dataframe.AggMean},
	}

	result, err := gb.Agg(aggFuncs)
	if err != nil {
		t.Fatalf("Agg failed: %v", err)
	}

	if result.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows, got %d", result.Shape()[0])
	}
}

func TestGroupByFilter(t *testing.T) {
	data := map[string][]interface{}{
		"group": {"A", "A", "B", "B", "B"},
		"value": {10.0, 20.0, 30.0, 40.0, 50.0},
	}
	df, err := dataframe.New(data)
	if err != nil {
		t.Fatalf("Failed to create DataFrame: %v", err)
	}

	gb, err := df.GroupBy("group")
	if err != nil {
		t.Fatalf("Failed to create GroupBy: %v", err)
	}

	// Filter groups with more than 2 elements
	filtered := gb.Filter(func(groupDF *dataframe.DataFrame) bool {
		return groupDF.Shape()[0] > 2
	})

	if filtered.Shape()[0] != 3 {
		t.Errorf("Expected 3 rows (group B), got %d", filtered.Shape()[0])
	}
}

func TestGroupByApply(t *testing.T) {
	data := map[string][]interface{}{
		"group": {"A", "A", "B", "B"},
		"value": {10.0, 20.0, 30.0, 40.0},
	}
	df, err := dataframe.New(data)
	if err != nil {
		t.Fatalf("Failed to create DataFrame: %v", err)
	}

	gb, err := df.GroupBy("group")
	if err != nil {
		t.Fatalf("Failed to create GroupBy: %v", err)
	}

	// Apply function that returns first row of each group
	result := gb.Apply(func(groupDF *dataframe.DataFrame) *dataframe.DataFrame {
		return groupDF.Head(1)
	})

	if result.Shape()[0] != 2 {
		t.Errorf("Expected 2 rows, got %d", result.Shape()[0])
	}
}

func TestConcat(t *testing.T) {
	data1 := map[string][]interface{}{
		"a": {1, 2},
		"b": {"x", "y"},
	}
	df1, _ := dataframe.New(data1)

	data2 := map[string][]interface{}{
		"a": {3, 4},
		"b": {"z", "w"},
	}
	df2, _ := dataframe.New(data2)

	result := dataframe.Concat(df1, df2)

	if result.Shape()[0] != 4 {
		t.Errorf("Expected 4 rows, got %d", result.Shape()[0])
	}
	if result.Shape()[1] != 2 {
		t.Errorf("Expected 2 columns, got %d", result.Shape()[1])
	}
}
