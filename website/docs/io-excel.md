---
sidebar_position: 6
title: Excel 读写
---

## 读取 Excel

```go
import (
    "github.com/datago/dataframe"
    "github.com/datago/io"
)

df, err := io.ReadExcel("testdata.xlsx", io.ExcelOptions{
    Sheet:     "Sheet1",
    HasHeader: true,
    SkipRows:  0,
    UseCols:   []string{"name", "age"},
    DTypes: map[string]dataframe.DType{
        "age": dataframe.DTypeInt64,
    },
})
```

### 读取选项 `ExcelOptions`

- `Sheet`：工作表名称（为空默认第一个）
- `HasHeader`：首行是否为表头
- `SkipRows`：跳过前 N 行
- `UseCols`：只读取指定列名
- `DTypes`：按列强制类型

## 写入 DataFrame

```go
err := io.WriteExcel("output.xlsx", df, io.ExcelWriteOptions{
    Sheet:        "Sheet1",
    IncludeIndex: false,
})
```

## 写入 Series

```go
err := io.WriteSeriesExcel("series.xlsx", s, io.ExcelWriteOptions{
    Sheet:        "Sheet1",
    IncludeIndex: false,
})
```

### 写入选项 `ExcelWriteOptions`

- `Sheet`：写入的工作表名称（默认 `Sheet1`）
- `IncludeHeader`：是否写表头（默认 `true`）
- `IncludeIndex`：是否写索引列
- `IndexName`：索引列名称（默认 `index`）
