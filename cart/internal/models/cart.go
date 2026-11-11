package models

type Cart map[Sku]uint32

func (c *Cart) GetSkuSlice() []Sku {
	skus := []Sku{}

	for sku := range *c {
		skus = append(skus, sku)
	}

	return skus
}
