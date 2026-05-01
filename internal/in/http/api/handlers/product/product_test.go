package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain/factory"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/rickferrdev/gostock/internal/core/ports/helpers"
	"github.com/rickferrdev/gostock/internal/infra/server"
	"github.com/rickferrdev/gostock/internal/platform/validator"
	"github.com/rickferrdev/gostock/internal/tests/mocks"
	"github.com/stretchr/testify/require"
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

func (suite *test) RequirePortError(t *testing.T, body []byte, expected *ports.Error) {
	var got *ports.Error
	err := json.Unmarshal(body, &got)
	require.NoError(t, err)

	suite.Require().Equal(expected.Code, got.Code)
	suite.Require().Equal(expected.Message, got.Message)
}

func (suite *test) TestByStockID() {
	stockID, productID := uuid.NewString(), uuid.NewString()
	productSlice := factory.NewProductSlice(5, productID, stockID)

	table := []struct {
		name   string
		desc   string
		url    string
		input  string
		status int
		err    *ports.Error
		setup  func() ports.DTOResponseByStockIDProduct
	}{
		{
			name:   "Success",
			desc:   "should return all products for a valid stock ID",
			url:    "/api/v1/stocks/%s/products",
			input:  stockID,
			status: fiber.StatusOK,
			setup: func() ports.DTOResponseByStockIDProduct {
				suite.service.EXPECT().ByStockID(gomock.Any(), stockID).Times(1).Return(productSlice, nil)
				return helpers.ToDTOResponseByStockIDProduct(productSlice)
			},
		},
		{
			name:   "InvalidStockID",
			desc:   "should return 400 when stock ID is not a valid UUID",
			url:    "/api/v1/stocks/%s/products",
			input:  "bad",
			status: fiber.StatusBadRequest,
			err:    ports.NewBadRequestError(nil),
			setup: func() ports.DTOResponseByStockIDProduct {
				suite.service.EXPECT().ByStockID(gomock.Any(), "bad").Times(1).Return(nil, ports.NewBadRequestError(nil))
				return ports.DTOResponseByStockIDProduct{}
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when service fails unexpectedly",
			url:    "/api/v1/stocks/%s/products",
			input:  stockID,
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() ports.DTOResponseByStockIDProduct {
				suite.service.EXPECT().ByStockID(gomock.Any(), stockID).Return(nil, ports.NewInternalError(nil)).Times(1)
				return ports.DTOResponseByStockIDProduct{}
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			url := fmt.Sprintf(tt.url, tt.input)
			expect := tt.setup()

			request := httptest.NewRequest(fiber.MethodGet, url, nil)
			resp, err := suite.app.Test(request)
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}

			var got ports.DTOResponseByStockIDProduct
			suite.Require().NoError(json.Unmarshal(body, &got))
			suite.Require().Equal(expect, got, tt.desc)
		})
	}
}

func (suite *test) TestByID() {
	id := uuid.NewString()
	product := factory.NewProduct(id, uuid.NewString(), "Product", 10, 1000)

	table := []struct {
		name   string
		desc   string
		url    string
		input  string
		status int
		err    *ports.Error
		setup  func() ports.DTOResponseByIDProduct
	}{
		{
			name:   "Success",
			desc:   "should return the product for a valid ID",
			url:    "/api/v1/products/%s",
			input:  id,
			status: fiber.StatusOK,
			setup: func() ports.DTOResponseByIDProduct {
				suite.service.EXPECT().ByID(gomock.Any(), id).Times(1).Return(product, nil)
				return helpers.ToDTOResponseByIDProduct(product)
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 when product does not exist",
			url:    "/api/v1/products/%s",
			input:  id,
			status: fiber.StatusNotFound,
			err:    ports.NewNotFoundError(nil),
			setup: func() ports.DTOResponseByIDProduct {
				suite.service.EXPECT().ByID(gomock.Any(), id).Times(1).Return(nil, ports.NewNotFoundError(nil))
				return ports.DTOResponseByIDProduct{}
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when service fails unexpectedly",
			url:    "/api/v1/products/%s",
			input:  id,
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() ports.DTOResponseByIDProduct {
				suite.service.EXPECT().ByID(gomock.Any(), id).Times(1).Return(nil, ports.NewInternalError(nil))
				return ports.DTOResponseByIDProduct{}
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			url := fmt.Sprintf(tt.url, tt.input)
			expect := tt.setup()

			request := httptest.NewRequest(fiber.MethodGet, url, nil)
			resp, err := suite.app.Test(request)
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}

			var got ports.DTOResponseByIDProduct
			suite.Require().NoError(json.Unmarshal(body, &got))
			suite.Require().Equal(expect, got, tt.desc)
		})
	}
}

func (suite *test) TestCreate() {
	product := factory.NewProduct("", uuid.NewString(), "Product", 10, 1000)

	table := []struct {
		name   string
		desc   string
		url    string
		input  ports.DTORequestCreateProduct
		status int
		err    *ports.Error
		setup  func() ports.DTOResponseCreateProduct
	}{
		{
			name:   "Success",
			desc:   "should create and return product",
			url:    "/api/v1/products",
			input:  helpers.ToDTORequestCreate(product),
			status: fiber.StatusCreated,
			setup: func() ports.DTOResponseCreateProduct {
				suite.service.EXPECT().Create(gomock.Any(), product).Times(1).Return(product, nil)
				return helpers.ToDTOResponseCreate(product)
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 if dependency not found",
			url:    "/api/v1/products",
			input:  helpers.ToDTORequestCreate(product),
			status: fiber.StatusNotFound,
			err:    ports.NewNotFoundError(nil),
			setup: func() ports.DTOResponseCreateProduct {
				suite.service.EXPECT().Create(gomock.Any(), product).Times(1).Return(nil, ports.NewNotFoundError(nil))
				return ports.DTOResponseCreateProduct{}
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when service fails",
			url:    "/api/v1/products",
			input:  helpers.ToDTORequestCreate(product),
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() ports.DTOResponseCreateProduct {
				suite.service.EXPECT().Create(gomock.Any(), product).Times(1).Return(nil, ports.NewInternalError(nil))
				return ports.DTOResponseCreateProduct{}
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expect := tt.setup()
			send, err := json.Marshal(tt.input)
			suite.Require().NoError(err)

			request := httptest.NewRequest(fiber.MethodPost, tt.url, bytes.NewBuffer(send))
			request.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(request)
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}

			var got ports.DTOResponseCreateProduct
			suite.Require().NoError(json.Unmarshal(body, &got))
			suite.Require().Equal(expect, got, tt.desc)
		})
	}
}

func (suite *test) TestUpdate() {
	input := factory.NewProduct(uuid.NewString(), uuid.NewString(), "Product", 10, 1000)

	table := []struct {
		name   string
		desc   string
		url    string
		input  ports.DTORequestUpdateProduct
		status int
		err    *ports.Error
		setup  func()
	}{
		{
			name:   "Success",
			desc:   "should update successfully",
			url:    "/api/v1/products/%s",
			input:  helpers.ToDTORequestUpdateProduct(input),
			status: fiber.StatusOK,
			setup: func() {
				suite.service.EXPECT().Update(gomock.Any(), input).Times(1).Return(nil)
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 when product does not exist",
			url:    "/api/v1/products/%s",
			input:  helpers.ToDTORequestUpdateProduct(input),
			status: fiber.StatusNotFound,
			err:    ports.NewNotFoundError(nil),
			setup: func() {
				suite.service.EXPECT().Update(gomock.Any(), input).Times(1).Return(ports.NewNotFoundError(nil))
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when service fails",
			url:    "/api/v1/products/%s",
			input:  helpers.ToDTORequestUpdateProduct(input),
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() {
				suite.service.EXPECT().Update(gomock.Any(), input).Times(1).Return(ports.NewInternalError(nil))
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()
			send, err := json.Marshal(tt.input)
			suite.Require().NoError(err)

			url := fmt.Sprintf(tt.url, input.ID)
			request := httptest.NewRequest(fiber.MethodPut, url, bytes.NewBuffer(send))
			request.Header.Set("Content-Type", "application/json")
			resp, err := suite.app.Test(request)
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}
		})
	}
}

func (suite *test) TestDelete() {
	id := uuid.NewString()

	table := []struct {
		name   string
		desc   string
		url    string
		input  string
		status int
		err    *ports.Error
		setup  func()
	}{
		{
			name:   "Success",
			desc:   "should delete successfully",
			url:    "/api/v1/products/%s",
			input:  id,
			status: fiber.StatusOK,
			setup: func() {
				suite.service.EXPECT().Delete(gomock.Any(), id).Times(1).Return(nil)
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404",
			url:    "/api/v1/products/%s",
			input:  id,
			status: fiber.StatusNotFound,
			err:    ports.NewNotFoundError(nil),
			setup: func() {
				suite.service.EXPECT().Delete(gomock.Any(), id).Times(1).Return(ports.NewNotFoundError(nil))
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500",
			url:    "/api/v1/products/%s",
			input:  id,
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() {
				suite.service.EXPECT().Delete(gomock.Any(), id).Times(1).Return(ports.NewInternalError(nil))
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()
			url := fmt.Sprintf(tt.url, tt.input)
			request := httptest.NewRequest(fiber.MethodDelete, url, nil)
			resp, err := suite.app.Test(request)
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}
		})
	}
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
