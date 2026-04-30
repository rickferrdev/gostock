package domain

type Stock struct {
	ID       string
	Name     string
	Capacity uint64
	Items    []*Product
}

func NewStock(id string, name string, capacity uint64) Stock {
	return Stock{
		ID:       id,
		Name:     name,
		Capacity: capacity,
		Items:    make([]*Product, 0),
	}
}

func (stock *Stock) Push(item *Product) {
	stock.Items = append(stock.Items, item)
}

func (stock *Stock) Occupancy() uint64 {
	var total uint64

	for _, product := range stock.Items {
		total += product.Qtd
	}

	return total
}

func (stock *Stock) AvailableSpace() uint64 {
	occupied := stock.Occupancy()
	if occupied >= stock.Capacity {
		return 0
	}

	return stock.Capacity - occupied
}

func (stock *Stock) CanAdd(quantity uint64) bool {
	return stock.AvailableSpace() >= quantity
}
