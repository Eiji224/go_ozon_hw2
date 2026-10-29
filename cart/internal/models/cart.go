package models

import (
	"cmp"
	"slices"
)

type Cart map[Sku]*Product

func (c Cart) GetCartSlice() []*Product {
	result := make([]*Product, 0, len(c))

	for _, product := range c {
		productSnapshot := *product
		result = append(result, &productSnapshot)
	}

	slices.SortFunc(result, func(a *Product, b *Product) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	return result
}
