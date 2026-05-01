package factory

import (
	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain"
)

// Return a stock with the given id, name, capacity and no products
func NewStock(id string, name string, capacity uint64) *domain.Stock {
	return &domain.Stock{
		ID:       id,
		Name:     name,
		Capacity: capacity,
		Items:    []*domain.Product{},
	}
}

// Return a product with the given id, stockID, name, qtd and price
func NewProduct(id string, stockID string, name string, qtd uint64, price int64) *domain.Product {
	return &domain.Product{
		ID:      id,
		Name:    name,
		Qtd:     qtd,
		Price:   price,
		StockID: stockID,
	}
}

// Return a stock with random ID, default name "Stock", default capacity 100 and no products
func NewStockDefault() *domain.Stock {
	return NewStock(uuid.NewString(), "Stock", 100)
}

// Return a product with random ID, default name "Product", default qtd 1, default price 100 and the given stockID
func NewProductDefault() *domain.Product {
	return NewProduct(uuid.NewString(), uuid.NewString(), "Product", 1, 100)
}

// Return a stock with random ID, default name "Stock", default capacity 100 and one product with random ID, default name "Product", default qtd 1, default price 100 and the stockID of the created stock
func NewStockWithOneProducts(qtd uint64, capacity uint64) (*domain.Stock, *domain.Product) {
	stockID := uuid.NewString()
	productID := uuid.NewString()
	stock := NewStock(stockID, "Stock", capacity)
	product := NewProduct(productID, stockID, "Product", qtd, 100)
	stock.Items = append(stock.Items, product)

	return stock, product
}

// Return a slice of stocks with the given count, all with the same stockID, default name "Stock", default capacity 100 and no products
func NewStockSlice(count int, stockID string) []*domain.Stock {
	list := make([]*domain.Stock, 0, count)
	for range count {
		list = append(list, NewStock(stockID, "Stock", 100))
	}
	return list
}

// Return a slice of products with the given count, all with the same stockID, default name "Product", default qtd 1 and default price 100
func NewProductSlice(count int, id, stockID string) []*domain.Product {
	list := make([]*domain.Product, 0, count)

	for range count {
		list = append(list, NewProduct(id, stockID, "Product", 1, 100))
	}

	return list
}
