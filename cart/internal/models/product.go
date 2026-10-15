package models

type Sku int64

type Product struct {
	Sku   Sku    `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}
