package domain

import "time"

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	OrderDate  time.Time   `json:"order_date"`
	Status     string      `json:"status"`
	Items      []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CreateOrderRequest struct {
	CustomerID string                   `json:"customer_id" binding:"required,uuid"`
	Items      []CreateOrderItemRequest `json:"items"       binding:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity"   binding:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}
