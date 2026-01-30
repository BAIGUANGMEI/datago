---
sidebar_position: 3
title: DataFrame 使用指南
---

## 创建 DataFrame

### 从 map 创建

```go
import "github.com/datago/dataframe"

df, err := dataframe.New(map[string][]interface{}{
    "name": {"alice", "bob"},
    "age":  {int64(30), int64(25)},
})
```

### 从二维记录创建

```go
records := [][]interface{}{
    {"alice", int64(30)},
    {"bob", int64(25)},
}
columns := []string{"name", "age"}

df, err := dataframe.FromRecords(records, columns)
```

## 基本信息

- `df.Columns()`：列名
- `df.Shape()`：返回 `[rows, cols]`
- `df.Index()`：返回索引

## 取值与切片

- `df.Head(n)` / `df.Tail(n)`
- `df.Select("name", "age")`
- `df.ILoc(rowStart, rowEnd, colStart, colEnd)`：位置选择
- `df.Loc(rowLabels, colLabels)`：标签选择
- `df.At(rowLabel, "col")`：按行标签 + 列名取值
- `df.Row(pos)`：按行位置返回 `Row`

示例：

```go
row, _ := df.Row(0)
value := row.Get("name")
```

## 结构操作

- `df.AddColumn(name, series)`：追加列
- `df.Drop(cols...)`：删除列
- `df.Rename(map[string]string)`：重命名列
- `df.SortBy(column, order)`：排序（`Ascending` / `Descending`）

## 统计分析

- `df.Describe()`：输出统计摘要（count/mean/std/min/max）

> `Describe` 只对可数值化列计算统计值。