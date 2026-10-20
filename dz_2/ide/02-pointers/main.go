package main

import (
	"errors"
	"fmt"
)

func changeCopy(value int) { value = 99 }

func changeOriginal(value *int) error {
	if value == nil {
		return errors.New("адрес не задан")
	}
	*value = 99
	return nil
}

func main() {
	number := 10
	changeCopy(number)
	fmt.Println("после копии:", number)

	pointer := &number
	fmt.Println("pointer == &number:", pointer == &number)
	fmt.Println("значение по адресу:", *pointer)
	if err := changeOriginal(pointer); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("после указателя:", number)

	var empty *int
	fmt.Println("нулевой указатель:", empty == nil)
	fmt.Println("вызов с nil:", changeOriginal(empty))
}
