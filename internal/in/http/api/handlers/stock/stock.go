package stock

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"go.uber.org/fx"
)

type Handler struct {
	service ports.StockService
}

type HandlerParams struct {
	fx.In
	Router  fiber.Router
	Service ports.StockService
}

func New(params HandlerParams) (*Handler, error) {
	handler := Handler{
		service: params.Service,
	}

	group := params.Router.Group("stocks")
	{
		group.Get("/", handler.All)
		group.Post("/", handler.Create)
		group.Get("/:id", handler.ByID)
		group.Put("/:id", handler.Update)
		group.Delete("/:id", handler.Delete)
	}

	return &handler, nil
}

// All godoc
// @Summary      List all stocks
// @Description  Returns a list of all available stocks with their current capacity usage
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Success      200  {object}  ports.DTOResponseAllStock
// @Failure      500  {object}  ports.Error  "Internal server error"
// @Router       /stocks [get]
func (handler *Handler) All(c fiber.Ctx) error {
	list, err := handler.service.All(c.Context())
	if err != nil {
		return err
	}

	items := make([]ports.DTOStock, 0, len(list))
	for _, stock := range list {
		from := ports.DTOStock{
			ID:           stock.ID,
			Name:         stock.Name,
			Capacity:     stock.Capacity,
			UsedCapacity: stock.Occupancy(),
		}

		items = append(items, from)
	}

	return c.Status(fiber.StatusOK).JSON(ports.DTOResponseAllStock{Data: items})
}

// ByID godoc
// @Summary      Get a stock by ID
// @Description  Returns a single stock with its current occupancy by its unique identifier
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Stock ID"
// @Success      200  {object}  ports.DTOResponseByIDStock
// @Failure      400  {object}  ports.Error  "Missing or invalid stock ID"
// @Failure      404  {object}  ports.Error  "Stock not found"
// @Failure      500  {object}  ports.Error  "Internal server error"
// @Router       /stocks/{id} [get]
func (handler *Handler) ByID(c fiber.Ctx) error {
	stockID := c.Params("id", "")
	if stockID == "" {
		return ports.NewBadRequestError(nil)
	}

	stock, err := handler.service.ByID(c.Context(), stockID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ports.DTOResponseByIDStock{Data: ports.DTOStock{
		ID:           stock.ID,
		Name:         stock.Name,
		Capacity:     stock.Capacity,
		UsedCapacity: stock.Occupancy(),
	}})
}

// Create godoc
// @Summary      Create a stock
// @Description  Creates a new stock with the given name and total capacity
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        body  body      ports.DTORequestCreateStock  true  "Stock payload"
// @Success      201   {object}  ports.DTOResponseCreateStock
// @Failure      400   {object}  ports.Error  "Invalid request body"
// @Failure      500   {object}  ports.Error  "Internal server error"
// @Router       /stocks [post]
func (handler *Handler) Create(c fiber.Ctx) error {
	var body ports.DTORequestCreateStock
	if err := c.Bind().JSON(&body); err != nil {
		return ports.NewBadRequestError(err)
	}

	domain := domain.Stock{
		Name:     body.Name,
		Capacity: body.Capacity,
	}

	stock, err := handler.service.Create(c.Context(), &domain)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(ports.DTOResponseCreateStock{
		Data: ports.DTOStock{
			ID:       stock.ID,
			Name:     stock.Name,
			Capacity: stock.Capacity,
		},
	})
}

// Update godoc
// @Summary      Update a stock
// @Description  Updates the name and/or capacity of an existing stock by its ID
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        id    path      string                       true  "Stock ID"
// @Param        body  body      ports.DTORequestUpdateStock  true  "Updated stock payload"
// @Success      200   "Stock updated successfully"
// @Failure      400   {object}  ports.Error  "Missing or invalid stock ID / request body"
// @Failure      404   {object}  ports.Error  "Stock not found"
// @Failure      500   {object}  ports.Error  "Internal server error"
// @Router       /stocks/{id} [put]
func (handler *Handler) Update(c fiber.Ctx) error {
	stockID := c.Params("id", "")
	if stockID == "" {
		return ports.NewBadRequestError(nil)
	}

	var body ports.DTORequestUpdateStock
	if err := c.Bind().JSON(&body); err != nil {
		return ports.NewBadRequestError(err)
	}

	fromDomain := domain.Stock{
		ID:       stockID,
		Name:     body.Name,
		Capacity: body.Capacity,
	}

	if err := handler.service.Update(c.Context(), &fromDomain); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// Delete godoc
// @Summary      Delete a stock
// @Description  Permanently removes a stock by its ID
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Stock ID"
// @Success      200  "Stock deleted successfully"
// @Failure      400  {object}  ports.Error  "Missing or invalid stock ID"
// @Failure      404  {object}  ports.Error  "Stock not found"
// @Failure      500  {object}  ports.Error  "Internal server error"
// @Router       /stocks/{id} [delete]
func (handler *Handler) Delete(c fiber.Ctx) error {
	stockID := c.Params("id", "")

	if stockID == "" {
		return ports.NewBadRequestError(nil)
	}

	if err := handler.service.Remove(c.Context(), stockID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
