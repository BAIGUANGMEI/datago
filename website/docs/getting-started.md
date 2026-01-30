---
sidebar_position: 2
title: 安装与快速使用
---

## 安装

在你的 Go 项目中直接引用模块即可：

```go
import (
    "github.com/datago/dataframe"
    "github.com/datago/io"
)
```

> 需要 Go 1.24+（见 go.mod）。

## 最小示例

```go
package main

import (
    "fmt"

    "github.com/datago/dataframe"
)

func main() {
    df, _ := dataframe.New(map[string][]interface{}{
        "name": {"alice", "bob"},
        "age":  {int64(30), int64(25)},
    })

    fmt.Println(df.Head(1))
}
```

## DataFrame 与 Series

- `DataFrame`：二维表格结构（列式存储）
- `Series`：单列数据（带索引）
- `Index`：行索引系统

建议先阅读 DataFrame 与 Series 的使用章节。
