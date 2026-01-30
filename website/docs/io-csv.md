---
sidebar_position: 7
title: CSV 读写
---

## 读取 CSV

```go
import (
    "github.com/datago/dataframe"
    "github.com/datago/io"
)

df, err := io.ReadCSV("data.csv", io.CSVOptions{
    Separator: ',',
    HasHeader: true,
    SkipRows:  0,
    UseCols:   []string{"name", "age"},
    DTypes: map[string]dataframe.DType{
        "age": dataframe.DTypeInt64,
    },
})
```

### 读取选项 `CSVOptions`

- `Separator`：分隔符（默认 `,`）
- `HasHeader`：首行是否为表头
- `SkipRows`：跳过前 N 行
- `UseCols`：只读取指定列名
- `DTypes`：按列强制类型

## 写入 DataFrame

```go
err := io.WriteCSV("output.csv", df, io.CSVWriteOptions{
    Separator:    ',',
    IncludeIndex: false,
})
```

## 写入 Series

```go
err := io.WriteSeriesCSV("series.csv", s, io.CSVWriteOptions{
    Separator:    ',',
    IncludeIndex: false,
})
```

### 写入选项 `CSVWriteOptions`

- `Separator`：分隔符（默认 `,`）
- `IncludeHeader`：是否写表头（默认 `true`）
- `IncludeIndex`：是否写索引列
- `IndexName`：索引列名称（默认 `index`）
