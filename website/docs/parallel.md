---
sidebar_position: 7
title: 并行处理
---

# 并行处理

DataGo 充分利用 Go 语言的并发特性，提供多种并行处理方法，可显著加速大数据量操作。

## 并行选项

```go
type ParallelOptions struct {
    NumWorkers int  // 工作协程数，0 = 自动（CPU 核心数）
    ChunkSize  int  // 每个工作块的最小大小
}

// 使用默认选项
opts := dataframe.DefaultParallelOptions()
// NumWorkers: 0 (自动检测)
// ChunkSize: 1000

// 自定义选项
opts := dataframe.ParallelOptions{
    NumWorkers: 8,
    ChunkSize:  500,
}
```

## Series 并行操作

### ParallelApply

对每个元素并行应用函数：

```go
// 创建大数据量 Series
data := make([]interface{}, 1000000)
for i := range data {
    data[i] = float64(i)
}
s := dataframe.NewSeries(data, "values")

// 并行计算平方
result := s.ParallelApply(func(v interface{}) interface{} {
    if f, ok := v.(float64); ok {
        return f * f
    }
    return v
})

// 使用自定义选项
result := s.ParallelApply(squareFunc, dataframe.ParallelOptions{
    NumWorkers: 4,
    ChunkSize:  10000,
})
```

### ChunkedApply

分块处理，节省内存：

```go
result := s.ChunkedApply(func(chunk []interface{}) []interface{} {
    out := make([]interface{}, len(chunk))
    for i, v := range chunk {
        if f, ok := v.(float64); ok {
            out[i] = f * 2
        }
    }
    return out
}, 50000) // 每块 50000 元素
```

### ParallelChunkedApply

分块 + 并行，适合超大数据：

```go
result := s.ParallelChunkedApply(func(chunk []interface{}) []interface{} {
    // 处理逻辑
    return processedChunk
}, 50000, dataframe.ParallelOptions{NumWorkers: 4})
```

## DataFrame 并行操作

### ParallelFilter

并行条件过滤：

```go
result := df.ParallelFilter(func(row dataframe.Row) bool {
    val := row.Get("value")
    if v, ok := val.(int); ok {
        return v > 1000
    }
    return false
})
```

### ParallelTransform

并行转换所有列：

```go
// 所有数值列乘以 2
result := df.ParallelTransform(func(s *dataframe.Series) *dataframe.Series {
    return s.Mul(2.0)
})
```

### 并行聚合

```go
// 并行计算各列求和
sums := df.ParallelSum()
fmt.Println(sums["sales"]) // 某列的总和

// 并行计算各列均值
means := df.ParallelMean()

// 并行计算最小/最大值
mins := df.ParallelMin()
maxs := df.ParallelMax()
```

## GroupBy 并行聚合

```go
gb, _ := df.GroupBy("category")

aggFuncs := map[string][]dataframe.AggFunc{
    "sales":    {dataframe.AggSum, dataframe.AggMean},
    "quantity": {dataframe.AggSum},
}

result, _ := gb.ParallelAgg(aggFuncs, dataframe.ParallelOptions{
    NumWorkers: 4,
})
```

## 批量处理

### ParallelMapSeries

并行处理多个 Series：

```go
seriesList := []*dataframe.Series{s1, s2, s3, s4, s5}

results := dataframe.ParallelMapSeries(seriesList, func(s *dataframe.Series) *dataframe.Series {
    return s.Mul(2.0)
})
```

### ParallelReadCSV

并行读取多个文件：

```go
paths := []string{"data1.csv", "data2.csv", "data3.csv", "data4.csv"}

combined, err := dataframe.ParallelReadCSV(paths, func(path string) (*dataframe.DataFrame, error) {
    return io.ReadCSV(path, io.CSVOptions{HasHeader: true})
})
// 返回合并后的 DataFrame
```

## 性能对比

### 何时使用并行

| 数据规模 | 建议 |
|----------|------|
| < 1,000 | 普通方法即可 |
| 1,000 - 10,000 | 可尝试并行 |
| 10,000 - 100,000 | 推荐并行 |
| > 100,000 | 强烈推荐并行 |

| 操作类型 | 并行收益 |
|----------|----------|
| 计算密集型 | 高 |
| I/O 密集型 | 中 |
| 简单操作 | 低（可能有开销） |

### 基准测试

```go
// 运行基准测试
// go test -bench=. -benchmem ./tests/...

func BenchmarkSeriesApply(b *testing.B) {
    s := createLargeSeries(100000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        s.Apply(doubleValue)
    }
}

func BenchmarkSeriesParallelApply(b *testing.B) {
    s := createLargeSeries(100000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        s.ParallelApply(doubleValue)
    }
}
```

## 调优指南

### NumWorkers 设置

```go
// 默认：自动检测 CPU 核心数（通常最优）
opts := dataframe.ParallelOptions{NumWorkers: 0}

// 计算密集型：使用 CPU 核心数
opts := dataframe.ParallelOptions{NumWorkers: runtime.NumCPU()}

// I/O 密集型：可以更高
opts := dataframe.ParallelOptions{NumWorkers: runtime.NumCPU() * 2}
```

### ChunkSize 设置

```go
// 太小：调度开销大
// 太大：并行度不足
// 推荐范围：500 - 10000

// 小数据量，简单计算
opts := dataframe.ParallelOptions{ChunkSize: 500}

// 大数据量，复杂计算
opts := dataframe.ParallelOptions{ChunkSize: 5000}
```

### 内存考虑

```go
// 并行处理会创建数据副本
// 内存受限时，使用分块处理

// 节省内存的方式
result := s.ChunkedApply(processFunc, 10000)

// 或使用更小的 worker 数
opts := dataframe.ParallelOptions{NumWorkers: 2}
```

## 完整示例

```go
package main

import (
    "fmt"
    "runtime"
    "time"

    "github.com/datago/dataframe"
)

func main() {
    // 创建大数据量
    n := 1000000
    data := make([]interface{}, n)
    for i := range data {
        data[i] = float64(i)
    }
    s := dataframe.NewSeries(data, "values")

    fmt.Printf("数据量: %d\n", n)
    fmt.Printf("CPU 核心数: %d\n", runtime.NumCPU())

    // 定义计算函数
    compute := func(v interface{}) interface{} {
        if f, ok := v.(float64); ok {
            // 模拟复杂计算
            return f*f + f*2 + 1
        }
        return v
    }

    // 1. 普通 Apply
    start := time.Now()
    _ = s.Apply(compute)
    normalTime := time.Since(start)
    fmt.Printf("\n普通 Apply: %v\n", normalTime)

    // 2. 并行 Apply
    start = time.Now()
    _ = s.ParallelApply(compute)
    parallelTime := time.Since(start)
    fmt.Printf("并行 Apply: %v\n", parallelTime)

    // 3. 分块并行 Apply
    start = time.Now()
    _ = s.ParallelChunkedApply(func(chunk []interface{}) []interface{} {
        result := make([]interface{}, len(chunk))
        for i, v := range chunk {
            if f, ok := v.(float64); ok {
                result[i] = f*f + f*2 + 1
            }
        }
        return result
    }, 50000)
    chunkedTime := time.Since(start)
    fmt.Printf("分块并行 Apply: %v\n", chunkedTime)

    // 性能对比
    fmt.Printf("\n加速比:\n")
    fmt.Printf("  并行 vs 普通: %.2fx\n", float64(normalTime)/float64(parallelTime))
    fmt.Printf("  分块并行 vs 普通: %.2fx\n", float64(normalTime)/float64(chunkedTime))

    // DataFrame 并行操作示例
    dfData := map[string][]interface{}{
        "a": make([]interface{}, 100000),
        "b": make([]interface{}, 100000),
        "c": make([]interface{}, 100000),
    }
    for i := 0; i < 100000; i++ {
        dfData["a"][i] = float64(i)
        dfData["b"][i] = float64(i * 2)
        dfData["c"][i] = float64(i * 3)
    }
    df, _ := dataframe.New(dfData)

    fmt.Println("\n=== DataFrame 并行聚合 ===")
    start = time.Now()
    sums := df.ParallelSum()
    fmt.Printf("并行求和耗时: %v\n", time.Since(start))
    fmt.Printf("结果: a=%.0f, b=%.0f, c=%.0f\n", sums["a"], sums["b"], sums["c"])
}
```

## 相关章节

- [DataFrame 使用指南](./dataframe) - 了解基本操作
- [Series 使用指南](./series) - 了解 Series 操作
- [GroupBy 分组聚合](./groupby) - 分组并行聚合
