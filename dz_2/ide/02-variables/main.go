package main

import "fmt"

func main() {
	var count int
	var price float64
	var name string
	var paid bool
	fmt.Printf("нулевые: %d %.1f %q %t\n", count, price, name, paid)
	count = 3
	customer := "Анна"
	customer = "Иван"
	const currency = "RUB"
	const (
		StatusNew = iota << 2
		StatusPaid
		StatusCancelled
	)
	fmt.Println(count, customer, currency)
	fmt.Println(StatusNew, StatusPaid, StatusCancelled)
	customer = "Ольга"
	count = 123
	{
		count := 10 // Новая переменная внутреннего блока.
		fmt.Println("внутри:", count)
	}
	fmt.Println("снаружи:", count)
}
