package product

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"go.uber.org/fx"
)

type Handler struct {
	service ports.ProductService
}

type HandlerParams struct {
	fx.In
	Router  fiber.Router
	Service ports.ProductService
}

func New(params HandlerParams) (*Handler, error) {
	handler := Handler{params.Service}

	group := params.Router.Group("/")
	group.Get("/stocks/:id/products", handler.ByStockID)
	group.Get("/products/:id", handler.ByID)
	group.Post("/products", handler.Create)
	group.Put("/products/:id", handler.Update)
	group.Delete("/products/:id", handler.Delete)

	return &handler, nil
}

// ByStockID godoc
// @Summary      List products by stock
// @Description  Returns all products associated with a given stock ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Stock ID"
// @Success      200  {object}  ports.DTOResponseByStockIDProduct
// @Failure      400  {object}  ports.Error  "Missing or invalid stock ID"
// @Failure      404  {object}  ports.Error  "Stock not found"
// @Failure      500  {object}  ports.Error  "Internal server error"
// @Router       /stocks/{id}/products [get]
func (handler *Handler) ByStockID(c fiber.Ctx) error {
	stockByID := c.Params("id", "")

	if stockByID == "" {
		return ports.NewBadRequestError(nil)
	}

	products, err := handler.service.ByStockID(c.Context(), stockByID)
	if err != nil {
		return err
	}

	list := make([]ports.DTOProduct, 0, len(products))

	for _, product := range products {
		list = append(list, ports.DTOProduct{
			ID:      product.ID,
			Name:    product.Name,
			Qtd:     product.Qtd,
			Price:   product.Price,
			StockID: product.StockID,
		})
	}

	return c.Status(fiber.StatusOK).JSON(ports.DTOResponseByStockIDProduct{
		Data: list,
	})
}

// ByID godoc
// @Summary      Get a product by ID
// @Description  Returns a single product by its unique identifier
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  ports.DTOResponseByIDProduct
// @Failure      400  {object}  ports.Error  "Missing or invalid product ID"
// @Failure      404  {object}  ports.Error  "Product not found"
// @Failure      500  {object}  ports.Error  "Internal server error"
// @Router       /products/{id} [get]
func (handler *Handler) ByID(c fiber.Ctx) error {
	productID := c.Params("id", "")

	if productID == "" {
		return ports.NewBadRequestError(nil)
	}

	product, err := handler.service.ByID(c.Context(), productID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ports.DTOResponseByIDProduct{
		Data: ports.DTOProduct{
			ID:      product.ID,
			Name:    product.Name,
			Qtd:     product.Qtd,
			Price:   product.Price,
			StockID: product.StockID,
		},
	})
}

// Create godoc
// @Summary      Create a product
// @Description  Creates a new product and associates it with a stock
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body  body      ports.DTORequestCreateProduct  true  "Product payload"
// @Success      201   {object}  ports.DTOResponseCreateProduct
// @Failure      400   {object}  ports.Error  "Invalid request body"
// @Failure      500   {object}  ports.Error  "Internal server error"
// @Router       /products [post]
func (handler *Handler) Create(c fiber.Ctx) error {
	var body ports.DTORequestCreateProduct
	if err := c.Bind().JSON(&body); err != nil {
		return ports.NewBadRequestError(err)
	}

	domain := domain.Product{
		Name:    body.Name,
		Qtd:     body.Qtd,
		Price:   body.Price,
		StockID: body.StockID,
	}

	product, err := handler.service.Create(c.Context(), &domain)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(ports.DTOResponseCreateProduct{
		Data: ports.DTOProduct{
			ID:      product.ID,
			Name:    product.Name,
			Qtd:     product.Qtd,
			Price:   product.Price,
			StockID: product.StockID,
		},
	})
}

// Update godoc
// @Summary      Update a product
// @Description  Updates an existing product's fields by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id    path      string                         true  "Product ID"
// @Param        body  body      ports.DTORequestUpdateProduct  true  "Updated product payload"
// @Success      200   "Product updated successfully"
// @Failure      400   {object}  ports.Error  "Missing or invalid product ID / request body"
// @Failure      404   {object}  ports.Error  "Product not found"
// @Failure      500   {object}  ports.Error  "Internal server error"
// @Router       /products/{id} [put]
func (handler *Handler) Update(c fiber.Ctx) error {
	var body ports.DTORequestUpdateProduct
	productID := c.Params("id", "")

	if productID == "" {
		return ports.NewBadRequestError(nil)
	}

	if err := c.Bind().JSON(&body); err != nil {
		return ports.NewBadRequestError(err)
	}

	domain := domain.Product{
		ID:      productID,
		Name:    body.Name,
		Qtd:     body.Qtd,
		Price:   body.Price,
		StockID: body.StockID,
	}

	err := handler.service.Update(c.Context(), &domain)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// Delete godoc
// @Summary      Delete a product
// @Description  Permanently removes a product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  "Product deleted successfully"
// @Failure      400  {object}  ports.Error  "Missing or invalid product ID"
// @Failure      404  {object}  ports.Error  "Product not found"
// @Failure      500  {object}  ports.Error  "Internal server error"
// @Router       /products/{id} [delete]
func (handler *Handler) Delete(c fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return ports.NewBadRequestError(nil)
	}

	if err := handler.service.Delete(c.Context(), productID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
