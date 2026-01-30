package tests

import (
	"math"
	"testing"

	"github.com/datago/dataframe"
)

func TestSeriesBasicStats(t *testing.T) {
	s := dataframe.NewSeries([]interface{}{1, 2, 3, 4, 5}, "nums")

	if got := s.Sum(); got != 15 {
		t.Fatalf("Sum() = %v, want 15", got)
	}
	if got := s.Mean(); got != 3 {
		t.Fatalf("Mean() = %v, want 3", got)
	}
	if got := s.Median(); got != 3 {
		t.Fatalf("Median() = %v, want 3", got)
	}
	if got := s.Min(); got != float64(1) {
		t.Fatalf("Min() = %v, want 1", got)
	}
	if got := s.Max(); got != float64(5) {
		t.Fatalf("Max() = %v, want 5", got)
	}
}

func TestSeriesStdVar(t *testing.T) {
	s := dataframe.NewSeries([]interface{}{1, 2, 3, 4, 5}, "nums")
	if got := s.Var(); math.Abs(got-2.5) > 1e-9 {
		t.Fatalf("Var() = %v, want 2.5", got)
	}
	if got := s.Std(); math.Abs(got-math.Sqrt(2.5)) > 1e-9 {
		t.Fatalf("Std() = %v, want sqrt(2.5)", got)
	}
}

func TestSeriesFillDropNA(t *testing.T) {
	s := dataframe.NewSeries([]interface{}{1, nil, 3, "NA"}, "nums")
	filled := s.FillNA(0)
	if v, _ := filled.Get(1); v != 0 {
		t.Fatalf("FillNA() at 1 = %v, want 0", v)
	}
	dropped := s.DropNA()
	if dropped.Len() != 2 {
		t.Fatalf("DropNA() len = %d, want 2", dropped.Len())
	}
}

func TestSeriesArithmetic(t *testing.T) {
	s := dataframe.NewSeries([]interface{}{1, 2, 3}, "nums")
	add := s.Add(1)
	if v, _ := add.Get(0); v != float64(2) {
		t.Fatalf("Add(1) first = %v, want 2", v)
	}
	mul := s.Mul(2)
	if v, _ := mul.Get(2); v != float64(6) {
		t.Fatalf("Mul(2) third = %v, want 6", v)
	}
}
