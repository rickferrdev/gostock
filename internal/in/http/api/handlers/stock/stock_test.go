package stock

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

func (suite *test) RequirePortError(t *testing.T, body []byte, expected *ports.Error) {
	var got *ports.Error
	err := json.Unmarshal(body, &got)
	require.NoError(t, err)

	suite.Require().Equal(expected.Code, got.Code)
	suite.Require().Equal(expected.Message, got.Message)
}

func (suite *test) TestAll() {
	stocks := factory.NewStockSlice(5, uuid.NewString())

	table := []struct {
		name   string
		desc   string
		status int
		err    *ports.Error
		setup  func() ports.DTOResponseAllStock
	}{
		{
			name:   "Success",
			desc:   "should return all stocks",
			status: fiber.StatusOK,
			setup: func() ports.DTOResponseAllStock {
				suite.service.EXPECT().All(gomock.Any()).Return(stocks, nil).Times(1)
				return helpers.ToDTOResponseAllStock(stocks)
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when service fails",
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() ports.DTOResponseAllStock {
				suite.service.EXPECT().All(gomock.Any()).Return(nil, ports.NewInternalError(nil)).Times(1)
				return ports.DTOResponseAllStock{}
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expect := tt.setup()

			req := httptest.NewRequest(fiber.MethodGet, "/api/v1/stocks", nil)
			resp, err := suite.app.Test(req)
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode, tt.desc)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}

			var got ports.DTOResponseAllStock
			suite.Require().NoError(json.Unmarshal(body, &got))
			suite.Require().Equal(expect, got)
		})
	}
}

func (suite *test) TestByID() {
	stockID := uuid.NewString()
	stock := factory.NewStock(stockID, "Test Stock", 100)

	table := []struct {
		name   string
		input  string
		status int
		err    *ports.Error
		setup  func() ports.DTOResponseByIDStock
	}{
		{
			name:   "Success",
			input:  stockID,
			status: fiber.StatusOK,
			setup: func() ports.DTOResponseByIDStock {
				suite.service.EXPECT().ByID(gomock.Any(), stockID).Return(stock, nil).Times(1)
				return helpers.ToDTOResponseByIDStock(stock)
			},
		},
		{
			name:   "NotFound",
			input:  stockID,
			status: fiber.StatusNotFound,
			err:    ports.NewNotFoundError(nil),
			setup: func() ports.DTOResponseByIDStock {
				suite.service.EXPECT().ByID(gomock.Any(), stockID).Return(nil, ports.NewNotFoundError(nil)).Times(1)
				return ports.DTOResponseByIDStock{}
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expect := tt.setup()

			url := fmt.Sprintf("/api/v1/stocks/%s", tt.input)
			resp, err := suite.app.Test(httptest.NewRequest(fiber.MethodGet, url, nil))
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}

			var got ports.DTOResponseByIDStock
			suite.Require().NoError(json.Unmarshal(body, &got))
			suite.Require().Equal(expect, got)
		})
	}
}

func (suite *test) TestCreate() {
	stock := factory.NewStock(uuid.NewString(), "Test Stock", 100)
	input := helpers.ToDTORequestCreateStock(stock)

	table := []struct {
		name   string
		status int
		err    *ports.Error
		setup  func() ports.DTOResponseCreateStock
	}{
		{
			name:   "Success",
			status: fiber.StatusCreated,
			setup: func() ports.DTOResponseCreateStock {
				suite.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(stock, nil).Times(1)
				return helpers.ToDTOResponseCreateStock(stock)
			},
		},
		{
			name:   "InternalError",
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() ports.DTOResponseCreateStock {
				suite.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, ports.NewInternalError(nil)).Times(1)
				return ports.DTOResponseCreateStock{}
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expect := tt.setup()

			b, err := json.Marshal(input)
			suite.Require().NoError(err)
			req := httptest.NewRequest(fiber.MethodPost, "/api/v1/stocks", bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req)
			suite.Require().NoError(err)

			body, err := io.ReadAll(resp.Body)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}

			var got ports.DTOResponseCreateStock
			suite.Require().NoError(json.Unmarshal(body, &got))
			suite.Require().Equal(expect, got)
		})
	}
}

func (suite *test) TestUpdate() {
	stockID := uuid.NewString()
	stock := factory.NewStock(stockID, "Test Stock", 100)
	input := helpers.ToDTORequestUpdateStock(stock)

	table := []struct {
		name   string
		status int
		err    *ports.Error
		setup  func()
	}{
		{
			name:   "Success",
			status: fiber.StatusOK,
			setup: func() {
				suite.service.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:   "NotFound",
			status: fiber.StatusNotFound,
			err:    ports.NewNotFoundError(nil),
			setup: func() {
				suite.service.EXPECT().Update(gomock.Any(), gomock.Any()).Return(ports.NewNotFoundError(nil)).Times(1)
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()
			b, err := json.Marshal(input)
			suite.Require().NoError(err)

			url := fmt.Sprintf("/api/v1/stocks/%s", stockID)
			req := httptest.NewRequest(fiber.MethodPut, url, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req)
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				body, err := io.ReadAll(resp.Body)
				suite.Require().NoError(err)
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}
		})
	}
}

func (suite *test) TestDelete() {
	stockID := uuid.NewString()

	table := []struct {
		name   string
		status int
		err    *ports.Error
		setup  func()
	}{
		{
			name:   "Success",
			status: fiber.StatusOK,
			setup: func() {
				suite.service.EXPECT().Remove(gomock.Any(), stockID).Return(nil).Times(1)
			},
		},
		{
			name:   "InternalError",
			status: fiber.StatusInternalServerError,
			err:    ports.NewInternalError(nil),
			setup: func() {
				suite.service.EXPECT().Remove(gomock.Any(), stockID).Return(ports.NewInternalError(nil)).Times(1)
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()
			url := fmt.Sprintf("/api/v1/stocks/%s", stockID)
			resp, err := suite.app.Test(httptest.NewRequest(fiber.MethodDelete, url, nil))
			suite.Require().NoError(err)
			suite.Require().Equal(tt.status, resp.StatusCode)

			if tt.err != nil {
				body, err := io.ReadAll(resp.Body)
				suite.Require().NoError(err)
				suite.RequirePortError(suite.T(), body, tt.err)
				return
			}
		})
	}
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
