package main

import (
	"slices"
	"testing"
)

func TestSum(t *testing.T) {
	numbers := []int{1, 2, 3}

	got := Sum(numbers)
	want := 6

	if got != want {
		t.Errorf("got %d want %d given, %v", got, want, numbers)
	}
}

func TestSumAll(t *testing.T) {
	numbers := [][]int{{1, 2, 3},
		{5, 3},
		{1}}

	got := SumAll(numbers)
	want := []int{6, 8, 1}

	if !slices.Equal(got, want) {
		t.Errorf("Got: %v, Want: %v", got, want)
	}
}
