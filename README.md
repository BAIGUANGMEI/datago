# Benchmark

## ReadExcel

The following table summarizes the benchmark results for reading Excel files using different programming languages and libraries. Each test was conducted over 5 rounds, and the average read time in seconds is reported.

### testdata.xlsx (15,430 rows × 11 columns)

| Language | Package | Rounds | Avg Read Time (s) |
|:--------|:--------|-------:|------------------:|
| Go      | datago  | 5      | 0.2120            |
| Python  | pandas  | 5      | 0.5135            |
| Python  | polars  | 5      | 0.1020            |


### testdatalarge.xlsx (271,114 rows × 16 columns)

| Language | Package | Rounds | Avg Read Time (s) |
|:--------|:--------|-------:|------------------:|
| Go      | datago  | 5      | 5.7648            |
| Python  | pandas  | 5      | 11.0093           |
| Python  | polars  | 5      | 2.1800            |
