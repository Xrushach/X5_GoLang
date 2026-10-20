package main

import "fmt"

type Item struct {
	Name         string
	PriceKopecks int64
}

func main() {
	var empty Item
	fmt.Printf("нулевое значение: %+v\n", empty)

	item := Item{Name: "Чай", PriceKopecks: 150_00}
	fmt.Println("поля:", item.Name, item.PriceKopecks)

	another := item
	another.Name = "Кофе"
	another.PriceKopecks = 200_00
	fmt.Printf("исходный: %+v\nкопия: %+v\n", item, another)
}
