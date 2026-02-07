package tests

import (
	"testing"

	"github.com/BAIGUANGMEI/datago/dataframe"
)

func TestDataFrameNewAndSelect(t *testing.T) {
	data := map[string][]interface{}{
		"name": {"Alice", "Bob", "Charlie"},
		"age":  {25, 30, 35},
	}
	df, err := dataframe.New(data)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if df.Shape()[0] != 3 || df.Shape()[1] != 2 {
		t.Fatalf("Shape() = %v, want [3 2]", df.Shape())
	}

	sub := df.Select("age")
	if sub.Shape()[1] != 1 {
		t.Fatalf("Select() cols = %d, want 1", sub.Shape()[1])
	}

	t.Log(df)
}

func TestDataFrameFilterSort(t *testing.T) {
	data := map[string][]interface{}{
		"name": {"Alice", "Bob", "Charlie"},
		"age":  {25, 30, 20},
	}
	df, _ := dataframe.New(data)

	filtered := df.Filter(func(r dataframe.Row) bool {
		return r.Get("age").(int) >= 25
	})
	if filtered.Shape()[0] != 2 {
		t.Fatalf("Filter() rows = %d, want 2", filtered.Shape()[0])
	}

	sorted := df.SortBy("age", dataframe.Ascending)
	v, _ := sorted.At(0, "age")
	if v != 20 {
		t.Fatalf("SortBy() first age = %v, want 20", v)
	}
}

func TestDataFrameDescribe(t *testing.T) {
	data := map[string][]interface{}{
		"a": {1, 2, 3},
		"b": {10, 20, 30},
	}
	df, _ := dataframe.New(data)
	desc := df.Describe()
	if desc.Shape()[0] != 2 || desc.Shape()[1] != 5 {
		t.Fatalf("Describe() shape = %v, want [2 5]", desc.Shape())
	}
}
