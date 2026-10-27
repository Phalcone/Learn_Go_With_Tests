package main

import "testing"

func TestAdd(t *testing.T) {
	sum := Add(2, 3)
	expected_sum := 5

	if sum != expected_sum {
		t.Errorf("Got: %d but expected: %d.", sum, expected_sum)
	}
}
