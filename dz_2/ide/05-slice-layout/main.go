package main

import "fmt"

// Наша учебная схема. Это не тип []int и не доступ к внутренним полям Go.
type SliceView struct {
	First *int
	Len   int
	Cap   int
}

func main() {
	numbers := [4]int{10, 20, 30, 40}
	values := numbers[1:3]
	fmt.Println("массив:", numbers, "слайс:", values)
	fmt.Println("len:", len(values), "cap:", cap(values))

	view := SliceView{First: &values[0], Len: len(values), Cap: cap(values)}
	fmt.Println("учебная схема: первый элемент:", *view.First, "len:", view.Len, "cap:", view.Cap)
	fmt.Println("первый элемент слайса является numbers[1]:", &values[0] == &numbers[1])
	// len = 2: доступны 20 и 30. cap = 3: от начала слайса до конца массива.

	values[0] = 99
	fmt.Println("запись через слайс:", numbers, values)
	// Элементы находятся в массиве. У values нет собственных полей First/Len/Cap.

	var empty []int
	fmt.Println("nil-слайс:", empty == nil, "len:", len(empty), "cap:", cap(empty))
	// Для пустого слайса values[0] недопустим. Схему выше строили для непустого.
}
