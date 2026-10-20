package main

import "fmt"

type Item struct {
	Name         string
	PriceKopecks int64
}

type Order struct {
	ID     int
	Items  []Item
	Status string
}

// Короткий мостик к методам, только после разбора полей и литерала.
func (order Order) Label() string { return fmt.Sprintf("заказ %d", order.ID) }

func main() {
	var empty Order
	fmt.Printf("нулевое значение: %+v; Items == nil: %t\n", empty, empty.Items == nil)

	original := Order{ID: 42, Status: "new", Items: []Item{{Name: "Чай", PriceKopecks: 150_00}}}
	fmt.Println("поле структуры, элемент слайса, поле позиции:", original.Items[0].Name)

	another := original
	another.Status = "paid"
	another.Items[0].Name = "Кофе"
	fmt.Printf("исходный: %+v\nкопия: %+v\n", original, another)
	// Status независим, Items описывают общий массив. Нового массива при := нет.

	independent := original
	independent.Items = make([]Item, len(original.Items))
	copy(independent.Items, original.Items)
	independent.Items[0].Name = "Сок"
	fmt.Println("копия позиций:", original.Items, independent.Items)
	// Для этой модели Item такого копирования достаточно, это не универсальная
	// глубокая копия произвольных вложенных данных.
	fmt.Println("метод:", original.Label())
}
