package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/domain/factory"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/rickferrdev/gostock/internal/core/ports/helpers"
	"github.com/rickferrdev/gostock/internal/infra/server"
	"github.com/rickferrdev/gostock/internal/platform/validator"
	"github.com/rickferrdev/gostock/internal/tests/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type test struct {
	suite.Suite
	ctrl    *gomock.Controller
	service *mocks.MockProductService
	app     *fiber.App
	router  fiber.Router
}

func (suite *test) SetupTest() {
	suite.app, suite.router = server.New(server.ServerParams{
		Validator: validator.New(),
	})

	suite.ctrl = gomock.NewController(suite.T())
	suite.service = mocks.NewMockProductService(suite.ctrl)
	_, err := New(HandlerParams{
		Router:  suite.router,
		Service: suite.service,
	})

	suite.Require().NoError(err)
}

func (suite *test) TearDownTest() {
	err := suite.app.Shutdown()
	suite.Require().NoError(err)
}

func (suite *test) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *test) TestByStockID() {
	stockID, productID := uuid.NewString(), uuid.NewString()
	product := factory.NewProductList(5, productID, stockID)
	responseByStockIDProduct := helpers.ToDTOResponseByStockIDProduct(product)

	table := []struct {
		name              string
		desc              string
		url               string
		input             string
		expectOutput      ports.DTOResponseByStockIDProduct
		expectStatus      int
		productMockReturn []*domain.Product
		stockMockErr      error
	}{
		{
			name:              "Success",
			desc:              "should return all products for a valid stock ID",
			url:               "/api/v1/stocks/%s/products",
			input:             stockID,
			productMockReturn: product,
			expectStatus:      fiber.StatusOK,
			expectOutput:      responseByStockIDProduct,
		},
		{
			name:         "InvalidStockID",
			desc:         "should return 400 when stock ID is not a valid UUID",
			url:          "/api/v1/stocks/%s/products",
			input:        "bad",
			stockMockErr: ports.NewBadRequestError(nil),
			expectOutput: ports.DTOResponseByStockIDProduct{},
			expectStatus: fiber.StatusBadRequest,
		},
		{
			name:         "InternalError",
			desc:         "should return 500 when service fails unexpectedly",
			url:          "/api/v1/stocks/%s/products",
			input:        "abc",
			stockMockErr: ports.NewInternalError(nil),
			expectOutput: ports.DTOResponseByStockIDProduct{},
			expectStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				ByStockID(gomock.Any(), tt.input).
				Return(tt.productMockReturn, tt.stockMockErr).
				Times(1)

			url := fmt.Sprintf(tt.url, tt.input)

			request := httptest.NewRequest(fiber.MethodGet, url, nil)
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func (suite *test) TestByID() {
	stockID := uuid.NewString()
	product := factory.NewProduct(uuid.NewString(), uuid.NewString())
	responseByIDProduct := helpers.ToDTOResponseByIDProduct(product)

	table := []struct {
		name              string
		desc              string
		url               string
		input             string
		expectOutput      ports.DTOResponseByIDProduct
		expectStatus      int
		productMockReturn *domain.Product
		stockMockErr      error
	}{
		{
			name:              "Success",
			desc:              "should return the product for a valid ID",
			url:               "/api/v1/products/%s",
			input:             stockID,
			productMockReturn: product,
			expectStatus:      fiber.StatusOK,
			expectOutput:      responseByIDProduct,
		},
		{
			name:         "NotFound",
			desc:         "should return 404 when product does not exist",
			url:          "/api/v1/products/%s",
			input:        stockID,
			stockMockErr: ports.NewNotFoundError(nil),
			expectOutput: ports.DTOResponseByIDProduct{},
			expectStatus: fiber.StatusNotFound,
		},
		{
			name:         "InternalError",
			desc:         "should return 500 when service fails unexpectedly",
			url:          "/api/v1/products/%s",
			input:        stockID,
			stockMockErr: ports.NewInternalError(nil),
			expectOutput: ports.DTOResponseByIDProduct{},
			expectStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				ByID(gomock.Any(), tt.input).
				Return(tt.productMockReturn, tt.stockMockErr).
				Times(1)

			url := fmt.Sprintf(tt.url, tt.input)

			request := httptest.NewRequest(fiber.MethodGet, url, nil)
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func (suite *test) TestCreate() {
	product := factory.NewProduct(uuid.NewString(), uuid.NewString())
	responseCreate := helpers.ToDTOResponseCreate(product)

	table := []struct {
		name              string
		desc              string
		url               string
		input             *domain.Product
		expectOutput      ports.DTOResponseCreateProduct
		expectStatus      int
		productMockReturn *domain.Product
		stockMockerr      error
	}{
		{
			name:              "Success",
			desc:              "should create the product and return 201",
			url:               "/api/v1/products",
			input:             product,
			productMockReturn: product,
			expectStatus:      fiber.StatusCreated,
			expectOutput:      responseCreate,
		},
		{
			name:              "BadRequest",
			desc:              "should return 400 when request payload is invalid",
			url:               "/api/v1/products",
			input:             product,
			productMockReturn: &domain.Product{},
			stockMockerr:      ports.NewBadRequestError(nil),
			expectOutput:      ports.DTOResponseCreateProduct{},
			expectStatus:      fiber.StatusBadRequest,
		},
		{
			name:              "InternalError",
			desc:              "should return 500 when service fails unexpectedly",
			url:               "/api/v1/products",
			input:             product,
			productMockReturn: &domain.Product{},
			stockMockerr:      ports.NewInternalError(nil),
			expectOutput:      ports.DTOResponseCreateProduct{},
			expectStatus:      fiber.StatusInternalServerError,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(tt.productMockReturn, tt.stockMockerr).
				Times(1)

			toJson, err := json.Marshal(ports.DTORequestCreateProduct{
				Name:    product.Name,
				Qtd:     product.Qtd,
				Price:   product.Price,
				StockID: product.StockID,
			})

			suite.NoError(err)

			body := bytes.NewBufferString(string(toJson))

			request := httptest.NewRequest(fiber.MethodPost, tt.url, body)
			request.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func (suite *test) TestUpdate() {
	product := factory.NewProduct(uuid.NewString(), uuid.NewString())

	table := []struct {
		name         string
		desc         string
		url          string
		input        *domain.Product
		expectStatus int
		stockMockErr error
	}{
		{
			name:         "Success",
			desc:         "should update the product and return 200",
			url:          "/api/v1/products/%s",
			input:        product,
			expectStatus: fiber.StatusOK,
		},
		{
			name:         "BadRequest",
			desc:         "should return 400 when request payload is invalid",
			url:          "/api/v1/products/%s",
			input:        product,
			stockMockErr: ports.NewBadRequestError(nil),
			expectStatus: fiber.StatusBadRequest,
		},
		{
			name:         "InternalError",
			desc:         "should return 500 when service fails unexpectedly",
			url:          "/api/v1/products/%s",
			input:        product,
			stockMockErr: ports.NewInternalError(nil),
			expectStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(tt.stockMockErr).
				Times(1)

			toJson, err := json.Marshal(ports.DTORequestCreateProduct{
				Name:    product.Name,
				Qtd:     product.Qtd,
				Price:   product.Price,
				StockID: product.StockID,
			})

			suite.NoError(err)

			url := fmt.Sprintf(tt.url, product.ID)

			body := bytes.NewBufferString(string(toJson))

			request := httptest.NewRequest(fiber.MethodPut, url, body)
			request.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func (suite *test) TestDelete() {
	product := factory.NewProduct(uuid.NewString(), uuid.NewString())

	table := []struct {
		name         string
		desc         string
		url          string
		input        *domain.Product
		expectStatus int
		stockMockErr error
	}{
		{
			name:         "Success",
			desc:         "should delete the product and return 200",
			url:          "/api/v1/products/%s",
			input:        product,
			expectStatus: fiber.StatusOK,
			stockMockErr: nil,
		},
		{
			name:         "NotFound",
			desc:         "should return 404 when product does not exist",
			url:          "/api/v1/products/%s",
			input:        product,
			stockMockErr: ports.NewNotFoundError(nil),
			expectStatus: fiber.StatusNotFound,
		},
		{
			name:         "InternalError",
			desc:         "should return 500 when service fails unexpectedly",
			url:          "/api/v1/products/%s",
			input:        product,
			stockMockErr: ports.NewInternalError(nil),
			expectStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				Delete(gomock.Any(), tt.input.ID).
				Return(tt.stockMockErr).
				Times(1)

			url := fmt.Sprintf(tt.url, product.ID)

			request := httptest.NewRequest(fiber.MethodDelete, url, nil)
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
