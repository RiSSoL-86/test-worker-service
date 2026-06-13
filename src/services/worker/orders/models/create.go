package models

type CreateOrderPayload struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	Address       string `json:"address"`
	Priority      string `json:"priority"`
}
