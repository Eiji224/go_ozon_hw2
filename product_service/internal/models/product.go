package models

type Sku int

type Product struct {
	Sku   Sku    `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
