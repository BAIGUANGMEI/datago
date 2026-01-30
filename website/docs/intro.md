---
sidebar_position: 1
title: DataGo 简介
---

DataGo 是一个面向 Go 的数据分析库，提供类似 pandas 的 `DataFrame` / `Series` 结构和常用的数据处理能力，并支持 Excel 读写。

## 功能概览

- **DataFrame/Series/Index**：二维/一维数据结构与索引体系
- **数据选择**：`Select`、`ILoc`、`Loc`、`At`、`Row`
- **常用操作**：`AddColumn`、`Drop`、`Rename`、`SortBy`
- **统计分析**：`Describe`、`Mean`、`Std`、`ValueCounts` 等
- **缺失值处理**：`FillNA`、`DropNA`、`IsNA`、`NotNA`
- **Excel 读写**：`ReadExcel`、`WriteExcel`、`WriteSeriesExcel`

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
    "name": {"alice", "bob"},
    "age":  {int64(30), int64(25)},
  })

  // 读取 Excel
  excelDF, _ := io.ReadExcel("testdata.xlsx", io.ExcelOptions{HasHeader: true})
  fmt.Println(excelDF.Head(5))

  // 写入 Excel
  _ = io.WriteExcel("output.xlsx", df, io.ExcelWriteOptions{IncludeIndex: false})
}
```
