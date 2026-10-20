package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		fmt.Println("счётчик:", i)
	}
	left := 3
	for left > 0 {
		fmt.Println("осталось:", left)
		left--
	}
	attempts := 0

	for {
		attempts++
		if attempts == 2 {
			break
		}
	}
	fmt.Println("попыток:", attempts)
	values := []int64{12_000, -100, 7_500, 0, 31_000, 4_500}
	var total int64
	for index, value := range values {
		if value <= 0 {
			continue
		}
		fmt.Println("элемент:", index, value)
		total += value
	}
	fmt.Println("сумма положительных:", total)
	for _, value := range values {
		value = 0 // Меняем копию, а не элемент среза.
		fmt.Print(value, " ")
	}
	fmt.Println("\nисходные:", values)
	for index := range values {
		values[index] = 0
	}
	fmt.Println("после изменения по индексу:", values)
}
