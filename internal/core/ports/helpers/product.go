package helpers

import (
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/ports"
)

func ToDTOProduct(product *domain.Product) ports.DTOProduct {
	return ports.DTOProduct{
		ID:      product.ID,
		Name:    product.Name,
		Qtd:     product.Qtd,
		Price:   product.Price,
		StockID: product.StockID,
	}
}

func ToDTOResponseByStockIDProduct(products []*domain.Product) ports.DTOResponseByStockIDProduct {
	list := make([]ports.DTOProduct, 0, len(products))
	for _, pro := range products {
		dto := ToDTOProduct(pro)
		list = append(list, dto)
	}

	return ports.DTOResponseByStockIDProduct{Data: list}
}

func ToDTOResponseByIDProduct(product *domain.Product) ports.DTOResponseByIDProduct {
	return ports.DTOResponseByIDProduct{Data: ToDTOProduct(product)}
}

func ToDTOResponseCreate(product *domain.Product) ports.DTOResponseCreateProduct {
	return ports.DTOResponseCreateProduct{Data: ToDTOProduct(product)}
}

func ToDTORequestCreate(product *domain.Product) ports.DTORequestCreateProduct {
	return ports.DTORequestCreateProduct{
		Name:    product.Name,
		Qtd:     product.Qtd,
		Price:   product.Price,
		StockID: product.StockID,
	}
}

func ToDTORequestUpdateProduct(product *domain.Product) ports.DTORequestUpdateProduct {
	return ports.DTORequestUpdateProduct{
		Name:    product.Name,
		Qtd:     product.Qtd,
		Price:   product.Price,
		StockID: product.StockID,
	}
}
