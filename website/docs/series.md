---
sidebar_position: 4
title: Series 使用指南
---

# Series 使用指南

`Series` 是一维标签数组，可以存储任何数据类型。每个 `Series` 都有一个名称和索引。

## 创建 Series

### 从类型化切片创建

```go
import "github.com/datago/dataframe"

// 从整数切片
intSeries := dataframe.NewSeriesFromInts([]int{1, 2, 3, 4, 5}, "numbers")

// 从 int64 切片
int64Series := dataframe.NewSeriesFromInt64s([]int64{1, 2, 3}, "ids")

// 从浮点数切片
floatSeries := dataframe.NewSeriesFromFloat64s([]float64{1.1, 2.2, 3.3}, "values")

// 从字符串切片
strSeries := dataframe.NewSeriesFromStrings([]string{"a", "b", "c"}, "letters")

// 从布尔切片
boolSeries := dataframe.NewSeriesFromBools([]bool{true, false, true}, "flags")
```

### 从 interface{} 切片创建

```go
// 自动推断类型
data := []interface{}{1, 2, 3, 4, 5}
s := dataframe.NewSeries(data, "mixed")

// 带自定义索引
index := dataframe.NewIndex([]interface{}{"a", "b", "c"}, "label")
s := dataframe.NewSeriesWithIndex(data, "values", index)
```

## 基本属性

```go
s := dataframe.NewSeriesFromInts([]int{10, 20, 30, 40, 50}, "scores")

// 名称
name := s.Name()           // "scores"
s.SetName("new_scores")    // 修改名称

// 数据类型
dtype := s.DType()         // DTypeInt64

// 长度
length := s.Len()          // 5

// 获取所有值
values := s.Values()       // []interface{}{10, 20, 30, 40, 50}

// 获取索引
index := s.Index()
```

## 数据访问

### 按位置访问

```go
// 获取单个值
val, err := s.Get(0)       // 10

// 设置值
err := s.Set(0, 100)

// 切片
head := s.Head(3)          // 前 3 个元素
tail := s.Tail(3)          // 后 3 个元素
slice := s.Slice(1, 4)     // 索引 1 到 3
```

### 按标签访问

```go
// 使用标签获取值
val, err := s.At("label_a")
```

## 统计方法

```go
s := dataframe.NewSeriesFromFloat64s([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "nums")

// 基本统计
sum := s.Sum()         // 55
mean := s.Mean()       // 5.5
median := s.Median()   // 5.5
std := s.Std()         // 标准差
variance := s.Var()    // 方差

// 极值
min := s.Min()         // 1
max := s.Max()         // 10

// 计数
count := s.Count()     // 非空值数量

// 唯一值
unique := s.Unique()   // 返回唯一值 Series
nunique := s.NUnique() // 唯一值数量

// 值计数
counts := s.ValueCounts() // 每个值的出现次数
```

## 数据变换

### Apply - 元素级变换

```go
// 对每个元素应用函数
doubled := s.Apply(func(v interface{}) interface{} {
    if f, ok := v.(float64); ok {
        return f * 2
    }
    return v
})

// 并行 Apply（大数据量推荐）
doubled := s.ParallelApply(func(v interface{}) interface{} {
    if f, ok := v.(float64); ok {
        return f * 2
    }
    return v
})
```

### Map - 映射替换

```go
// 使用映射表替换值
mapping := map[interface{}]interface{}{
    "A": "优秀",
    "B": "良好",
    "C": "及格",
}
grades := s.Map(mapping)
```

### 类型转换

```go
// 转换为 float64
floatSeries, err := s.AsType(dataframe.DTypeFloat64)

// 转换为 string
strSeries, err := s.AsType(dataframe.DTypeString)
```

### 排序

```go
// 升序排序
sorted := s.SortValues(true)

// 降序排序
sorted := s.SortValues(false)
```

## 缺失值处理

```go
s := dataframe.NewSeries([]interface{}{1, nil, 3, nil, 5}, "data")

// 检测缺失值
isNA := s.IsNA()       // [false, true, false, true, false]
notNA := s.NotNA()     // [true, false, true, false, true]

// 填充缺失值
filled := s.FillNA(0)  // [1, 0, 3, 0, 5]

// 删除缺失值
dropped := s.DropNA()  // [1, 3, 5]
```

## 算术运算

支持与标量或另一个 Series 进行运算：

```go
s := dataframe.NewSeriesFromFloat64s([]float64{10, 20, 30}, "values")

// 与标量运算
added := s.Add(5)      // [15, 25, 35]
subbed := s.Sub(5)     // [5, 15, 25]
mulled := s.Mul(2)     // [20, 40, 60]
divided := s.Div(2)    // [5, 10, 15]

// 与另一个 Series 运算
s2 := dataframe.NewSeriesFromFloat64s([]float64{1, 2, 3}, "other")
result := s.Add(s2)    // [11, 22, 33]
result := s.Mul(s2)    // [10, 40, 90]
```

## 复制

```go
// 创建副本
sCopy := s.Copy()
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 创建成绩数据
    scores := dataframe.NewSeriesFromFloat64s(
        []float64{85, 92, 78, 95, 88, 76, 91, 83},
        "scores",
    )

    fmt.Println("成绩数据:")
    fmt.Println(scores)

    // 统计分析
    fmt.Printf("\n统计信息:\n")
    fmt.Printf("  总分: %.1f\n", scores.Sum())
    fmt.Printf("  平均: %.1f\n", scores.Mean())
    fmt.Printf("  中位数: %.1f\n", scores.Median())
    fmt.Printf("  标准差: %.2f\n", scores.Std())
    fmt.Printf("  最低: %v\n", scores.Min())
    fmt.Printf("  最高: %v\n", scores.Max())

    // 成绩分级
    grades := scores.Apply(func(v interface{}) interface{} {
        score := v.(float64)
        switch {
        case score >= 90:
            return "A"
        case score >= 80:
            return "B"
        case score >= 70:
            return "C"
        default:
            return "D"
        }
    })
    fmt.Println("\n等级分布:")
    fmt.Println(grades.ValueCounts())

    // 标准化分数（z-score）
    mean := scores.Mean()
    std := scores.Std()
    zScores := scores.Sub(mean).Div(std)
    fmt.Println("\n标准化分数:")
    fmt.Println(zScores)
}
```

## 相关章节

- [DataFrame 使用指南](./dataframe) - Series 是 DataFrame 的组成部分
- [Index 使用指南](./index) - Series 的索引系统
- [并行处理](./parallel) - 大数据量的并行操作
