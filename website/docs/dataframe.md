---
sidebar_position: 3
title: DataFrame 使用指南
---

# DataFrame 使用指南

`DataFrame` 是 DataGo 的核心数据结构，表示二维表格数据，类似于电子表格或 SQL 表。

## 创建 DataFrame

### 从 Map 创建

```go
import "github.com/datago/dataframe"

df, err := dataframe.New(map[string][]interface{}{
    "name":   {"Alice", "Bob", "Charlie"},
    "age":    {25, 30, 35},
    "salary": {50000.0, 60000.0, 70000.0},
})
if err != nil {
    // 处理错误（如列长度不一致）
}
```

### 从二维记录创建

```go
records := [][]interface{}{
    {"Alice", 25, 50000.0},
    {"Bob", 30, 60000.0},
    {"Charlie", 35, 70000.0},
}
columns := []string{"name", "age", "salary"}

df, err := dataframe.FromRecords(records, columns)
```

## 基本信息

```go
// 获取形状 [行数, 列数]
shape := df.Shape() // [3, 3]

// 获取列名
cols := df.Columns() // ["name", "age", "salary"]

// 获取索引
index := df.Index()

// 打印 DataFrame
fmt.Println(df)
// DataFrame: rows=3, cols=3
// index  name     age  salary
// 0      Alice    25   50000
// 1      Bob      30   60000
// 2      Charlie  35   70000
```

## 数据选择

### 选择行

```go
// 前 n 行
first5 := df.Head(5)

// 后 n 行
last5 := df.Tail(5)

// 按位置获取单行
row, err := df.Row(0)
name := row.Get("name") // "Alice"
```

### 选择列

```go
// 选择指定列（返回新 DataFrame）
subset := df.Select("name", "age")

// 获取单列（返回 Series）
series, ok := df.GetSeries("age")
if ok {
    fmt.Println(series.Mean()) // 30
}
```

### 切片操作

```go
// ILoc - 按位置切片 [rowStart:rowEnd, colStart:colEnd]
subset := df.ILoc(0, 2, 0, 2) // 前2行，前2列

// Loc - 按标签选择
subset := df.Loc([]interface{}{0, 1}, []string{"name", "age"})

// At - 获取单个值
value, err := df.At(0, "name") // "Alice"
```

## 数据操作

### 添加/修改列

```go
// 添加新列
bonusSeries := dataframe.NewSeriesFromFloat64s(
    []float64{5000, 6000, 7000}, 
    "bonus",
)
newDF := df.AddColumn("bonus", bonusSeries)

// 设置/替换列
err := df.SetColumn("age", newAgeSeries)
```

### 删除列

```go
// 删除单列
df2 := df.Drop("salary")

// 删除多列
df2 := df.Drop("salary", "bonus")
```

### 重命名列

```go
df2 := df.Rename(map[string]string{
    "name":   "employee_name",
    "salary": "annual_salary",
})
```

### 数据过滤

```go
// 使用 Filter 方法
filtered := df.Filter(func(row dataframe.Row) bool {
    age := row.Get("age")
    if v, ok := age.(int); ok {
        return v >= 30
    }
    return false
})

// 并行过滤（大数据量推荐）
filtered := df.ParallelFilter(func(row dataframe.Row) bool {
    return row.Get("age").(int) >= 30
})
```

### 排序

```go
// 升序排序
sorted := df.SortBy("age", dataframe.Ascending)

// 降序排序
sorted := df.SortBy("salary", dataframe.Descending)
```

## 统计分析

### Describe - 统计摘要

```go
stats := df.Describe()
fmt.Println(stats)
// 输出每列的 count, mean, std, min, max
```

### 并行聚合

```go
// 并行计算所有列的和
sums := df.ParallelSum()
fmt.Println(sums["salary"]) // 180000

// 并行计算均值
means := df.ParallelMean()

// 并行计算最小/最大值
mins := df.ParallelMin()
maxs := df.ParallelMax()
```

## 复制与转换

```go
// 浅拷贝
dfCopy := df.Copy()

// 并行转换所有列
transformed := df.ParallelTransform(func(s *dataframe.Series) *dataframe.Series {
    return s.Mul(2) // 所有数值乘以 2
})
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 创建销售数据
    df, _ := dataframe.New(map[string][]interface{}{
        "product":  {"A", "B", "A", "B", "A"},
        "region":   {"East", "East", "West", "West", "East"},
        "sales":    {100.0, 150.0, 200.0, 120.0, 180.0},
        "quantity": {10, 15, 20, 12, 18},
    })
    fmt.Println("原始数据:")
    fmt.Println(df)

    // 筛选 East 地区
    eastSales := df.Filter(func(r dataframe.Row) bool {
        return r.Get("region") == "East"
    })
    fmt.Println("\nEast 地区销售:")
    fmt.Println(eastSales)

    // 按销售额降序排序
    sorted := df.SortBy("sales", dataframe.Descending)
    fmt.Println("\n按销售额排序:")
    fmt.Println(sorted)

    // 统计摘要
    fmt.Println("\n统计摘要:")
    fmt.Println(df.Describe())

    // 选择特定列
    subset := df.Select("product", "sales")
    fmt.Println("\n产品销售:")
    fmt.Println(subset)
}
```

## 相关章节

- [Series 使用指南](./series) - 单列数据操作
- [GroupBy 分组聚合](./groupby) - 数据分组与聚合
- [Merge/Join 数据合并](./merge) - 多表关联
- [并行处理](./parallel) - 大数据量加速
