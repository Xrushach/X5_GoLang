package main

import (
	"fmt"
	"os"
)

func addStatus(names map[string]string) { names["paid"] = "оплачен" }

func main() {
	var names map[string]string
	name, ok := names["new"]
	fmt.Printf("чтение nil-map: значение=%q найдено=%t len=%d\n", name, ok, len(names))

	// Отдельный намеренно падающий запуск: go run ./ide/08-map panic
	if len(os.Args) > 1 && os.Args[1] == "panic" {
		names["new"] = "новый"
		return
	}

	names = make(map[string]string)
	names["new"] = "новый"
	names["empty"] = ""
	for _, code := range []string{"new", "empty", "missing"} {
		name, ok := names[code]
		fmt.Printf("ключ=%q значение=%q найдено=%t\n", code, name, ok)
	}

	addStatus(names)
	fmt.Println("изменение из функции:", names["paid"])
	delete(names, "empty")
	fmt.Println("после delete, len:", len(names))
	for code, name := range names {
		fmt.Println("пара:", code, name)
	}
	// Порядок обхода map не задан, одинакового порядка в запусках не обещаем.
}
