package main

import "fmt"

// Необязательная часть: variadic, именованный результат, функция как значение.
func sumPrices(prices ...int64) (total int64) {
	for _, price := range prices {
		total += price
	}
	return total
}

func double(value int64) int64 { return value * 2 }

func main() {
	fmt.Println("variadic:", sumPrices(100, 200, 300))
	fmt.Println("без аргументов:", sumPrices())
	operation := double
	fmt.Println("функция в переменной:", operation(7))
}
