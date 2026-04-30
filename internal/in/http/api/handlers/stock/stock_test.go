package stock

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
	service *mocks.MockStockService
	app     *fiber.App
	router  fiber.Router
}

func (suite *test) SetupTest() {
	suite.app, suite.router = server.New(server.ServerParams{
		Validator: validator.New(),
	})

	suite.ctrl = gomock.NewController(suite.T())
	suite.service = mocks.NewMockStockService(suite.ctrl)
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

func (suite *test) TestAll() {
	stockID := uuid.NewString()
	stock := factory.NewStockList(5, stockID)
	responseAllStock := helpers.ToDTOResponseAllStock(stock)

	table := []struct {
		name          string
		desc          string
		url           string
		expectStatus  int
		expectBody    ports.DTOResponseAllStock
		stockMockList []*domain.Stock
		stockMockErr  error
	}{
		{
			name:          "Success",
			desc:          "should return all stocks and 200",
			url:           "/api/v1/stocks",
			expectStatus:  fiber.StatusOK,
			expectBody:    responseAllStock,
			stockMockList: stock,
			stockMockErr:  nil,
		},
		{
			name:          "InternalError",
			desc:          "should return 500 when service fails unexpectedly",
			url:           "/api/v1/stocks",
			expectStatus:  fiber.StatusInternalServerError,
			stockMockList: []*domain.Stock{},
			stockMockErr:  ports.NewInternalError(nil),
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				All(gomock.Any()).
				Return(tt.stockMockList, tt.stockMockErr).
				Times(1)

			request := httptest.NewRequest(fiber.MethodGet, tt.url, nil)
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func (suite *test) TestByID() {
	stockID := uuid.NewString()
	stock := factory.NewStock(stockID)
	responseByIDStock := helpers.ToDTOResponseByIDStock(stock)

	table := []struct {
		name         string
		desc         string
		url          string
		input        string
		expectStatus int
		expectBody   ports.DTOResponseByIDStock
		stockMock    *domain.Stock
		stockMockErr error
	}{
		{
			name:         "Success",
			desc:         "should return the stock for a valid ID",
			url:          "/api/v1/stocks/%s",
			input:        stockID,
			stockMock:    stock,
			stockMockErr: nil,
			expectBody:   responseByIDStock,
			expectStatus: fiber.StatusOK,
		},
		{
			name:         "InvalidStockID",
			desc:         "should return 400 when stock ID is not a valid UUID",
			url:          "/api/v1/stocks/%s",
			input:        "abc",
			stockMock:    &domain.Stock{},
			stockMockErr: ports.NewBadRequestError(nil),
			expectStatus: fiber.StatusBadRequest,
		},
		{
			name:         "NotFound",
			desc:         "should return 404 when stock does not exist",
			url:          "/api/v1/stocks/%s",
			input:        stockID,
			stockMock:    &domain.Stock{},
			stockMockErr: ports.NewNotFoundError(nil),
			expectStatus: fiber.StatusNotFound,
		},
		{
			name:         "InternalError",
			desc:         "should return 500 when service fails unexpectedly",
			url:          "/api/v1/stocks/%s",
			input:        stockID,
			stockMock:    &domain.Stock{},
			stockMockErr: ports.NewInternalError(nil),
			expectStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				ByID(gomock.Any(), tt.input).
				Return(tt.stockMock, tt.stockMockErr).
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
	stockID := uuid.NewString()
	stock := factory.NewStock(stockID)
	requestCreateStock := helpers.ToDTORequestCreateStock(stock)
	responseCreateStock := helpers.ToDTOResponseCreateStock(stock)

	table := []struct {
		name         string
		desc         string
		url          string
		expectStatus int
		expectBody   ports.DTOResponseCreateStock
		stockMock    *domain.Stock
		stockMockErr error
	}{
		{
			name:         "Success",
			desc:         "should create the stock and return 201",
			url:          "/api/v1/stocks",
			stockMock:    stock,
			stockMockErr: nil,
			expectBody:   responseCreateStock,
			expectStatus: fiber.StatusCreated,
		},
		{
			name:         "BadRequest",
			desc:         "should return 400 when request payload is invalid",
			url:          "/api/v1/stocks",
			stockMock:    &domain.Stock{},
			stockMockErr: ports.NewBadRequestError(nil),
			expectStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(tt.stockMock, tt.stockMockErr).
				Times(1)

			toJSON, err := json.Marshal(requestCreateStock)
			suite.NoError(err)

			request := httptest.NewRequest(fiber.MethodPost, tt.url, bytes.NewBufferString(string(toJSON)))
			request.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func (suite *test) TestUpdate() {
	stockID := uuid.NewString()
	stock := factory.NewStock(stockID)
	requestUpdateStock := helpers.ToDTORequestUpdateStock(stock)

	table := []struct {
		name         string
		desc         string
		url          string
		inputID      string
		expectStatus int
		stockMockErr error
	}{
		{
			name:         "Success",
			desc:         "should update the stock and return 200",
			url:          "/api/v1/stocks/%s",
			inputID:      stockID,
			expectStatus: fiber.StatusOK,
			stockMockErr: nil,
		},
		{
			name:         "NotFound",
			desc:         "should return 404 when stock does not exist",
			url:          "/api/v1/stocks/%s",
			inputID:      stockID,
			expectStatus: fiber.StatusNotFound,
			stockMockErr: ports.NewNotFoundError(nil),
		},
		{
			name:         "InternalError",
			desc:         "should return 500 when service fails unexpectedly",
			url:          "/api/v1/stocks/%s",
			inputID:      stockID,
			expectStatus: fiber.StatusInternalServerError,
			stockMockErr: ports.NewInternalError(nil),
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				Update(gomock.Any(), gomock.Any()).
				Return(tt.stockMockErr).
				Times(1)

			toJSON, err := json.Marshal(requestUpdateStock)
			suite.NoError(err)

			url := fmt.Sprintf(tt.url, tt.inputID)

			request := httptest.NewRequest(fiber.MethodPut, url, bytes.NewBufferString(string(toJSON)))
			request.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expectStatus, resp.StatusCode, tt.desc)
		})
	}
}

func (suite *test) TestDelete() {
	stockID := uuid.NewString()

	table := []struct {
		name         string
		desc         string
		url          string
		inputID      string
		expectStatus int
		stockMockErr error
	}{
		{
			name:         "Success",
			desc:         "should delete the stock and return 200",
			url:          "/api/v1/stocks/%s",
			inputID:      stockID,
			expectStatus: fiber.StatusOK,
			stockMockErr: nil,
		},
		{
			name:         "InvalidStockID",
			desc:         "should return 400 when stock ID is not a valid UUID",
			url:          "/api/v1/stocks/%s",
			inputID:      "abc",
			expectStatus: fiber.StatusBadRequest,
			stockMockErr: ports.NewBadRequestError(nil),
		},
		{
			name:         "NotFound",
			desc:         "should return 404 when stock does not exist",
			url:          "/api/v1/stocks/%s",
			inputID:      stockID,
			expectStatus: fiber.StatusNotFound,
			stockMockErr: ports.NewNotFoundError(nil),
		},
		{
			name:         "InternalError",
			desc:         "should return 500 when service fails unexpectedly",
			url:          "/api/v1/stocks/%s",
			inputID:      stockID,
			expectStatus: fiber.StatusInternalServerError,
			stockMockErr: ports.NewInternalError(nil),
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			suite.service.EXPECT().
				Remove(gomock.Any(), tt.inputID).
				Return(tt.stockMockErr).
				Times(1)

			url := fmt.Sprintf(tt.url, tt.inputID)
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
