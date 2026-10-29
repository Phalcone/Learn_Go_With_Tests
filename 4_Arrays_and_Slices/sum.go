package main

func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAll(numbers [][]int) (result []int) {
	result = make([]int, len(numbers))
	for i := range numbers {
		result[i] = Sum(numbers[i])
	}
	return result
}
