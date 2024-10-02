package main

import (
	"testing"
	"math"
)

func TestBinarySearch(t *testing.T) {
	testCases := []struct {
		arr []int
		target int
		expected int
	}{
		{nil, 0, 0},
		{nil, 1, 0},
		{nil, -1, 0},
		{nil, math.MaxInt, 0},
		{nil, math.MinInt, 0},
		{[]int{}, 0, 0},
		{[]int{}, 1, 0},
		{[]int{}, -1, 0},
		{[]int{}, math.MaxInt, 0},
		{[]int{}, math.MinInt, 0},
		{[]int{0}, 0, 0},
		{[]int{0}, 1, 0},
		{[]int{0}, -1, 0},
		{[]int{0}, math.MaxInt, 0},
		{[]int{0}, math.MinInt, 0},
	}

	for _, tc := range testCases {
		result := BinarySearch(tc.arr, tc.target)
		if result != tc.expected {
			t.Errorf("BinarySearch expected %v, got %v", tc.expected, result)
		}
	}
}

