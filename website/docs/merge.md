---
sidebar_position: 6
title: Merge/Join 数据合并
---

# Merge/Join 数据合并

Merge 和 Join 功能允许你基于一个或多个键将两个 DataFrame 合并，类似于 SQL 的 JOIN 操作。

## 合并类型

| 类型 | SQL 等价 | 说明 |
|------|----------|------|
| `InnerJoin` | `INNER JOIN` | 只返回两表都有的键 |
| `LeftJoin` | `LEFT JOIN` | 保留左表所有行 |
| `RightJoin` | `RIGHT JOIN` | 保留右表所有行 |
| `OuterJoin` | `FULL OUTER JOIN` | 保留两表所有行 |

```
左表: id=[1,2,3]    右表: id=[2,3,4]

InnerJoin → id=[2,3]      (交集)
LeftJoin  → id=[1,2,3]    (左表全部)
RightJoin → id=[2,3,4]    (右表全部)
OuterJoin → id=[1,2,3,4]  (并集)
```

## 基本用法

### 使用 Merge 函数

```go
import "github.com/datago/dataframe"

// 员工表
employees, _ := dataframe.New(map[string][]interface{}{
    "emp_id": {1, 2, 3, 4},
    "name":   {"Alice", "Bob", "Charlie", "David"},
})

// 薪资表
salaries, _ := dataframe.New(map[string][]interface{}{
    "emp_id": {2, 3, 5},
    "salary": {50000, 60000, 70000},
})

// Inner Join - 只保留两表都有的员工
result, err := dataframe.Merge(employees, salaries, dataframe.MergeOptions{
    How: dataframe.InnerJoin,
    On:  []string{"emp_id"},
})
// emp_id  name     salary
// 2       Bob      50000
// 3       Charlie  60000
```

### 使用 DataFrame 方法

```go
// Join 方法（简洁语法）
result, err := employees.Join(salaries, []string{"emp_id"}, dataframe.LeftJoin)

// MergeOn 方法（支持不同列名）
result, err := employees.MergeOn(salaries,
    []string{"emp_id"},      // 左表列
    []string{"employee_id"}, // 右表列
    dataframe.LeftJoin,
)
```

## 合并选项

### MergeOptions 详解

```go
type MergeOptions struct {
    How         JoinType   // 合并类型：InnerJoin/LeftJoin/RightJoin/OuterJoin
    On          []string   // 两表共同的键列名
    LeftOn      []string   // 左表的键列（与 RightOn 配合使用）
    RightOn     []string   // 右表的键列
    Suffixes    [2]string  // 重名列的后缀，默认 ["_x", "_y"]
    Indicator   bool       // 是否添加 _merge 列显示来源
}
```

### 处理重复列名

当两表有同名非键列时，自动添加后缀区分：

```go
left, _ := dataframe.New(map[string][]interface{}{
    "id":    {1, 2},
    "value": {100, 200},
})

right, _ := dataframe.New(map[string][]interface{}{
    "id":    {1, 2},
    "value": {10, 20},
})

result, _ := dataframe.Merge(left, right, dataframe.MergeOptions{
    How:      dataframe.InnerJoin,
    On:       []string{"id"},
    Suffixes: [2]string{"_left", "_right"},
})
// 结果列: id, value_left, value_right
```

### 添加来源指示列

```go
result, _ := dataframe.Merge(left, right, dataframe.MergeOptions{
    How:       dataframe.OuterJoin,
    On:        []string{"id"},
    Indicator: true,
})
// 添加 _merge 列，值为: "left_only", "right_only", "both"
```

## 各种 Join 示例

### Inner Join

```go
// 只保留两表都有的记录
result, _ := dataframe.Merge(left, right, dataframe.MergeOptions{
    How: dataframe.InnerJoin,
    On:  []string{"id"},
})
```

### Left Join

```go
// 保留左表所有记录，右表无匹配则填充 nil
result, _ := dataframe.Merge(employees, departments, dataframe.MergeOptions{
    How: dataframe.LeftJoin,
    On:  []string{"dept_id"},
})
```

### Right Join

```go
// 保留右表所有记录
result, _ := dataframe.Merge(orders, products, dataframe.MergeOptions{
    How: dataframe.RightJoin,
    On:  []string{"product_id"},
})
```

### Outer Join

```go
// 保留两表所有记录
result, _ := dataframe.Merge(table1, table2, dataframe.MergeOptions{
    How: dataframe.OuterJoin,
    On:  []string{"key"},
})
```

## 多键合并

支持基于多列进行合并：

```go
// 销售数据
sales, _ := dataframe.New(map[string][]interface{}{
    "year":    {2023, 2023, 2024, 2024},
    "quarter": {1, 2, 1, 2},
    "sales":   {1000, 1200, 1100, 1300},
})

// 目标数据
targets, _ := dataframe.New(map[string][]interface{}{
    "year":    {2023, 2024},
    "quarter": {1, 1},
    "target":  {950, 1050},
})

// 按年份和季度合并
result, _ := dataframe.Merge(sales, targets, dataframe.MergeOptions{
    How: dataframe.LeftJoin,
    On:  []string{"year", "quarter"},
})
```

## 不同列名合并

当两表的键列名不同时：

```go
employees, _ := dataframe.New(map[string][]interface{}{
    "emp_id": {1, 2, 3},
    "name":   {"Alice", "Bob", "Charlie"},
})

reviews, _ := dataframe.New(map[string][]interface{}{
    "employee_id": {1, 2},
    "rating":      {4.5, 4.8},
})

result, _ := dataframe.Merge(employees, reviews, dataframe.MergeOptions{
    How:     dataframe.LeftJoin,
    LeftOn:  []string{"emp_id"},
    RightOn: []string{"employee_id"},
})
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/datago/dataframe"
)

func main() {
    // === 示例：员工数据库 ===
    
    // 员工基本信息
    employees, _ := dataframe.New(map[string][]interface{}{
        "emp_id":  {1, 2, 3, 4, 5},
        "name":    {"Alice", "Bob", "Charlie", "David", "Eve"},
        "dept_id": {10, 20, 10, 30, 20},
    })

    // 部门信息
    departments, _ := dataframe.New(map[string][]interface{}{
        "dept_id":   {10, 20, 40},
        "dept_name": {"工程部", "市场部", "人事部"},
    })

    // 薪资信息
    salaries, _ := dataframe.New(map[string][]interface{}{
        "emp_id": {1, 2, 3, 6},
        "salary": {80000, 75000, 85000, 90000},
    })

    fmt.Println("=== 员工表 ===")
    fmt.Println(employees)

    fmt.Println("\n=== 部门表 ===")
    fmt.Println(departments)

    fmt.Println("\n=== 薪资表 ===")
    fmt.Println(salaries)

    // 1. Left Join: 员工 + 部门
    empDept, _ := dataframe.Merge(employees, departments, dataframe.MergeOptions{
        How: dataframe.LeftJoin,
        On:  []string{"dept_id"},
    })
    fmt.Println("\n=== 员工部门信息 (Left Join) ===")
    fmt.Println(empDept)

    // 2. Inner Join: 员工 + 薪资（只显示有薪资记录的员工）
    empSalary, _ := dataframe.Merge(employees, salaries, dataframe.MergeOptions{
        How: dataframe.InnerJoin,
        On:  []string{"emp_id"},
    })
    fmt.Println("\n=== 有薪资记录的员工 (Inner Join) ===")
    fmt.Println(empSalary)

    // 3. Outer Join + Indicator: 查看完整匹配情况
    fullJoin, _ := dataframe.Merge(employees, salaries, dataframe.MergeOptions{
        How:       dataframe.OuterJoin,
        On:        []string{"emp_id"},
        Indicator: true,
    })
    fmt.Println("\n=== 完整匹配情况 (Outer Join + Indicator) ===")
    fmt.Println(fullJoin)

    // 4. 多表合并
    result, _ := dataframe.Merge(empDept, salaries, dataframe.MergeOptions{
        How: dataframe.LeftJoin,
        On:  []string{"emp_id"},
    })
    fmt.Println("\n=== 完整员工信息 ===")
    fmt.Println(result)
}
```

## 性能提示

1. **哈希索引**：内部使用哈希表加速查找，大数据量效率高
2. **内存管理**：Outer Join 可能产生大量数据，注意内存
3. **键选择**：使用高区分度的列作为键可提升效率
4. **数据预处理**：合并前清理重复数据可减少结果行数

## 相关章节

- [DataFrame 使用指南](./dataframe) - 了解基本数据操作
- [GroupBy 分组聚合](./groupby) - 合并后常需聚合分析
- [并行处理](./parallel) - 大数据量处理优化
