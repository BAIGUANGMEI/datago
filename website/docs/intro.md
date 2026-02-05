---
sidebar_position: 1
title: DataGo 简介
---

# DataGo

DataGo 是一个高性能的 Go 语言数据分析库，提供类似 Python pandas 的 `DataFrame` 和 `Series` 数据结构，专为 Go 开发者设计。

## 为什么选择 DataGo？

- **熟悉的 API**：借鉴 pandas 的设计理念，降低学习成本
- **高性能**：原生 Go 实现，Excel 读取速度是 pandas 的 2 倍
- **并发支持**：充分利用 Go 的并发特性，支持并行数据处理
- **类型安全**：Go 的静态类型系统提供编译时检查
- **易于集成**：无缝融入现有 Go 项目

## 功能概览

### 核心数据结构

| 结构 | 说明 |
|------|------|
| `DataFrame` | 二维表格数据，支持多种数据类型 |
| `Series` | 一维标签数组，支持统计计算 |
| `Index` | 行索引系统，支持标签访问 |

### 数据操作

- **选择与切片**：`Select`、`ILoc`、`Loc`、`At`、`Head`、`Tail`
- **结构操作**：`AddColumn`、`Drop`、`Rename`、`Filter`、`SortBy`
- **统计分析**：`Describe`、`Mean`、`Std`、`Min`、`Max`、`ValueCounts`
- **缺失值处理**：`FillNA`、`DropNA`、`IsNA`、`NotNA`

### 高级功能

- **GroupBy 分组聚合**：支持 `Sum`、`Mean`、`Count` 等聚合操作
- **Merge/Join 合并**：支持 Inner、Left、Right、Outer Join
- **并行处理**：`ParallelApply`、`ParallelFilter`、`ParallelTransform`
- **数据合并**：`Concat` 垂直合并多个 DataFrame

### 数据 I/O

- **Excel**：`ReadExcel`、`WriteExcel`（基于 excelize）
- **CSV**：`ReadCSV`、`WriteCSV`

## 性能基准

读取 Excel 文件性能对比：

| 数据规模 | DataGo | pandas | polars |
|----------|--------|--------|--------|
| 15K行 × 11列 | **0.21s** | 0.51s | 0.10s |
| 271K行 × 16列 | **5.76s** | 11.01s | 2.18s |

> DataGo 的 Excel 读取性能约为 pandas 的 **2 倍**。

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
    fmt.Println(df)

    // 数据筛选
    filtered := df.Filter(func(r dataframe.Row) bool {
        age, _ := r.Get("age").(int)
        return age >= 30
    })
    fmt.Println("Age >= 30:", filtered)

    // 分组聚合
    gb, _ := df.GroupBy("age")
    stats := gb.Mean("salary")
    fmt.Println("平均薪资:", stats)

    // 读取 Excel
    excelDF, _ := io.ReadExcel("data.xlsx", io.ExcelOptions{HasHeader: true})
    fmt.Println(excelDF.Head(5))

    // 写入 CSV
    _ = io.WriteCSV("output.csv", df, io.CSVWriteOptions{})
}
```

## 下一步

- [安装与快速使用](./getting-started) - 开始使用 DataGo
- [DataFrame 指南](./dataframe) - 学习核心数据结构
- [示例](./examples) - 查看更多使用示例
