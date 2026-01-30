---
sidebar_position: 4
title: Series 使用指南
---

## 创建 Series

```go
import "github.com/datago/dataframe"

s1 := dataframe.NewSeriesFromStrings([]string{"a", "b"}, "letter")
s2 := dataframe.NewSeriesFromInts([]int{1, 2, 3}, "num")
```

## 基本操作

- `s.Name()` / `s.SetName(name)`
- `s.DType()`：数据类型推断结果
- `s.Index()` / `s.SetIndex(index)`
- `s.Len()` / `s.Values()`
- `s.Get(pos)` / `s.At(label)` / `s.Set(pos, value)`
- `s.Copy()` / `s.Head(n)` / `s.Tail(n)` / `s.Slice(start, end)`

## 统计与聚合

- `s.Sum()` / `s.Mean()` / `s.Median()` / `s.Std()` / `s.Var()`
- `s.Min()` / `s.Max()` / `s.Count()`
- `s.Unique()` / `s.NUnique()` / `s.ValueCounts()`

## 变换与缺失值

- `s.Apply(fn)`：按元素映射
- `s.Map(mapping)`：用映射表替换值
- `s.FillNA(value)`：填充缺失
- `s.DropNA()`：删除缺失
- `s.IsNA()` / `s.NotNA()`：缺失值标记
- `s.AsType(dtype)`：类型转换
- `s.SortValues(ascending)`：排序

## 算术运算

支持与标量或其他 `Series` 进行逐元素运算：

- `s.Add(other)`
- `s.Sub(other)`
- `s.Mul(other)`
- `s.Div(other)`
