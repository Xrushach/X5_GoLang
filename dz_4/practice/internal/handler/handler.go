package handler

import (
	"example/internal/model"
)

type createService interface {
	Create(customer string) (model.Order, error)
}

type Handler struct {
	orders createService
}

func NewHandler(orders createService) *Handler {
	return &Handler{orders: orders}
}

func (h *Handler) Create(customer string) (model.Order, error) {
	return h.orders.Create(customer)
}
