---
sidebar_position: 11
title: 示例集锦
---

# 示例集锦

本页面提供 DataGo 的常见使用场景和完整示例。

## 1. 数据筛选与排序

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    df, _ := dataframe.New(map[string][]interface{}{
        "name":   {"Alice", "Bob", "Charlie", "David", "Eve"},
        "age":    {28, 35, 42, 31, 26},
        "salary": {50000.0, 75000.0, 90000.0, 65000.0, 48000.0},
        "dept":   {"IT", "HR", "IT", "Sales", "IT"},
    })

    // 筛选年龄 >= 30 的员工
    over30 := df.Filter(func(r dataframe.Row) bool {
        age := r.Get("age").(int)
        return age >= 30
    })
    fmt.Println("年龄 >= 30:")
    fmt.Println(over30)

    // 按薪资降序排序
    bySalary := df.SortBy("salary", dataframe.Descending)
    fmt.Println("\n按薪资排序:")
    fmt.Println(bySalary)

    // 组合条件：IT 部门且薪资 > 45000
    itHighPay := df.Filter(func(r dataframe.Row) bool {
        dept := r.Get("dept").(string)
        salary := r.Get("salary").(float64)
        return dept == "IT" && salary > 45000
    })
    fmt.Println("\nIT 部门高薪员工:")
    fmt.Println(itHighPay)
}
```

## 2. 分组聚合分析

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 销售数据
    df, _ := dataframe.New(map[string][]interface{}{
        "date":     {"2024-01", "2024-01", "2024-02", "2024-02", "2024-03"},
        "product":  {"A", "B", "A", "B", "A"},
        "region":   {"East", "East", "West", "West", "East"},
        "sales":    {1000.0, 1500.0, 1200.0, 1800.0, 1100.0},
        "quantity": {10, 15, 12, 18, 11},
    })

    // 按产品分组统计
    gbProduct, _ := df.GroupBy("product")
    fmt.Println("=== 按产品统计 ===")
    fmt.Println("销售总额:")
    fmt.Println(gbProduct.Sum("sales"))
    fmt.Println("\n平均销量:")
    fmt.Println(gbProduct.Mean("quantity"))

    // 多列分组
    gbMulti, _ := df.GroupBy("product", "region")
    fmt.Println("\n=== 按产品和地区统计 ===")
    fmt.Println(gbMulti.Sum("sales"))

    // 多指标聚合
    aggFuncs := map[string][]dataframe.AggFunc{
        "sales":    {dataframe.AggSum, dataframe.AggMean, dataframe.AggMax},
        "quantity": {dataframe.AggSum, dataframe.AggCount},
    }
    result, _ := gbProduct.Agg(aggFuncs)
    fmt.Println("\n=== 多指标聚合 ===")
    fmt.Println(result)
}
```

## 3. 多表合并

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 订单表
    orders, _ := dataframe.New(map[string][]interface{}{
        "order_id":    {1, 2, 3, 4, 5},
        "customer_id": {101, 102, 101, 103, 102},
        "amount":      {250.0, 180.0, 320.0, 150.0, 420.0},
    })

    // 客户表
    customers, _ := dataframe.New(map[string][]interface{}{
        "customer_id": {101, 102, 104},
        "name":        {"Alice", "Bob", "Charlie"},
        "city":        {"Beijing", "Shanghai", "Guangzhou"},
    })

    // 产品表（使用不同的键名）
    products, _ := dataframe.New(map[string][]interface{}{
        "prod_id": {1, 2, 3},
        "name":    {"Laptop", "Phone", "Tablet"},
    })

    // Left Join: 所有订单 + 客户信息
    orderCustomer, _ := dataframe.Merge(orders, customers, dataframe.MergeOptions{
        How: dataframe.LeftJoin,
        On:  []string{"customer_id"},
    })
    fmt.Println("=== 订单详情 (Left Join) ===")
    fmt.Println(orderCustomer)

    // Inner Join: 只显示有客户信息的订单
    validOrders, _ := dataframe.Merge(orders, customers, dataframe.MergeOptions{
        How: dataframe.InnerJoin,
        On:  []string{"customer_id"},
    })
    fmt.Println("\n=== 有效订单 (Inner Join) ===")
    fmt.Println(validOrders)

    // 使用不同列名合并
    orderProducts, _ := dataframe.New(map[string][]interface{}{
        "order_id":   {1, 2, 3},
        "product_id": {1, 2, 1},
    })
    withProducts, _ := orderProducts.MergeOn(products,
        []string{"product_id"},
        []string{"prod_id"},
        dataframe.LeftJoin,
    )
    fmt.Println("\n=== 订单产品 (不同键名) ===")
    fmt.Println(withProducts)
}
```

## 4. 统计分析

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 考试成绩
    scores := dataframe.NewSeriesFromFloat64s(
        []float64{85, 92, 78, 95, 88, 76, 91, 83, 79, 94},
        "scores",
    )

    fmt.Println("=== 成绩统计 ===")
    fmt.Printf("总分: %.1f\n", scores.Sum())
    fmt.Printf("平均: %.2f\n", scores.Mean())
    fmt.Printf("中位数: %.1f\n", scores.Median())
    fmt.Printf("标准差: %.2f\n", scores.Std())
    fmt.Printf("最低: %v\n", scores.Min())
    fmt.Printf("最高: %v\n", scores.Max())
    fmt.Printf("人数: %d\n", scores.Count())

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
    
    fmt.Println("\n=== 等级分布 ===")
    fmt.Println(grades.ValueCounts())

    // 标准化分数 (Z-Score)
    mean := scores.Mean()
    std := scores.Std()
    zScores := scores.Sub(mean).Div(std)
    fmt.Println("\n=== Z-Score ===")
    fmt.Println(zScores.Head(5))
}
```

## 5. 缺失值处理

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 带缺失值的数据
    df, _ := dataframe.New(map[string][]interface{}{
        "name":   {"Alice", "Bob", nil, "David", "Eve"},
        "age":    {25, nil, 30, 35, nil},
        "salary": {50000.0, 60000.0, nil, 70000.0, 55000.0},
    })

    fmt.Println("原始数据:")
    fmt.Println(df)

    // 检查缺失值
    ageSeries, _ := df.GetSeries("age")
    fmt.Println("\nage 列缺失值标记:")
    fmt.Println(ageSeries.IsNA())

    // 填充缺失值
    filledAge := ageSeries.FillNA(0)
    fmt.Println("\n填充后的 age:")
    fmt.Println(filledAge)

    // 用均值填充
    meanAge := ageSeries.Mean()
    filledWithMean := ageSeries.FillNA(meanAge)
    fmt.Println("\n用均值填充的 age:")
    fmt.Println(filledWithMean)

    // 删除缺失值
    cleanAge := ageSeries.DropNA()
    fmt.Println("\n删除缺失值后:")
    fmt.Println(cleanAge)
}
```

## 6. 并行处理大数据

```go
package main

import (
    "fmt"
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

    // 复杂计算函数
    compute := func(v interface{}) interface{} {
        if f, ok := v.(float64); ok {
            return f*f + f*2 + 1
        }
        return v
    }

    // 普通处理
    start := time.Now()
    _ = s.Apply(compute)
    normalTime := time.Since(start)
    fmt.Printf("普通 Apply: %v\n", normalTime)

    // 并行处理
    start = time.Now()
    _ = s.ParallelApply(compute)
    parallelTime := time.Since(start)
    fmt.Printf("并行 Apply: %v\n", parallelTime)

    // 加速比
    fmt.Printf("加速比: %.2fx\n", float64(normalTime)/float64(parallelTime))

    // DataFrame 并行操作
    dfData := map[string][]interface{}{
        "a": make([]interface{}, 100000),
        "b": make([]interface{}, 100000),
    }
    for i := 0; i < 100000; i++ {
        dfData["a"][i] = float64(i)
        dfData["b"][i] = float64(i * 2)
    }
    df, _ := dataframe.New(dfData)

    // 并行聚合
    start = time.Now()
    sums := df.ParallelSum()
    fmt.Printf("\n并行求和耗时: %v\n", time.Since(start))
    fmt.Printf("a 列总和: %.0f\n", sums["a"])
}
```

## 7. 文件读写

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
    "github.com/datago/io"
)

func main() {
    // === 从 Excel 读取 ===
    excelDF, err := io.ReadExcel("sales.xlsx", io.ExcelOptions{
        Sheet:     "Q1",
        HasHeader: true,
        UseCols:   []string{"product", "amount", "date"},
    })
    if err != nil {
        fmt.Printf("读取 Excel 失败: %v\n", err)
    } else {
        fmt.Println("Excel 数据:")
        fmt.Println(excelDF.Head(3))
    }

    // === 从 CSV 读取 ===
    csvDF, err := io.ReadCSV("customers.csv", io.CSVOptions{
        Separator: ',',
        HasHeader: true,
        DTypes: map[string]dataframe.DType{
            "age": dataframe.DTypeInt64,
        },
    })
    if err != nil {
        fmt.Printf("读取 CSV 失败: %v\n", err)
    } else {
        fmt.Println("\nCSV 数据:")
        fmt.Println(csvDF.Head(3))
    }

    // === 创建并写入数据 ===
    df, _ := dataframe.New(map[string][]interface{}{
        "name":    {"产品A", "产品B", "产品C"},
        "price":   {99.9, 199.9, 299.9},
        "stock":   {100, 50, 30},
    })

    // 写入 Excel
    err = io.WriteExcel("products.xlsx", df, io.ExcelWriteOptions{
        Sheet:        "产品列表",
        IncludeIndex: false,
    })
    if err == nil {
        fmt.Println("\n已写入 products.xlsx")
    }

    // 写入 CSV
    err = io.WriteCSV("products.csv", df, io.CSVWriteOptions{
        Separator:    ',',
        IncludeIndex: true,
        IndexName:    "id",
    })
    if err == nil {
        fmt.Println("已写入 products.csv")
    }
}
```

## 8. 数据透视分析

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // 电商订单数据
    orders, _ := dataframe.New(map[string][]interface{}{
        "date":     {"2024-01-01", "2024-01-01", "2024-01-02", "2024-01-02", "2024-01-03"},
        "category": {"电子", "服装", "电子", "食品", "服装"},
        "region":   {"华东", "华北", "华东", "华南", "华东"},
        "amount":   {1500.0, 800.0, 2200.0, 500.0, 1200.0},
    })

    fmt.Println("=== 订单数据 ===")
    fmt.Println(orders)

    // 按类别统计
    gbCategory, _ := orders.GroupBy("category")
    fmt.Println("\n=== 按类别统计 ===")
    fmt.Println("销售总额:")
    fmt.Println(gbCategory.Sum("amount"))

    // 按地区统计
    gbRegion, _ := orders.GroupBy("region")
    fmt.Println("\n=== 按地区统计 ===")
    fmt.Println("销售总额:")
    fmt.Println(gbRegion.Sum("amount"))
    fmt.Println("平均订单金额:")
    fmt.Println(gbRegion.Mean("amount"))

    // 按类别和地区交叉统计
    gbCross, _ := orders.GroupBy("category", "region")
    fmt.Println("\n=== 交叉统计 ===")
    fmt.Println(gbCross.Sum("amount"))

    // 找出各类别销售最高的地区
    topByCategory := gbCross.Apply(func(g *dataframe.DataFrame) *dataframe.DataFrame {
        return g.SortBy("amount", dataframe.Descending).Head(1)
    })
    fmt.Println("\n=== 各类别销售最高的地区 ===")
    fmt.Println(topByCategory)
}
```

## 相关章节

- [DataFrame 使用指南](./dataframe)
- [Series 使用指南](./series)
- [GroupBy 分组聚合](./groupby)
- [Merge/Join 数据合并](./merge)
- [并行处理](./parallel)
- [Excel 读写](./io-excel)
- [CSV 读写](./io-csv)
