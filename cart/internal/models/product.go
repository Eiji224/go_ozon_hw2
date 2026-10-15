package models

type Sku int

type Product struct {
	Sku   Sku    `json:"sku"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Price int    `json:"price"`
}
