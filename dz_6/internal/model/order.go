package model

// Order продолжает модель из занятия 6.
type Order struct {
	ID       int    `json:"id"`
	Customer string `json:"customer"`
	Status   string `json:"status"`
}
