package helpers

import (
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/ports"
)

func ToDTOResponseAllStock(stocks []*domain.Stock) ports.DTOResponseAllStock {
	list := make([]ports.DTOStock, 0, len(stocks))
	for _, pro := range stocks {
		dto := ToDTOStock(pro)
		list = append(list, dto)
	}

	return ports.DTOResponseAllStock{Data: list}
}

func ToDTOStock(stock *domain.Stock) ports.DTOStock {
	return ports.DTOStock{
		ID:           stock.ID,
		Name:         stock.Name,
		Capacity:     stock.Capacity,
		UsedCapacity: stock.Occupancy(),
	}
}

func ToDTORequestCreateStock(stock *domain.Stock) ports.DTORequestCreateStock {
	return ports.DTORequestCreateStock{
		Name:     stock.Name,
		Capacity: stock.AvailableSpace(),
	}
}

func ToDTOResponseCreateStock(stock *domain.Stock) ports.DTOResponseCreateStock {
	return ports.DTOResponseCreateStock{
		Data: ToDTOStock(stock),
	}
}

func ToDTOResponseByIDStock(stock *domain.Stock) ports.DTOResponseByIDStock {
	return ports.DTOResponseByIDStock{
		Data: ToDTOStock(stock),
	}
}

func ToDTORequestUpdateStock(stock *domain.Stock) ports.DTORequestUpdateStock {
	return ports.DTORequestUpdateStock{
		Name:     stock.Name,
		Capacity: stock.Capacity,
	}
}
