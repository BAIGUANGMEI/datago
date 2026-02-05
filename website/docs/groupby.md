---
sidebar_position: 5
title: GroupBy 分组聚合
---

# GroupBy 分组聚合

`GroupBy` 功能允许你按一个或多个列对 DataFrame 进行分组，然后对每个分组应用聚合函数。这是数据分析中最常用的操作之一，类似于 SQL 的 `GROUP BY`。

## 基本用法

### 创建 GroupBy 对象

```go
import "github.com/datago/dataframe"

// 创建示例数据
df, _ := dataframe.New(map[string][]interface{}{
    "region":   {"East", "East", "West", "West", "East"},
    "product":  {"A", "B", "A", "B", "A"},
    "sales":    {100.0, 150.0, 200.0, 120.0, 180.0},
    "quantity": {10, 15, 20, 12, 18},
})

// 按单列分组
gb, err := df.GroupBy("region")
if err != nil {
    // 处理错误（如列不存在）
}

// 按多列分组
gb, err := df.GroupBy("region", "product")
```

### 查看分组信息

```go
// 分组数量
numGroups := gb.NGroups() // 2 (East, West)

// 每组大小
sizeDF := gb.Size()
fmt.Println(sizeDF)
// region  size
// East    3
// West    2
```

## 聚合方法

### 内置聚合函数

| 方法 | 说明 | 示例 |
|------|------|------|
| `Sum()` | 求和 | `gb.Sum("sales")` |
| `Mean()` | 均值 | `gb.Mean("sales")` |
| `Min()` | 最小值 | `gb.Min("sales")` |
| `Max()` | 最大值 | `gb.Max("sales")` |
| `Count()` | 计数 | `gb.Count("sales")` |
| `Std()` | 标准差 | `gb.Std("sales")` |
| `First()` | 首个值 | `gb.First("product")` |
| `Last()` | 末尾值 | `gb.Last("product")` |

```go
// 计算每个地区的销售总额
sumDF := gb.Sum("sales")
fmt.Println(sumDF)
// region  sales_sum
// East    430
// West    320

// 计算每个地区的平均销量
meanDF := gb.Mean("quantity")

// 不指定列名时，对所有非分组列进行聚合
allSums := gb.Sum()
```

### 自定义多重聚合

使用 `Agg` 方法同时应用多个聚合函数：

```go
aggFuncs := map[string][]dataframe.AggFunc{
    "sales":    {dataframe.AggSum, dataframe.AggMean, dataframe.AggMax},
    "quantity": {dataframe.AggSum, dataframe.AggCount},
}

result, err := gb.Agg(aggFuncs)
// 结果包含: sales_0(sum), sales_1(mean), sales_2(max), quantity_0(sum), quantity_1(count)
```

### 预定义聚合函数

| 函数 | 说明 |
|------|------|
| `AggSum` | 求和 |
| `AggMean` | 均值 |
| `AggMin` | 最小值 |
| `AggMax` | 最大值 |
| `AggCount` | 非空计数 |
| `AggStd` | 标准差 |
| `AggVar` | 方差 |
| `AggFirst` | 第一个值 |
| `AggLast` | 最后一个值 |

## 高级操作

### Apply - 自定义分组函数

对每个分组应用任意函数：

```go
// 获取每组销售额最高的记录
result := gb.Apply(func(groupDF *dataframe.DataFrame) *dataframe.DataFrame {
    sorted := groupDF.SortBy("sales", dataframe.Descending)
    return sorted.Head(1)
})
```

### Filter - 分组过滤

根据条件过滤整个分组：

```go
// 只保留销售记录数 >= 2 的地区
filtered := gb.Filter(func(groupDF *dataframe.DataFrame) bool {
    return groupDF.Shape()[0] >= 2
})
```

### Transform - 分组转换

对分组数据进行转换，保持原始索引结构：

```go
// 计算每条记录相对于组内均值的偏差
transformed, err := gb.Transform("sales", func(s *dataframe.Series) *dataframe.Series {
    mean := s.Mean()
    return s.Sub(mean)
})
```

## 并行聚合

对于大数据量，使用并行聚合提升性能：

```go
aggFuncs := map[string][]dataframe.AggFunc{
    "sales":    {dataframe.AggSum, dataframe.AggMean},
    "quantity": {dataframe.AggSum},
}

result, err := gb.ParallelAgg(aggFuncs, dataframe.ParallelOptions{
    NumWorkers: 4,   // 使用 4 个工作协程
    ChunkSize:  100, // 每个工作块最小 100 个分组
})
```

## Concat - 合并 DataFrame

垂直合并多个 DataFrame（类似 SQL UNION）：

```go
df1, _ := dataframe.New(map[string][]interface{}{
    "name": {"Alice", "Bob"},
    "age":  {25, 30},
})

df2, _ := dataframe.New(map[string][]interface{}{
    "name": {"Charlie", "David"},
    "age":  {35, 40},
})

// 合并
combined := dataframe.Concat(df1, df2)
// name     age
// Alice    25
// Bob      30
// Charlie  35
// David    40
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 电商销售数据
    df, _ := dataframe.New(map[string][]interface{}{
        "category": {"电子", "电子", "服装", "服装", "电子", "服装"},
        "product":  {"手机", "电脑", "T恤", "裤子", "耳机", "外套"},
        "sales":    {5000.0, 8000.0, 200.0, 300.0, 500.0, 800.0},
        "quantity": {10, 5, 50, 30, 100, 20},
    })

    fmt.Println("=== 原始数据 ===")
    fmt.Println(df)

    // 1. 按类别分组统计
    gb, _ := df.GroupBy("category")
    
    fmt.Println("\n=== 按类别统计 ===")
    fmt.Println("分组数:", gb.NGroups())
    fmt.Println(gb.Size())

    // 2. 各类别销售总额
    fmt.Println("\n=== 各类别销售总额 ===")
    fmt.Println(gb.Sum("sales"))

    // 3. 各类别多指标统计
    fmt.Println("\n=== 多指标统计 ===")
    aggFuncs := map[string][]dataframe.AggFunc{
        "sales":    {dataframe.AggSum, dataframe.AggMean},
        "quantity": {dataframe.AggSum, dataframe.AggMax},
    }
    stats, _ := gb.Agg(aggFuncs)
    fmt.Println(stats)

    // 4. 筛选销售额 > 1000 的类别
    fmt.Println("\n=== 高销售额类别 ===")
    highSales := gb.Filter(func(g *dataframe.DataFrame) bool {
        s, _ := g.GetSeries("sales")
        return s.Sum() > 1000
    })
    fmt.Println(highSales)

    // 5. 获取各类别销售额最高的产品
    fmt.Println("\n=== 各类别销售冠军 ===")
    topProducts := gb.Apply(func(g *dataframe.DataFrame) *dataframe.DataFrame {
        return g.SortBy("sales", dataframe.Descending).Head(1)
    })
    fmt.Println(topProducts)
}
```

## 相关章节

- [DataFrame 使用指南](./dataframe) - GroupBy 的数据来源
- [Merge/Join 数据合并](./merge) - 合并不同数据源
- [并行处理](./parallel) - 更多并行操作
