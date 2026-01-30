# DataGo

DataGo 是一个面向 Go 的数据分析库，提供类似 pandas 的 `DataFrame` / `Series` 结构和常用的数据处理能力，并支持 Excel 读写。

## 功能概览

- DataFrame/Series/Index：二维/一维数据结构与索引体系
- 数据选择：`Select`、`ILoc`、`Loc`、`At`、`Row`
- 常用操作：`AddColumn`、`Drop`、`Rename`、`SortBy`
- 统计分析：`Describe`、`Mean`、`Std`、`ValueCounts`
- 缺失值处理：`FillNA`、`DropNA`、`IsNA`、`NotNA`
- Excel 读写：`ReadExcel`、`WriteExcel`、`WriteSeriesExcel`

## 安装

```bash
go get github.com/datago
```

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

## 文档

在线文档：https://baiguangmei.github.io/datago/

## 基准测试

### ReadExcel

以下结果为读取 Excel 的基准测试，均为 5 轮测试的平均耗时（秒）。

#### testdata.xlsx (15,430 rows × 11 columns)

| Language | Package | Rounds | Avg Read Time (s) |
|:--------|:--------|-------:|------------------:|
| Go      | datago  | 5      | 0.2120            |
| Python  | pandas  | 5      | 0.5135            |
| Python  | polars  | 5      | 0.1020            |

#### testdatalarge.xlsx (271,114 rows × 16 columns)

| Language | Package | Rounds | Avg Read Time (s) |
|:--------|:--------|-------:|------------------:|
| Go      | datago  | 5      | 5.7648            |
| Python  | pandas  | 5      | 11.0093           |
| Python  | polars  | 5      | 2.1800            |
