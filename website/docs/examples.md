---
sidebar_position: 7
title: 示例
---

## 1. 条件筛选与排序

```go
import (
    "github.com/datago/dataframe"
)

// 过滤年龄 >= 30，并按年龄降序排序
filtered := df.Filter(func(r dataframe.Row) bool {
    v := r.Get("age")
    if v == nil {
        return false
    }
    age, _ := v.(int64)
    return age >= 30
}).SortBy("age", dataframe.Descending)
```

## 2. 描述性统计

```go
stats := df.Describe()
fmt.Println(stats)
```

## 3. Series 缺失值处理

```go
s := dataframe.NewSeriesFromStrings([]string{"a", "", "c"}, "col")
filled := s.FillNA("N/A")
```

## 4. 写入 Excel

```go
import "github.com/datago/io"

_ = io.WriteExcel("output.xlsx", df, io.ExcelWriteOptions{IncludeIndex: false})
```
