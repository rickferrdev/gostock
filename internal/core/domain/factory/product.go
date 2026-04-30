package factory

import (
	"github.com/rickferrdev/gostock/internal/core/domain"
)

func NewProduct(id string, stockID string) *domain.Product {
	return &domain.Product{
		ID:      id,
		Name:    "Product",
		Qtd:     100,
		Price:   50,
		StockID: stockID,
	}
}

func NewProductList(count int, id, stockID string) []*domain.Product {
	list := make([]*domain.Product, 0, count)

	for range count {
		list = append(list, NewProduct(id, stockID))
	}

	return list
}
