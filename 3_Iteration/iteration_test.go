package main

import "testing"

func TestRepeat(t *testing.T) {
	got := Repeat("a", -2)
	expected := ""

	if got != expected {
		t.Errorf("Got: %q , expected: %q. ", got, expected)
	}
}

/*
	Benchmarking:
	=>	You'll see the code is very similar to a test. The testing.B gives you access to the loop function. Loop() returns true
		as long as the benchmark should continue running. When the benchmark code is executed, it measures how long it takes.
		After Loop() returns false, b.N contains the total number of iterations that ran. The number of times the code is run
		shouldn't matter to you, the framework will determine what is a "good" value for that to let you have some decent results.
	=>	To run the benchmarks do go test -bench=.
		Gives 2 values which are the number of loops run and the time per operation carried out.
	=>	Structure of benchmark:
			func Benchmark(b *testing.B) {
			//... setup ...
			for b.Loop() {
				//... code to measure ...
			}
			//... cleanup ...
			}
		Only the body of the loop is timed; it automatically excludes setup and cleanup code from benchmark timing.
	=>	go test -bench=. -benchmem gives:
			B/op: the number of bytes allocated per iteration
			allocs/op: the number of memory allocations per iteration
*/

func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		Repeat("a", 5)
	}
}
