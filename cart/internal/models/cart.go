package models

type Cart map[Sku]*Product

func (c Cart) GetCartSlice() []*Product {
	result := make([]*Product, 0, len(c))

	for _, product := range c {
		result = append(result, product)
	}

	return result
}
