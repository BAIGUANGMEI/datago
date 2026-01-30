package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/datago/io"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Get the path to testdata.xlsx (in parent directory: tests/)
	_, filename, _, _ := runtime.Caller(0)
	benchmarkDir := filepath.Dir(filename)
	parentDir := filepath.Dir(benchmarkDir)
	excelFile := filepath.Join(parentDir, "testdatalarge.xlsx")

	iterations := 5

	fmt.Printf("Benchmarking Excel read: %s\n", excelFile)
	fmt.Printf("Iterations: %d\n", iterations)
	fmt.Println("----------------------------------------")

	var times []time.Duration
	var excelizeTimes []time.Duration

	for i := 0; i < iterations; i++ {
		fmt.Printf("\n--- Run %d ---\n", i+1)

		start := time.Now()
		df, err := io.ReadExcel(excelFile, io.ExcelOptions{HasHeader: true})
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("datago: %.4fs, shape: %v\n", elapsed.Seconds(), df.Shape())
		times = append(times, elapsed)

		startExcelize := time.Now()
		rowsCount, colsCount, err := readWithExcelize(excelFile)
		elapsedExcelize := time.Since(startExcelize)
		if err != nil {
			fmt.Printf("excelize error: %v\n", err)
			return
		}
		fmt.Printf("excelize: %.4fs, shape: [%d %d]\n", elapsedExcelize.Seconds(), rowsCount, colsCount)
		excelizeTimes = append(excelizeTimes, elapsedExcelize)
	}

	// Calculate average
	var total time.Duration
	for _, t := range times {
		total += t
	}
	avg := total / time.Duration(len(times))

	var excelizeTotal time.Duration
	for _, t := range excelizeTimes {
		excelizeTotal += t
	}
	excelizeAvg := excelizeTotal / time.Duration(len(excelizeTimes))

	fmt.Println("\n========================================")
	fmt.Println("Results:")
	fmt.Printf("datago avg: %.4fs\n", avg.Seconds())
	fmt.Printf("excelize avg: %.4fs\n", excelizeAvg.Seconds())
}

func readWithExcelize(path string) (int, int, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		_ = file.Close()
	}()

	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return 0, 0, fmt.Errorf("no sheets found")
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return 0, 0, err
	}

	rowCount := len(rows)
	colCount := 0
	for _, row := range rows {
		if len(row) > colCount {
			colCount = len(row)
		}
	}

	return rowCount, colCount, nil
}
