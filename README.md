# DataGo

高性能 Go 语言数据分析库，灵感来自 Python pandas。

## 特性

- **DataFrame & Series**：为 Go 开发者提供熟悉的 pandas 风格 API
- **高性能**：Excel 读取速度是 pandas 的 2 倍
- **GroupBy**：强大的分组聚合操作（Sum、Mean、Count 等）
- **Merge/Join**：SQL 风格的表连接（Inner、Left、Right、Outer）
- **并行处理**：利用 Go 并发优势加速大数据处理
- **Excel & CSV**：完整支持数据文件读写

## 安装

```bash
go get github.com/datago
```

需要 Go 1.24+

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
    "github.com/datago/io"
)

func main() {
    // 创建 DataFrame
    df, _ := dataframe.New(map[string][]interface{}{
        "name":   {"Alice", "Bob", "Charlie"},
        "age":    {25, 30, 35},
        "salary": {50000.0, 60000.0, 70000.0},
    })

    // 筛选和排序
    filtered := df.Filter(func(r dataframe.Row) bool {
        return r.Get("age").(int) >= 30
    }).SortBy("salary", dataframe.Descending)

    // GroupBy 聚合
    gb, _ := df.GroupBy("age")
    stats := gb.Mean("salary")

    // 读写文件
    excelDF, _ := io.ReadExcel("data.xlsx", io.ExcelOptions{HasHeader: true})
    _ = io.WriteCSV("output.csv", df, io.CSVWriteOptions{})
}
```

## 性能基准

### ReadExcel 性能测试

| 数据集 | DataGo | pandas | polars |
|--------|--------|--------|--------|
| 15K 行 × 11 列 | **0.21s** | 0.51s | 0.10s |
| 271K 行 × 16 列 | **5.76s** | 11.01s | 2.18s |

DataGo 读取 Excel 速度约为 pandas 的 **2 倍**。

## 文档

在线文档：https://baiguangmei.github.io/datago/

- [简介](./website/docs/intro.md)
- [快速入门](./website/docs/getting-started.md)
- [DataFrame 指南](./website/docs/dataframe.md)
- [Series 指南](./website/docs/series.md)
- [GroupBy 分组聚合](./website/docs/groupby.md)
- [Merge/Join 表连接](./website/docs/merge.md)
- [并行处理](./website/docs/parallel.md)
- [Excel 读写](./website/docs/io-excel.md)
- [CSV 读写](./website/docs/io-csv.md)
- [示例](./website/docs/examples.md)

## API 概览

### DataFrame 操作

```go
// 创建
df, _ := dataframe.New(data)
df, _ := dataframe.FromRecords(records, columns)

// 选择
df.Head(n) / df.Tail(n)
df.Select("col1", "col2")
df.Filter(func(row Row) bool { ... })
df.ILoc(rowStart, rowEnd, colStart, colEnd)

// 操作
df.AddColumn("name", series)
df.Drop("col1", "col2")
df.Rename(map[string]string{"old": "new"})
df.SortBy("col", Ascending)

// 统计
df.Describe()
df.ParallelSum() / df.ParallelMean()
```

### GroupBy 操作

```go
gb, _ := df.GroupBy("category")
gb.Sum("value")
gb.Mean("value")
gb.Agg(map[string][]AggFunc{...})
gb.Apply(func(*DataFrame) *DataFrame { ... })
gb.Filter(func(*DataFrame) bool { ... })
```

### Merge/Join 操作

```go
// Inner/Left/Right/Outer 连接
result, _ := dataframe.Merge(left, right, MergeOptions{
    How: InnerJoin,
    On:  []string{"key"},
})

// 不同列名连接
result, _ := left.MergeOn(right, 
    []string{"left_key"}, 
    []string{"right_key"}, 
    LeftJoin,
)
```

### 并行处理

```go
// 并行 Apply
result := series.ParallelApply(func(v interface{}) interface{} { ... })

// 并行筛选
result := df.ParallelFilter(func(row Row) bool { ... })

// 并行聚合
sums := df.ParallelSum()
gb.ParallelAgg(aggFuncs)
```

### I/O 操作

```go
// Excel
df, _ := io.ReadExcel("file.xlsx", ExcelOptions{HasHeader: true})
io.WriteExcel("output.xlsx", df, ExcelWriteOptions{})

// CSV
df, _ := io.ReadCSV("file.csv", CSVOptions{HasHeader: true})
io.WriteCSV("output.csv", df, CSVWriteOptions{})
```

## 许可证

MIT License
