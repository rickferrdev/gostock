package factory

import "github.com/rickferrdev/gostock/internal/core/domain"

func NewStock(id string) *domain.Stock {
	return &domain.Stock{
		ID:       id,
		Name:     "Product",
		Capacity: 100,
		Items:    []*domain.Product{},
	}
}

func NewStockList(count int, stockID string) []*domain.Stock {
	list := make([]*domain.Stock, 0, count)
	for range count {
		list = append(list, NewStock(stockID))
	}
	return list
}
