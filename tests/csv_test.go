package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/datago/dataframe"
	"github.com/datago/io"
)

func TestWriteCSVDataFrame(t *testing.T) {
	data := map[string][]interface{}{
		"name": {"alice", "bob"},
		"age":  {int64(30), int64(25)},
	}
	df, err := dataframe.New(data)
	if err != nil {
		t.Fatalf("DataFrame create error: %v", err)
	}

	outputDir := filepath.Join(".", "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("Create output dir error: %v", err)
	}
	path := filepath.Join(outputDir, "df.csv")
	if err := io.WriteCSV(path, df, io.CSVWriteOptions{IncludeIndex: false}); err != nil {
		t.Fatalf("WriteCSV error: %v", err)
	}

	readBack, err := io.ReadCSV(path, io.CSVOptions{HasHeader: true})
	if err != nil {
		t.Fatalf("ReadCSV error: %v", err)
	}

	if readBack.Shape()[0] != 2 || readBack.Shape()[1] != 2 {
		t.Fatalf("unexpected shape: %v", readBack.Shape())
	}

	nameSeries, ok := readBack.GetSeries("name")
	if !ok {
		t.Fatalf("missing column 'name'")
	}
	val, _ := nameSeries.Get(0)
	if val != "alice" {
		t.Fatalf("unexpected value: %v", val)
	}
}

func TestWriteCSVSeries(t *testing.T) {
	s := dataframe.NewSeriesFromStrings([]string{"x", "y", "z"}, "letter")

	outputDir := filepath.Join(".", "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("Create output dir error: %v", err)
	}
	path := filepath.Join(outputDir, "series.csv")
	if err := io.WriteSeriesCSV(path, s, io.CSVWriteOptions{IncludeIndex: false}); err != nil {
		t.Fatalf("WriteSeriesCSV error: %v", err)
	}

	readBack, err := io.ReadCSV(path, io.CSVOptions{HasHeader: true})
	if err != nil {
		t.Fatalf("ReadCSV error: %v", err)
	}

	if readBack.Shape()[0] != 3 || readBack.Shape()[1] != 1 {
		t.Fatalf("unexpected shape: %v", readBack.Shape())
	}

	series, ok := readBack.GetSeries("letter")
	if !ok {
		t.Fatalf("missing column 'letter'")
	}
	val, _ := series.Get(2)
	if val != "z" {
		t.Fatalf("unexpected value: %v", val)
	}
}
