package domain

type Product struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Qtd     uint64 `json:"qtd"`
	Price   int64  `json:"price"`
	StockID string `json:"stock_id"`
}

func NewProduct(id string, name string, qtd uint64, price int64) Product {
	return Product{
		ID:    id,
		Name:  name,
		Qtd:   qtd,
		Price: price,
	}
}
