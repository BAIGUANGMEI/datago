---
sidebar_position: 9
title: Excel 读写
---

# Excel 读写

DataGo 提供高性能的 Excel 文件读写功能，基于 [excelize](https://github.com/xuri/excelize) 库实现。

## 性能特点

| 数据规模 | DataGo | pandas |
|----------|--------|--------|
| 15K行 × 11列 | **0.21s** | 0.51s |
| 271K行 × 16列 | **5.76s** | 11.01s |

> DataGo 的 Excel 读取速度约为 pandas 的 **2 倍**。

## 读取 Excel

### 基本用法

```go
import (
    "github.com/datago/dataframe"
    "github.com/datago/io"
)

// 最简单的读取
df, err := io.ReadExcel("data.xlsx", io.ExcelOptions{
    HasHeader: true,
})
if err != nil {
    // 处理错误
}
fmt.Println(df.Head(10))
```

### 读取选项

```go
df, err := io.ReadExcel("data.xlsx", io.ExcelOptions{
    Sheet:     "Sheet1",              // 工作表名（空=第一个）
    HasHeader: true,                  // 首行是否为表头
    SkipRows:  2,                     // 跳过前 N 行
    UseCols:   []string{"name", "age", "salary"}, // 只读取指定列
    DTypes: map[string]dataframe.DType{           // 指定列类型
        "age":    dataframe.DTypeInt64,
        "salary": dataframe.DTypeFloat64,
    },
})
```

### ExcelOptions 详解

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `Sheet` | `string` | 第一个工作表 | 要读取的工作表名称 |
| `HasHeader` | `bool` | `false` | 首行是否为表头 |
| `SkipRows` | `int` | `0` | 跳过开头的行数 |
| `UseCols` | `[]string` | 全部列 | 只读取指定列 |
| `DTypes` | `map[string]DType` | 自动推断 | 强制指定列的数据类型 |

### 读取多个工作表

```go
// 读取第一个工作表
df1, _ := io.ReadExcel("multi_sheet.xlsx", io.ExcelOptions{
    Sheet:     "Sheet1",
    HasHeader: true,
})

// 读取另一个工作表
df2, _ := io.ReadExcel("multi_sheet.xlsx", io.ExcelOptions{
    Sheet:     "销售数据",
    HasHeader: true,
})
```

## 写入 DataFrame

### 基本用法

```go
df, _ := dataframe.New(map[string][]interface{}{
    "name":   {"Alice", "Bob", "Charlie"},
    "age":    {25, 30, 35},
    "salary": {50000.0, 60000.0, 70000.0},
})

err := io.WriteExcel("output.xlsx", df, io.ExcelWriteOptions{})
if err != nil {
    // 处理错误
}
```

### 写入选项

```go
err := io.WriteExcel("output.xlsx", df, io.ExcelWriteOptions{
    Sheet:         "员工数据",    // 工作表名（默认 Sheet1）
    IncludeHeader: &trueVal,     // 是否写入表头（默认 true）
    IncludeIndex:  true,         // 是否写入索引列
    IndexName:     "row_id",     // 索引列名称（默认 "index"）
})

// 注意：IncludeHeader 使用指针，需要这样设置
trueVal := true
falseVal := false
```

### ExcelWriteOptions 详解

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `Sheet` | `string` | `"Sheet1"` | 工作表名称 |
| `IncludeHeader` | `*bool` | `true` | 是否写入表头 |
| `IncludeIndex` | `bool` | `false` | 是否写入索引列 |
| `IndexName` | `string` | `"index"` | 索引列的名称 |

## 写入 Series

```go
s := dataframe.NewSeriesFromFloat64s(
    []float64{100, 200, 300, 400, 500},
    "monthly_sales",
)

err := io.WriteSeriesExcel("series.xlsx", s, io.ExcelWriteOptions{
    Sheet:        "销售",
    IncludeIndex: true,
})
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
    // === 读取 Excel ===
    fmt.Println("=== 读取 Excel ===")
    
    df, err := io.ReadExcel("employees.xlsx", io.ExcelOptions{
        HasHeader: true,
        UseCols:   []string{"name", "department", "salary"},
        DTypes: map[string]dataframe.DType{
            "salary": dataframe.DTypeFloat64,
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("读取数据:")
    fmt.Println(df.Head(5))
    fmt.Printf("形状: %v\n", df.Shape())

    // === 数据处理 ===
    fmt.Println("\n=== 数据处理 ===")
    
    // 按部门分组统计
    gb, _ := df.GroupBy("department")
    stats := gb.Mean("salary")
    fmt.Println("各部门平均薪资:")
    fmt.Println(stats)

    // === 写入 Excel ===
    fmt.Println("\n=== 写入 Excel ===")
    
    // 写入处理后的数据
    err = io.WriteExcel("department_stats.xlsx", stats, io.ExcelWriteOptions{
        Sheet:        "部门统计",
        IncludeIndex: false,
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("已写入 department_stats.xlsx")

    // 写入原始数据（包含索引）
    err = io.WriteExcel("employees_backup.xlsx", df, io.ExcelWriteOptions{
        Sheet:        "员工数据",
        IncludeIndex: true,
        IndexName:    "序号",
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("已写入 employees_backup.xlsx")
}
```

## 错误处理

```go
df, err := io.ReadExcel("data.xlsx", io.ExcelOptions{HasHeader: true})
if err != nil {
    switch {
    case os.IsNotExist(err):
        fmt.Println("文件不存在")
    default:
        fmt.Printf("读取错误: %v\n", err)
    }
    return
}
```

## 性能提示

1. **只读取需要的列**：使用 `UseCols` 减少内存占用
2. **指定数据类型**：使用 `DTypes` 避免类型推断开销
3. **大文件处理**：考虑分批读取或使用流式处理
4. **并行读取**：多个文件可使用 `ParallelReadCSV` 类似方式

## 相关章节

- [CSV 读写](./io-csv) - 另一种常用格式
- [DataFrame 使用指南](./dataframe) - 数据处理
- [并行处理](./parallel) - 批量文件处理
