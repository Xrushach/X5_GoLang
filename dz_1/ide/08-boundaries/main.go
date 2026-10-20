package main

import "fmt"

func main() {
	values := []int64{12_000, 7_500, 31_000, 4_500}
	// values = []int64{} // Включить для проверки пустого ввода.
	if len(values) == 0 {
		fmt.Println("нет данных для среднего")
		return
	}
	minimum := values[0]
	maximum := values[0]
	var total int64
	for _, value := range values {
		if value < minimum {
			minimum = value
		}
		if value > maximum {
			maximum = value
		}
		total += value
	}
	average := float64(total) / float64(len(values))
	fmt.Println("min, max, sum, avg:", minimum, maximum, total, average)
}
