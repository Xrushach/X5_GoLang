package main

import "fmt"

func main() {
	var empty [3]int
	fmt.Println("нулевой массив:", empty, "длина:", len(empty))

	numbers := [3]int{10, 20, 30}
	fmt.Println("по индексу:", numbers[0], numbers[2])
	for i, value := range numbers {
		fmt.Println("индекс:", i, "значение:", value)
	}

	another := numbers
	another[0] = 99
	fmt.Println("исходный:", numbers, "копия:", another)

	longer := [4]int{10, 20, 30, 40}
	fmt.Printf("длина входит в тип: %T и %T\n", numbers, longer)
	// numbers = longer // Не скомпилируется: [3]int и [4]int являются разными типами.
}
