package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/datago/io"
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
	}

	// Calculate average
	var total time.Duration
	for _, t := range times {
		total += t
	}
	avg := total / time.Duration(len(times))

	fmt.Println("\n========================================")
	fmt.Println("Results:")
	fmt.Printf("datago avg: %.4fs\n", avg.Seconds())
}
