---
sidebar_position: 10
title: CSV 读写
---

# CSV 读写

DataGo 提供完整的 CSV 文件读写功能，支持自定义分隔符、表头处理和类型转换。

## 读取 CSV

### 基本用法

```go
import (
    "github.com/datago/dataframe"
    "github.com/datago/io"
)

// 最简单的读取
df, err := io.ReadCSV("data.csv", io.CSVOptions{
    HasHeader: true,
})
if err != nil {
    // 处理错误
}
fmt.Println(df.Head(10))
```

### 读取选项

```go
df, err := io.ReadCSV("data.csv", io.CSVOptions{
    Separator: ',',                               // 分隔符（默认逗号）
    HasHeader: true,                              // 首行是否为表头
    SkipRows:  1,                                 // 跳过前 N 行
    UseCols:   []string{"name", "age", "email"},  // 只读取指定列
    DTypes: map[string]dataframe.DType{           // 指定列类型
        "age":   dataframe.DTypeInt64,
        "score": dataframe.DTypeFloat64,
    },
})
```

### CSVOptions 详解

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `Separator` | `rune` | `,` | 字段分隔符 |
| `HasHeader` | `bool` | `false` | 首行是否为表头 |
| `SkipRows` | `int` | `0` | 跳过开头的行数 |
| `UseCols` | `[]string` | 全部列 | 只读取指定列 |
| `DTypes` | `map[string]DType` | 自动推断 | 强制指定列的数据类型 |

### 读取不同分隔符的文件

```go
// Tab 分隔的 TSV 文件
df, _ := io.ReadCSV("data.tsv", io.CSVOptions{
    Separator: '\t',
    HasHeader: true,
})

// 分号分隔（欧洲常用）
df, _ := io.ReadCSV("data.csv", io.CSVOptions{
    Separator: ';',
    HasHeader: true,
})

// 管道符分隔
df, _ := io.ReadCSV("data.txt", io.CSVOptions{
    Separator: '|',
    HasHeader: true,
})
```

### 跳过行

```go
// 跳过前 3 行（如注释或元数据）
df, _ := io.ReadCSV("data.csv", io.CSVOptions{
    SkipRows:  3,
    HasHeader: true, // 第 4 行是表头
})
```

## 写入 DataFrame

### 基本用法

```go
df, _ := dataframe.New(map[string][]interface{}{
    "name":   {"Alice", "Bob", "Charlie"},
    "age":    {25, 30, 35},
    "email":  {"alice@example.com", "bob@example.com", "charlie@example.com"},
})

err := io.WriteCSV("output.csv", df, io.CSVWriteOptions{})
if err != nil {
    // 处理错误
}
```

### 写入选项

```go
trueVal := true

err := io.WriteCSV("output.csv", df, io.CSVWriteOptions{
    Separator:     ',',       // 分隔符
    IncludeHeader: &trueVal,  // 是否写入表头
    IncludeIndex:  true,      // 是否写入索引列
    IndexName:     "id",      // 索引列名称
})
```

### CSVWriteOptions 详解

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `Separator` | `rune` | `,` | 字段分隔符 |
| `IncludeHeader` | `*bool` | `true` | 是否写入表头 |
| `IncludeIndex` | `bool` | `false` | 是否写入索引列 |
| `IndexName` | `string` | `"index"` | 索引列的名称 |

## 写入 Series

```go
s := dataframe.NewSeriesFromFloat64s(
    []float64{100.5, 200.3, 300.8, 400.2},
    "values",
)

err := io.WriteSeriesCSV("series.csv", s, io.CSVWriteOptions{
    IncludeIndex: true,
    IndexName:    "row",
})
```

## 并行读取多个文件

```go
paths := []string{"data1.csv", "data2.csv", "data3.csv", "data4.csv"}

combined, err := dataframe.ParallelReadCSV(paths, func(path string) (*dataframe.DataFrame, error) {
    return io.ReadCSV(path, io.CSVOptions{HasHeader: true})
})
if err != nil {
    // 处理错误
}

fmt.Printf("合并后: %v 行\n", combined.Shape()[0])
```

## 完整示例

```go
package main

import (
    "fmt"

    "github.com/datago/dataframe"
    "github.com/datago/io"
)

func main() {
    // === 读取 CSV ===
    fmt.Println("=== 读取 CSV ===")
    
    df, err := io.ReadCSV("sales.csv", io.CSVOptions{
        Separator: ',',
        HasHeader: true,
        DTypes: map[string]dataframe.DType{
            "amount": dataframe.DTypeFloat64,
            "qty":    dataframe.DTypeInt64,
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("原始数据:")
    fmt.Println(df.Head(5))
    fmt.Printf("形状: %v\n", df.Shape())

    // === 数据处理 ===
    fmt.Println("\n=== 数据处理 ===")
    
    // 筛选金额 > 1000 的记录
    filtered := df.Filter(func(r dataframe.Row) bool {
        amount := r.Get("amount")
        if v, ok := amount.(float64); ok {
            return v > 1000
        }
        return false
    })
    fmt.Printf("金额 > 1000 的记录: %d 条\n", filtered.Shape()[0])

    // 按产品分组统计
    gb, _ := df.GroupBy("product")
    stats := gb.Sum("amount")
    fmt.Println("\n各产品销售总额:")
    fmt.Println(stats)

    // === 写入 CSV ===
    fmt.Println("\n=== 写入 CSV ===")
    
    // 写入筛选后的数据
    err = io.WriteCSV("filtered_sales.csv", filtered, io.CSVWriteOptions{})
    if err != nil {
        panic(err)
    }
    fmt.Println("已写入 filtered_sales.csv")

    // 写入统计结果（TSV 格式）
    err = io.WriteCSV("product_stats.tsv", stats, io.CSVWriteOptions{
        Separator: '\t',
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("已写入 product_stats.tsv")

    // 写入带索引的数据
    err = io.WriteCSV("sales_with_index.csv", df, io.CSVWriteOptions{
        IncludeIndex: true,
        IndexName:    "row_num",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("已写入 sales_with_index.csv")
}
```

## 处理特殊情况

### 处理空值

```go
// CSV 中的空值会被读取为空字符串
df, _ := io.ReadCSV("data_with_nulls.csv", io.CSVOptions{HasHeader: true})

// 获取某列
s, _ := df.GetSeries("value")

// 填充空值
filled := s.FillNA(0)

// 删除空值行
cleaned := s.DropNA()
```

### 无表头的文件

```go
// 如果没有表头，列名会自动生成为 col_0, col_1, ...
df, _ := io.ReadCSV("no_header.csv", io.CSVOptions{
    HasHeader: false,
})
// 列名: col_0, col_1, col_2, ...

// 可以重命名
df = df.Rename(map[string]string{
    "col_0": "name",
    "col_1": "age",
    "col_2": "city",
})
```

## 性能提示

1. **使用 UseCols**：只读取需要的列
2. **指定 DTypes**：避免类型推断开销
3. **并行读取**：多个文件使用 `ParallelReadCSV`
4. **流式写入**：大数据量考虑分批写入

## 相关章节

- [Excel 读写](./io-excel) - 另一种常用格式
- [DataFrame 使用指南](./dataframe) - 数据处理
- [并行处理](./parallel) - 批量文件处理
