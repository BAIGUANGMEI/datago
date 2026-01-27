"""
Benchmark: pandas vs polars Excel reading performance
"""

import time
import os

# Get the path to testdata.xlsx (in parent directory: tests/)
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PARENT_DIR = os.path.dirname(SCRIPT_DIR)

EXCEL_FILE = os.path.join(PARENT_DIR, "testdatalarge.xlsx")


def benchmark_pandas():
    import pandas as pd
    
    start = time.perf_counter()
    df = pd.read_excel(EXCEL_FILE)
    elapsed = time.perf_counter() - start
    
    print(f"pandas: {elapsed:.4f}s, shape: {df.shape}")
    return elapsed


def benchmark_polars():
    import polars as pl
    
    start = time.perf_counter()
    df = pl.read_excel(EXCEL_FILE)
    elapsed = time.perf_counter() - start
    
    print(f"polars: {elapsed:.4f}s, shape: {df.shape}")
    return elapsed


def run_benchmark(iterations=5):
    print(f"Benchmarking Excel read: {EXCEL_FILE}")
    print(f"Iterations: {iterations}")
    print("-" * 40)
    
    pandas_times = []
    polars_times = []
    
    for i in range(iterations):
        print(f"\n--- Run {i + 1} ---")
        pandas_times.append(benchmark_pandas())
        polars_times.append(benchmark_polars())
    
    print("\n" + "=" * 40)
    print("Results:")
    print(f"pandas avg: {sum(pandas_times) / len(pandas_times):.4f}s")
    print(f"polars avg: {sum(polars_times) / len(polars_times):.4f}s")


if __name__ == "__main__":
    if not os.path.exists(EXCEL_FILE):
        print(f"Error: {EXCEL_FILE} not found")
        print("Please create testdata.xlsx in the benchmark folder")
    else:
        run_benchmark()
