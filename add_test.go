package main

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		a        int
		b        int
		expected int
	}{
		{0, 0, 0},
		{0, 1, 1},
		{0, -1, -1},
		{0, math.MaxInt, 2147483647},
		{0, math.MinInt, -2147483648},
		{1, 0, 1},
		{1, 1, 2},
		{1, -1, 0},
		{1, math.MaxInt, -2147483648},
		{1, math.MinInt, -2147483647},
		{-1, 0, -1},
		{-1, 1, 0},
		{-1, -1, -2},
		{-1, math.MaxInt, 2147483646},
		{-1, math.MinInt, 2147483647},
		{math.MaxInt, 0, 2147483647},
		{math.MaxInt, 1, -2147483648},
		{math.MaxInt, -1, 2147483646},
		{math.MaxInt, math.MaxInt, -2},
		{math.MaxInt, math.MinInt, -1},
		{math.MinInt, 0, -2147483648},
		{math.MinInt, 1, -2147483647},
		{math.MinInt, -1, 2147483647},
		{math.MinInt, math.MaxInt, -1},
		{math.MinInt, math.MinInt, 0},
	}

	for _, tc := range testCases {
		result := Add(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Add expected %v, got %v", tc.expected, result)
		}
	}
}

func TestSubtract(t *testing.T) {
	testCases := []struct {
		a        int
		b        int
		expected int
	}{
		{0, 0, 0},
		{0, 1, -1},
		{0, -1, 1},
		{0, math.MaxInt, -2147483647},
		{0, math.MinInt, -2147483648},
		{1, 0, 1},
		{1, 1, 0},
		{1, -1, 2},
		{1, math.MaxInt, -2147483646},
		{1, math.MinInt, -2147483647},
		{-1, 0, -1},
		{-1, 1, -2},
		{-1, -1, 0},
		{-1, math.MaxInt, -2147483648},
		{-1, math.MinInt, 2147483647},
		{math.MaxInt, 0, 2147483647},
		{math.MaxInt, 1, 2147483646},
		{math.MaxInt, -1, -2147483648},
		{math.MaxInt, math.MaxInt, 0},
		{math.MaxInt, math.MinInt, -1},
		{math.MinInt, 0, -2147483648},
		{math.MinInt, 1, 2147483647},
		{math.MinInt, -1, -2147483647},
		{math.MinInt, math.MaxInt, 1},
		{math.MinInt, math.MinInt, 0},
	}

	for _, tc := range testCases {
		result := Subtract(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Subtract expected %v, got %v", tc.expected, result)
		}
	}
}

func TestMultiply(t *testing.T) {
	testCases := []struct {
		a        int
		b        int
		expected int
	}{
		{0, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, math.MaxInt, 0},
		{0, math.MinInt, 0},
		{1, 0, 0},
		{1, 1, 1},
		{1, -1, -1},
		{1, math.MaxInt, 2147483647},
		{1, math.MinInt, -2147483648},
		{-1, 0, 0},
		{-1, 1, -1},
		{-1, -1, 1},
		{-1, math.MaxInt, -2147483647},
		{-1, math.MinInt, -2147483648},
		{math.MaxInt, 0, 0},
		{math.MaxInt, 1, 2147483647},
		{math.MaxInt, -1, -2147483647},
		{math.MaxInt, math.MaxInt, 1},
		{math.MaxInt, math.MinInt, -2147483648},
		{math.MinInt, 0, 0},
		{math.MinInt, 1, -2147483648},
		{math.MinInt, -1, -2147483648},
		{math.MinInt, math.MaxInt, -2147483648},
		{math.MinInt, math.MinInt, 0},
	}

	for _, tc := range testCases {
		result := Multiply(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Multiply expected %v, got %v", tc.expected, result)
		}
	}
}

func TestDivide(t *testing.T) {
	testCases := []struct {
		a        int
		b        int
		expected int
	}{
		{0, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, math.MaxInt, 0},
		{0, math.MinInt, 0},
		{1, 0, 0},
		{1, 1, 1},
		{1, -1, -1},
		{1, math.MaxInt, 0},
		{1, math.MinInt, 0},
		{-1, 0, 0},
		{-1, 1, -1},
		{-1, -1, 1},
		{-1, math.MaxInt, 0},
		{-1, math.MinInt, 0},
		{math.MaxInt, 0, 0},
		{math.MaxInt, 1, 2147483647},
		{math.MaxInt, -1, -2147483647},
		{math.MaxInt, math.MaxInt, 1},
		{math.MaxInt, math.MinInt, 0},
		{math.MinInt, 0, 0},
		{math.MinInt, 1, -2147483648},
		{math.MinInt, -1, -2147483648},
		{math.MinInt, math.MaxInt, -1},
		{math.MinInt, math.MinInt, 1},
	}

	for _, tc := range testCases {
		result := Divide(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Divide expected %v, got %v", tc.expected, result)
		}
	}
}
