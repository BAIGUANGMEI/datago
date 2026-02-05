# DataGo

A high-performance data analysis library for Go, inspired by Python's pandas.

## Features

- **DataFrame & Series**: Familiar pandas-like API for Go developers
- **High Performance**: 2x faster than pandas for Excel reading
- **GroupBy**: Powerful aggregation operations (Sum, Mean, Count, etc.)
- **Merge/Join**: SQL-like table joins (Inner, Left, Right, Outer)
- **Parallel Processing**: Leverage Go's concurrency for big data
- **Excel & CSV**: Full support for reading and writing data files

## Installation

```bash
go get github.com/datago
```

Requires Go 1.24+

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
    "github.com/datago/io"
)

func main() {
    // Create DataFrame
    df, _ := dataframe.New(map[string][]interface{}{
        "name":   {"Alice", "Bob", "Charlie"},
        "age":    {25, 30, 35},
        "salary": {50000.0, 60000.0, 70000.0},
    })

    // Filter and sort
    filtered := df.Filter(func(r dataframe.Row) bool {
        return r.Get("age").(int) >= 30
    }).SortBy("salary", dataframe.Descending)

    // GroupBy aggregation
    gb, _ := df.GroupBy("age")
    stats := gb.Mean("salary")

    // Read/Write files
    excelDF, _ := io.ReadExcel("data.xlsx", io.ExcelOptions{HasHeader: true})
    _ = io.WriteCSV("output.csv", df, io.CSVWriteOptions{})
}
```

## Benchmark

### ReadExcel Performance

| Dataset | DataGo | pandas | polars |
|---------|--------|--------|--------|
| 15K rows × 11 cols | **0.21s** | 0.51s | 0.10s |
| 271K rows × 16 cols | **5.76s** | 11.01s | 2.18s |

DataGo is approximately **2x faster** than pandas for Excel reading.

## Documentation

Visit [documentation site](./website) for detailed guides:

- [Introduction](./website/docs/intro.md)
- [Getting Started](./website/docs/getting-started.md)
- [DataFrame Guide](./website/docs/dataframe.md)
- [Series Guide](./website/docs/series.md)
- [GroupBy](./website/docs/groupby.md)
- [Merge/Join](./website/docs/merge.md)
- [Parallel Processing](./website/docs/parallel.md)
- [Excel I/O](./website/docs/io-excel.md)
- [CSV I/O](./website/docs/io-csv.md)
- [Examples](./website/docs/examples.md)

## API Overview

### DataFrame Operations

```go
// Creation
df, _ := dataframe.New(data)
df, _ := dataframe.FromRecords(records, columns)

// Selection
df.Head(n) / df.Tail(n)
df.Select("col1", "col2")
df.Filter(func(row Row) bool { ... })
df.ILoc(rowStart, rowEnd, colStart, colEnd)

// Manipulation
df.AddColumn("name", series)
df.Drop("col1", "col2")
df.Rename(map[string]string{"old": "new"})
df.SortBy("col", Ascending)

// Statistics
df.Describe()
df.ParallelSum() / df.ParallelMean()
```

### GroupBy Operations

```go
gb, _ := df.GroupBy("category")
gb.Sum("value")
gb.Mean("value")
gb.Agg(map[string][]AggFunc{...})
gb.Apply(func(*DataFrame) *DataFrame { ... })
gb.Filter(func(*DataFrame) bool { ... })
```

### Merge/Join Operations

```go
// Inner/Left/Right/Outer Join
result, _ := dataframe.Merge(left, right, MergeOptions{
    How: InnerJoin,
    On:  []string{"key"},
})

// Different column names
result, _ := left.MergeOn(right, 
    []string{"left_key"}, 
    []string{"right_key"}, 
    LeftJoin,
)
```

### Parallel Processing

```go
// Parallel Apply
result := series.ParallelApply(func(v interface{}) interface{} { ... })

// Parallel Filter
result := df.ParallelFilter(func(row Row) bool { ... })

// Parallel Aggregation
sums := df.ParallelSum()
gb.ParallelAgg(aggFuncs)
```

### I/O Operations

```go
// Excel
df, _ := io.ReadExcel("file.xlsx", ExcelOptions{HasHeader: true})
io.WriteExcel("output.xlsx", df, ExcelWriteOptions{})

// CSV
df, _ := io.ReadCSV("file.csv", CSVOptions{HasHeader: true})
io.WriteCSV("output.csv", df, CSVWriteOptions{})
```

## License

MIT License
