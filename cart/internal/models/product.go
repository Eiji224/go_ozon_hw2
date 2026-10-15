package models

type Sku int

type Product struct {
	Sku   Sku
	Name  string
	Count int `json:"count"`
	Price int
}
