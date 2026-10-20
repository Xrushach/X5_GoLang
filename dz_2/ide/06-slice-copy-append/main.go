package main

import "fmt"

func appendLocally(values []int) {
	values = append(values, 99)
	fmt.Println("длина внутри функции:", len(values))
}

func main() {
	numbers := [3]int{10, 20, 30}
	values := numbers[:]
	another := values
	another[0] = 99
	shorter := another[:1]
	fmt.Println("копия слайса, общий массив:", numbers, values, another)
	fmt.Println("длины независимы:", len(values), len(shorter))

	items := make([]int, 0, 2)
	fmt.Println("make:", items, "len:", len(items), "cap:", cap(items))
	items = append(items, 10)
	items = append(items, 20)
	fmt.Println("два append:", items, "len:", len(items), "cap:", cap(items))
	items = append(items, 30)
	fmt.Println("рост:", items, "len:", len(items), "cap:", cap(items))
	// Конкретный cap после роста не фиксируем. Здесь достаточно cap >= len.

	base := [4]int{10, 20, 30, 40}
	left, right := base[:2], base[2:]
	left = append(left, 777)
	fmt.Println("append в общий массив:", "base:", base, "left:", left, "right:", right)

	original := make([]int, 2, 2)
	original[0], original[1] = 10, 20
	grown := append(original, 30)
	grown[0] = -1
	fmt.Println("append в новый массив:", "original:", original, "grown:", grown)

	independent := make([]int, len(values))
	copy(independent, values)
	independent[0] = 555
	fmt.Println("копия элементов:", "values:", values, "independent:", independent)

	appendLocally(items)
	fmt.Println("длина снаружи функции:", len(items))
}
