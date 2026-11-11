package models

type Sku int

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
