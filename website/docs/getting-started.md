---
sidebar_position: 2
title: 安装与快速使用
---

# 安装与快速使用

## 环境要求

- Go 1.24 或更高版本

## 安装

### 方式一：直接引用

在你的 Go 项目中直接导入：

```go
import (
    "github.com/datago/dataframe"
    "github.com/datago/io"
)
```

### 方式二：使用 go get

```bash
go get github.com/datago
```

## 核心概念

DataGo 提供三个核心数据结构：

| 结构 | 说明 | 类比 pandas |
|------|------|-------------|
| `DataFrame` | 二维表格数据（列式存储） | `pd.DataFrame` |
| `Series` | 一维标签数组 | `pd.Series` |
| `Index` | 行索引系统 | `pd.Index` |

## 最小示例

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 创建 DataFrame
    df, err := dataframe.New(map[string][]interface{}{
        "name": {"Alice", "Bob", "Charlie"},
        "age":  {25, 30, 35},
    })
    if err != nil {
        panic(err)
    }

    // 打印 DataFrame
    fmt.Println(df)

    // 查看形状
    fmt.Printf("Shape: %v\n", df.Shape()) // [3, 2]

    // 获取前 2 行
    fmt.Println(df.Head(2))

    // 选择列
    fmt.Println(df.Select("name"))
}
```

## 数据类型

DataGo 支持以下数据类型：

| 类型 | Go 类型 | 说明 |
|------|---------|------|
| `DTypeInt64` | `int64` | 64位整数 |
| `DTypeFloat64` | `float64` | 64位浮点数 |
| `DTypeString` | `string` | 字符串 |
| `DTypeBool` | `bool` | 布尔值 |
| `DTypeDateTime` | `time.Time` | 日期时间 |
| `DTypeObject` | `interface{}` | 任意类型 |

类型会根据数据自动推断，也可以手动指定：

```go
// 读取时指定类型
df, _ := io.ReadCSV("data.csv", io.CSVOptions{
    HasHeader: true,
    DTypes: map[string]dataframe.DType{
        "age":    dataframe.DTypeInt64,
        "salary": dataframe.DTypeFloat64,
    },
})

// Series 类型转换
s, _ := series.AsType(dataframe.DTypeFloat64)
```

## 常用操作速查

### DataFrame 操作

```go
// 创建
df, _ := dataframe.New(data)
df, _ := dataframe.FromRecords(records, columns)

// 查看
df.Head(n)           // 前 n 行
df.Tail(n)           // 后 n 行
df.Shape()           // [rows, cols]
df.Columns()         // 列名列表

// 选择
df.Select("col1", "col2")  // 选择列
df.ILoc(0, 5, 0, 2)        // 按位置切片
df.At(rowLabel, "col")     // 单个值

// 操作
df.Filter(fn)              // 条件过滤
df.SortBy("col", Ascending) // 排序
df.Drop("col")             // 删除列
df.Rename(mapping)         // 重命名

// 分组
gb, _ := df.GroupBy("col")
gb.Sum("value")
gb.Mean("value")

// 合并
dataframe.Merge(left, right, opts)
dataframe.Concat(df1, df2)
```

### Series 操作

```go
// 创建
s := dataframe.NewSeriesFromInts([]int{1, 2, 3}, "nums")
s := dataframe.NewSeriesFromStrings([]string{"a", "b"}, "chars")

// 统计
s.Sum()
s.Mean()
s.Min()
s.Max()

// 变换
s.Apply(fn)
s.FillNA(value)
s.DropNA()

// 算术
s.Add(other)
s.Sub(other)
s.Mul(other)
s.Div(other)
```

### I/O 操作

```go
// Excel
df, _ := io.ReadExcel("file.xlsx", io.ExcelOptions{HasHeader: true})
io.WriteExcel("output.xlsx", df, io.ExcelWriteOptions{})

// CSV
df, _ := io.ReadCSV("file.csv", io.CSVOptions{HasHeader: true})
io.WriteCSV("output.csv", df, io.CSVWriteOptions{})
```

## 下一步

建议按以下顺序阅读文档：

1. [DataFrame 使用指南](./dataframe) - 核心数据结构
2. [Series 使用指南](./series) - 一维数据操作
3. [GroupBy 分组聚合](./groupby) - 数据聚合分析
4. [Merge/Join 数据合并](./merge) - 多表关联
5. [并行处理](./parallel) - 大数据加速
