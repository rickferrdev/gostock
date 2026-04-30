package ports

type HealthResponse struct {
	Status   string         `json:"status"`
	Services map[string]any `json:"services"`
}

type DTOStock struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Capacity     uint64 `json:"capacity"`
	UsedCapacity uint64 `json:"used_capacity"`
}

type DTOResponseAllStock struct {
	Data []DTOStock `json:"data"`
}

type DTORequestCreateStock struct {
	Name     string `json:"name" validate:"required,min=3"`
	Capacity uint64 `json:"capacity" validate:"required,min=1"`
}

type DTOResponseCreateStock struct {
	Data DTOStock `json:"data"`
}

type DTOResponseByIDStock struct {
	Data DTOStock `json:"data"`
}

type DTORequestUpdateStock struct {
	Name     string `json:"name" validate:"required,min=3"`
	Capacity uint64 `json:"capacity" validate:"required,min=1"`
}

type DTOProduct struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Qtd     uint64 `json:"qtd"`
	Price   int64  `json:"price"`
	StockID string `json:"stock_id"`
}

type DTORequestCreateProduct struct {
	Name    string `json:"name" validate:"required,min=2"`
	Qtd     uint64 `json:"qtd" validate:"gte=0"`
	Price   int64  `json:"price" validate:"required,min=1"`
	StockID string `json:"stock_id" validate:"required,uuid"`
}

type DTOResponseCreateProduct struct {
	Data DTOProduct `json:"data"`
}

type DTOResponseByIDProduct struct {
	Data DTOProduct `json:"data"`
}

type DTOResponseByStockIDProduct struct {
	Data []DTOProduct `json:"data"`
}

type DTORequestUpdateProduct struct {
	Name    string `json:"name" validate:"min=2"`
	Qtd     uint64 `json:"qtd" validate:"gte=0"`
	Price   int64  `json:"price" validate:"min=1"`
	StockID string `json:"stock_id" validate:"uuid"`
}
