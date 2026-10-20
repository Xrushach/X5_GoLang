package main

import "fmt"

type Item struct{ PriceKopecks int64 }

func main() {
	items := []Item{{PriceKopecks: 100}, {PriceKopecks: 200}}
	for _, item := range items {
		item.PriceKopecks += 50
	}
	fmt.Println("меняли переменную range:", items)

	for i := range items {
		items[i].PriceKopecks += 50
	}
	fmt.Println("меняли элементы по индексу:", items)
}
