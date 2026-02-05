---
sidebar_position: 8
title: Index 使用指南
---

# Index 使用指南

`Index` 是 DataFrame 和 Series 的行索引系统，支持标签访问和集合操作。

## 创建 Index

### 默认整数索引

```go
import "github.com/datago/dataframe"

// 创建 0 到 n-1 的整数索引
index := dataframe.NewRangeIndex(5) // [0, 1, 2, 3, 4]
```

### 自定义标签索引

```go
// 使用字符串标签
labels := []interface{}{"a", "b", "c", "d", "e"}
index := dataframe.NewIndex(labels, "my_index")

// 使用混合类型标签
mixedLabels := []interface{}{2020, 2021, 2022, "latest"}
index := dataframe.NewIndex(mixedLabels, "year")
```

## 基本属性

```go
index := dataframe.NewIndex([]interface{}{"a", "b", "c"}, "letters")

// 长度
length := index.Len() // 3

// 名称
name := index.Name()        // "letters"
index.SetName("new_name")   // 修改名称

// 获取所有标签
labels := index.Labels()    // []interface{}{"a", "b", "c"}
```

## 数据访问

### 按位置获取

```go
// 获取指定位置的标签
label, err := index.Get(0) // "a"
label, err := index.Get(2) // "c"
```

### 按标签查找

```go
// 获取标签的位置
pos, err := index.GetLoc("b") // 1

// 检查标签是否存在
exists := index.Contains("b") // true
exists := index.Contains("x") // false
```

## 切片与修改

### 切片

```go
// 获取子索引
sub := index.Slice(1, 3) // ["b", "c"]

// 与 Series/DataFrame 的切片对应
s := series.Slice(1, 3)
```

### 追加

```go
// 追加新标签
newIndex := index.Append("d") // ["a", "b", "c", "d"]
```

### 复制与重置

```go
// 创建副本
indexCopy := index.Copy()

// 重置为默认整数索引
resetIndex := index.Reset() // [0, 1, 2]
```

## 转换

```go
// 转换为字符串切片
strLabels := index.ToStringSlice() // []string{"a", "b", "c"}
```

## 集合操作

### 相等比较

```go
idx1 := dataframe.NewIndex([]interface{}{1, 2, 3}, "")
idx2 := dataframe.NewIndex([]interface{}{1, 2, 3}, "")
idx3 := dataframe.NewIndex([]interface{}{1, 2, 4}, "")

idx1.Equals(idx2) // true
idx1.Equals(idx3) // false
```

### 并集

```go
idx1 := dataframe.NewIndex([]interface{}{1, 2, 3}, "")
idx2 := dataframe.NewIndex([]interface{}{3, 4, 5}, "")

union := idx1.Union(idx2) // [1, 2, 3, 4, 5]
```

### 交集

```go
idx1 := dataframe.NewIndex([]interface{}{1, 2, 3, 4}, "")
idx2 := dataframe.NewIndex([]interface{}{3, 4, 5, 6}, "")

intersection := idx1.Intersection(idx2) // [3, 4]
```

### 差集

```go
idx1 := dataframe.NewIndex([]interface{}{1, 2, 3, 4}, "")
idx2 := dataframe.NewIndex([]interface{}{3, 4, 5, 6}, "")

diff := idx1.Difference(idx2) // [1, 2] (在 idx1 但不在 idx2)
```

## 与 DataFrame/Series 配合使用

### 创建带自定义索引的 Series

```go
data := []interface{}{100, 200, 300}
index := dataframe.NewIndex([]interface{}{"a", "b", "c"}, "label")

s := dataframe.NewSeriesWithIndex(data, "values", index)

// 按标签访问
val, _ := s.At("b") // 200
```

### 在 DataFrame 中使用

```go
df, _ := dataframe.New(map[string][]interface{}{
    "name":  {"Alice", "Bob", "Charlie"},
    "score": {85, 90, 78},
})

// 使用 Loc 按索引标签选择
subset := df.Loc([]interface{}{0, 2}, nil)

// 使用 At 按标签获取单个值
val, _ := df.At(1, "name") // "Bob"
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 创建月份索引
    months := dataframe.NewIndex(
        []interface{}{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
        "month",
    )

    fmt.Println("=== 基本信息 ===")
    fmt.Printf("长度: %d\n", months.Len())
    fmt.Printf("名称: %s\n", months.Name())
    fmt.Printf("标签: %v\n", months.Labels())

    // 访问操作
    fmt.Println("\n=== 访问操作 ===")
    label, _ := months.Get(2)
    fmt.Printf("位置 2 的标签: %v\n", label)

    pos, _ := months.GetLoc("Apr")
    fmt.Printf("'Apr' 的位置: %d\n", pos)

    fmt.Printf("包含 'Feb': %v\n", months.Contains("Feb"))
    fmt.Printf("包含 'Dec': %v\n", months.Contains("Dec"))

    // 切片操作
    fmt.Println("\n=== 切片操作 ===")
    q1 := months.Slice(0, 3)
    fmt.Printf("Q1 月份: %v\n", q1.Labels())

    // 集合操作
    fmt.Println("\n=== 集合操作 ===")
    set1 := dataframe.NewIndex([]interface{}{"Jan", "Feb", "Mar"}, "")
    set2 := dataframe.NewIndex([]interface{}{"Mar", "Apr", "May"}, "")

    fmt.Printf("集合1: %v\n", set1.Labels())
    fmt.Printf("集合2: %v\n", set2.Labels())
    fmt.Printf("并集: %v\n", set1.Union(set2).Labels())
    fmt.Printf("交集: %v\n", set1.Intersection(set2).Labels())
    fmt.Printf("差集: %v\n", set1.Difference(set2).Labels())

    // 与 Series 配合
    fmt.Println("\n=== 与 Series 配合 ===")
    sales := dataframe.NewSeriesWithIndex(
        []interface{}{1000, 1200, 1100, 1300, 1400, 1500},
        "sales",
        months,
    )
    
    val, _ := sales.At("Mar")
    fmt.Printf("3月销售额: %v\n", val)
}
```

## 相关章节

- [DataFrame 使用指南](./dataframe) - Index 作为 DataFrame 的行索引
- [Series 使用指南](./series) - Index 作为 Series 的索引
