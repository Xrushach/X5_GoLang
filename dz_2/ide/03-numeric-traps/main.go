package main

import "fmt"

func main() {
	//var small int8 = 127
	//small++
	//fmt.Println("переполнение int8:", small)
	//var impossible int8 = 128
	//value := 300
	//fmt.Println("300 в uint8:", uint8(value))
	//fraction := 19.9
	//fmt.Println("отбрасывание дроби:", int(fraction))
	//a, b := 0.1, 0.2
	//result := a + b
	//fmt.Printf("сумма: %.1f\n", result)
	//fmt.Println("равна 0.3:", result == 0.3)
	//fmt.Println("деление до преобразования:", float64(5/2))
	//fmt.Println("преобразование до деления:", float64(5)/2)
	priceKopecks := int64(19_990)
	quantity := int64(3)
	total := priceKopecks * quantity
	fmt.Println("копейки:", total)
	fmt.Printf("рубли: %d.%02d\n", total/100, total%100)
}
