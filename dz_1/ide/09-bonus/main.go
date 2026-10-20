package main

import "fmt"

func main() {
Outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Println(i, j)
			if i == 1 && j == 1 {
				break Outer
			}
		}
	}
	switch 1 {
	case 1:
		fmt.Println("первая ветка")
		fallthrough
	case 2:
		fmt.Println("вторая, без проверки условия")
	}
}
