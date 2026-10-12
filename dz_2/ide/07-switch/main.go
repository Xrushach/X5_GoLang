package main

import "fmt"

func main() {
	status := "paid"
	switch status {
	case "new":
		fmt.Println("ждём оплаты")
	case "paid", "shipped":
		fmt.Println("в работе")
	default:
		fmt.Println("неизвестный статус")
	}
	amount := int64(50_000)
	switch {
	case amount < 0:
		fmt.Println("invalid")

	case amount < 10_000:
		fmt.Println("small")
	case amount < 50_000:
		fmt.Println("regular")
	default:
		fmt.Println("large")
	}
	for n := 1; n <= 30; n++ {
		switch {
		case n%15 == 0:
			fmt.Println("FizzBuzz")
		case n%3 == 0:
			fmt.Println("Fizz")
		case n%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(n)
		}
	}
}
