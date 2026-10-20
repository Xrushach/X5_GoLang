package main

import (
	"errors"
	"fmt"
)

// Цена и результат в копейках; по условию произведение помещается в int64.
func lineTotal(price int64, quantity int) (int64, error) {
	if price < 0 {
		return 0, errors.New("цена не может быть отрицательной")
	}
	if quantity <= 0 {
		return 0, errors.New("количество должно быть положительным")
	}
	return price * int64(quantity), nil
}

func main() {
	total, err := lineTotal(199_90, 2)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("сумма в копейках:", total)

	total, err = lineTotal(199_90, 0)
	fmt.Println("неверное количество:", total, err)
}
