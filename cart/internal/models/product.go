package models

import (
	"cmp"
	"slices"
)

type Sku int64

type Product struct {
	Sku   Sku    `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

func SortProducts(products []*Product) {
	slices.SortFunc(products, func(a *Product, b *Product) int {
		return cmp.Compare(a.Sku, b.Sku)
	})
}
