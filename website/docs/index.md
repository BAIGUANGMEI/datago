---
sidebar_position: 5
title: Index 使用指南
---

`Index` 是 `DataFrame` 与 `Series` 的行索引系统。

## 创建与基础操作

- `NewIndex(labels, name)`：用自定义标签创建
- `NewRangeIndex(size)`：默认整数索引
- `Len()` / `Name()` / `SetName(name)`
- `Labels()` / `Get(pos)` / `GetLoc(label)` / `Contains(label)`

## 切片与复制

- `Slice(start, end)`
- `Append(label)`
- `Copy()` / `Reset()`
- `ToStringSlice()`

## 集合操作

- `Equals(other)`
- `Union(other)`
- `Intersection(other)`
- `Difference(other)`
