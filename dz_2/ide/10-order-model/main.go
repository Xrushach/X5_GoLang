package main

import (
	"errors"
	"fmt"
)

type Item struct {
	Name         string
	PriceKopecks int64
}

type Order struct {
	ID     int
	Items  []Item
	Status string
}

// NewOrder создаёт заказ с положительным ID и статусом new.
func NewOrder(id int) (Order, error) {
	if id <= 0 {
		return Order{}, errors.New("id должен быть положительным")
	}
	return Order{ID: id, Status: "new"}, nil
}

// AddItem проверяет вход до изменения заказа. Нулевая цена допустима.
func AddItem(order *Order, item Item) error {
	if order == nil {
		return errors.New("заказ не задан")
	}
	if item.Name == "" {
		return errors.New("название не задано")
	}
	if item.PriceKopecks < 0 {
		return errors.New("цена не может быть отрицательной")
	}
	order.Items = append(order.Items, item)
	return nil
}

// Total не изменяет заказ. Сумма по условию помещается в int64.
func Total(order Order) int64 {
	var total int64
	for _, item := range order.Items {
		total += item.PriceKopecks
	}
	return total
}

// StatusName возвращает название и признак наличия кода в справочнике.
func StatusName(code string) (string, bool) {
	names := map[string]string{"new": "новый", "paid": "оплачен", "cancelled": "отменён"}
	name, ok := names[code]
	return name, ok
}

func main() {
	order, err := NewOrder(42)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := AddItem(&order, Item{Name: "Чай", PriceKopecks: 150_00}); err != nil {
		fmt.Println(err)
		return
	}
	if err := AddItem(&order, Item{Name: "Кофе", PriceKopecks: 200_00}); err != nil {
		fmt.Println(err)
		return
	}
	name, ok := StatusName(order.Status)
	fmt.Println("заказ:", order.ID, "статус:", name, "найден:", ok, "сумма:", Total(order))
	fmt.Println("nil-заказ:", AddItem(nil, Item{Name: "Чай", PriceKopecks: 100}))
}
