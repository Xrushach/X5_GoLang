package main

import "fmt"

func main() {
	amount := int64(12_000)
	paid := true
	if amount < 0 {
		fmt.Println("некорректная сумма")
	} else if amount == 0 {
		fmt.Println("нулевая сумма")
	} else {
		fmt.Println("положительная сумма")
	}
	if paid && amount > 0 {
		fmt.Println("можно собирать заказ")
	}
	if discount := amount / 10; discount >= 1_000 {
		fmt.Println("скидка:", discount)
	} else {
		fmt.Println("скидка меньше 1000:", discount)
	}
	// fmt.Println(discount) // Переменная видна только внутри if и else.
}
