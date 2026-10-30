package main

import (
	"fmt"

	"example/internal/config"
	"example/internal/handler"
	"example/internal/repository"
	"example/internal/service"
)

func main() {
	repo := repository.NewMemory()
	orders := service.NewService(repo)
	handler := handler.NewHandler(orders)

	address := config.GetAddress()
	order, err := handler.Create("Покупатель")
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Printf("сервис слушал бы %s, создан заказ: %+v\n", address, order)
}
