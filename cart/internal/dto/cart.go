package dto

import "github.com/Eiji224/go_ozon/hw2/cart/internal/models"

type GetCartResponse struct {
	Items      []*models.Product `json:"items"`
	TotalPrice uint32            `json:"total_price"`
}
