package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BAIGUANGMEI/datago/dataframe"
	"github.com/BAIGUANGMEI/datago/io"
)

func TestReadExcelBasic(t *testing.T) {
	path := "testdata.xlsx"
	df, err := io.ReadExcel(path, io.ExcelOptions{HasHeader: true})
	if err != nil {
		t.Fatalf("ReadExcel error: %v", err)
	}

	t.Log(df)
}

func TestWriteExcelDataFrame(t *testing.T) {
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
	path := filepath.Join(outputDir, "df.xlsx")
	if err := io.WriteExcel(path, df, io.ExcelWriteOptions{IncludeIndex: false}); err != nil {
		t.Fatalf("WriteExcel error: %v", err)
	}

	readBack, err := io.ReadExcel(path, io.ExcelOptions{HasHeader: true})
	if err != nil {
		t.Fatalf("ReadExcel error: %v", err)
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

func TestWriteExcelSeries(t *testing.T) {
	s := dataframe.NewSeriesFromStrings([]string{"x", "y", "z"}, "letter")

	outputDir := filepath.Join(".", "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("Create output dir error: %v", err)
	}
	path := filepath.Join(outputDir, "series.xlsx")
	if err := io.WriteSeriesExcel(path, s, io.ExcelWriteOptions{IncludeIndex: false}); err != nil {
		t.Fatalf("WriteSeriesExcel error: %v", err)
	}

	readBack, err := io.ReadExcel(path, io.ExcelOptions{HasHeader: true})
	if err != nil {
		t.Fatalf("ReadExcel error: %v", err)
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
