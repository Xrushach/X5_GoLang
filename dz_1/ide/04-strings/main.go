package main

import (
	"fmt"
)

func main() {
	word := "Goшник"
	original := word
	word = "Go"
	fmt.Println("после присваивания:", word, "копия:", original)
	word = original
	fmt.Println("байт:", len(word))
	fmt.Printf("первый байт: %d, %c\n", word[0], word[0])
	count := 0
	for index, symbol := range word {
		fmt.Printf("индекс байта=%d код=%d знак=%c\n", index, symbol, symbol)
		count++
	}
	fmt.Println("кодовых точек:", count)
	bytes := []byte(word)
	runes := []rune(word)
	fmt.Println("байты:", bytes)
	fmt.Println("руны:", runes)
	fmt.Println("обратно:", string(runes))

	text := string(bytes)
	bytes[0] = 'g'
	fmt.Println("изменённые байты:", string(bytes), "строка:", text)
	fmt.Println("первые два байта:", word[:3])
	n := 65
	fmt.Printf("string(65): %q\n", string(n))
}
