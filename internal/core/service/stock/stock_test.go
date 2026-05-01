package stock

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/domain/factory"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/rickferrdev/gostock/internal/tests/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type test struct {
	suite.Suite
	service ports.StockService
	ctrl    *gomock.Controller
	stock   *mocks.MockStockRepository
}

func (suite *test) SetupTest() {
	var errService error
	suite.ctrl = gomock.NewController(suite.T())
	suite.stock = mocks.NewMockStockRepository(suite.ctrl)
	suite.service, errService = New(ServiceParams{Stock: suite.stock})
	suite.NoError(errService)
}

func (suite *test) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *test) RequirePortError(err error, expected *ports.Error) {
	var got *ports.Error

	suite.Require().ErrorAs(err, &got)
	suite.Require().Equal(expected.Code, got.Code)
	suite.Require().Equal(expected.Status, got.Status)
	suite.Require().Equal(expected.Message, got.Message)
}

func (suite *test) TestValidateStock() {
	table := []struct {
		name   string
		desc   string
		expect error
		setup  func() *domain.Stock
	}{
		{
			name: "Success",
			desc: "should pass validation for a well-formed stock",
			setup: func() *domain.Stock {
				return factory.NewStock(uuid.NewString(), "Stock", 100)
			},
		},
		{
			name:   "InvalidProduct",
			desc:   "should return 400 when product fields are not valid UUIDs",
			expect: ports.NewBadRequestError(nil),
			setup: func() *domain.Stock {
				return &domain.Stock{ID: "abc", Name: "Stock", Capacity: 0}
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()
			err := suite.service.ValidateStock(input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}

			suite.Require().NoError(err, tt.desc)
		})
	}
}

func (suite *test) TestAll() {
	stocks := factory.NewStockSlice(5, uuid.NewString())
	table := []struct {
		name        string
		desc        string
		setup       func()
		expectStock []*domain.Stock
		expect      error
	}{
		{
			name:        "Success",
			desc:        "should return all stocks successfully",
			expectStock: stocks,
			setup: func() {
				suite.stock.EXPECT().All(gomock.Any()).Times(1).Return(stocks, nil)
			},
		},
		{
			name:   "StocksNotFound",
			desc:   "should return 404 when no stocks are found",
			expect: ports.NewNotFoundError(nil),
			setup: func() {
				suite.stock.EXPECT().All(gomock.Any()).Times(1).Return(nil, sql.ErrNoRows)
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails unexpectedly",
			expect: ports.NewInternalError(nil),
			setup: func() {
				suite.stock.EXPECT().All(gomock.Any()).Times(1).Return(nil, ports.NewInternalError(nil))
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()

			actual, err := suite.service.All(context.Background())
			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}

			suite.Require().NoError(err)
			suite.Require().Equal(stocks, actual)
		})
	}
}

func (suite *test) TestByID() {
	stockID := uuid.NewString()

	table := []struct {
		name   string
		desc   string
		input  string
		setup  func() *domain.Stock
		expect error
	}{
		{
			name:  "Success",
			desc:  "should return the stock for a valid ID",
			input: stockID,
			setup: func() *domain.Stock {
				stock := factory.NewStock(stockID, "Stock", 100)
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(stock, nil)
				return stock
			},
		},
		{
			name:   "InvalidID",
			desc:   "should return 400 when stock ID is not a valid UUID",
			expect: ports.NewBadRequestError(nil),
			input:  "abc",
			setup: func() *domain.Stock {
				return nil
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 when stock does not exist",
			input:  stockID,
			expect: ports.NewNotFoundError(nil),
			setup: func() *domain.Stock {
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(nil, sql.ErrNoRows)
				return nil
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails unexpectedly",
			input:  stockID,
			expect: ports.NewInternalError(nil),
			setup: func() *domain.Stock {
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(nil, ports.NewInternalError(nil))
				return nil
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expect := tt.setup()
			actual, err := suite.service.ByID(context.Background(), tt.input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}

			suite.Require().NoError(err)
			suite.Require().Equal(expect, actual)
		})
	}
}

func (suite *test) TestCreate() {
	table := []struct {
		name   string
		desc   string
		setup  func() *domain.Stock
		expect error
	}{
		{
			name: "Success",
			desc: "should create the stock and return it",
			setup: func() *domain.Stock {
				stock := factory.NewStock("", "Stock", 100)
				suite.stock.EXPECT().CreateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return stock
			},
		},
		{
			name:   "InvalidStock",
			desc:   "should return 400 when stock fields are invalid",
			expect: ports.NewBadRequestError(nil),
			setup: func() *domain.Stock {
				return &domain.Stock{ID: "abc"}
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails on create",
			expect: ports.NewInternalError(nil),
			setup: func() *domain.Stock {
				stock := factory.NewStock("", "Stock", 100)
				suite.stock.EXPECT().CreateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(ports.NewInternalError(nil))
				return stock
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			actual, err := suite.service.Create(context.Background(), input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}
			suite.Require().NoError(err)
			suite.Require().NotNil(input)
			suite.Require().NotNil(actual)
		})
	}
}

func (suite *test) TestOccupancy() {
	stockID := uuid.NewString()
	productID := uuid.NewString()

	table := []struct {
		name   string
		desc   string
		input  string
		setup  func() uint64
		expect error
	}{
		{
			name:  "Success",
			desc:  "should return occupancy for a valid stock ID",
			input: stockID,
			setup: func() uint64 {
				stock := &domain.Stock{ID: stockID, Name: "Stock", Capacity: 100, Items: factory.NewProductSlice(5, productID, stockID)}
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(stock, nil)

				return 5
			},
		},
		{
			name:   "InvalidID",
			desc:   "should return 400 when stock ID is not a valid UUID",
			input:  "abc",
			expect: ports.NewBadRequestError(nil),
			setup: func() uint64 {
				return 0
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 when stock does not exist",
			expect: ports.NewNotFoundError(nil),
			input:  stockID,
			setup: func() uint64 {
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(nil, sql.ErrNoRows)

				return 0
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails unexpectedly",
			expect: ports.NewInternalError(nil),
			input:  stockID,
			setup: func() uint64 {
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(nil, ports.NewInternalError(nil))

				return 0
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expect := tt.setup()

			actual, err := suite.service.Occupancy(context.Background(), tt.input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}

			suite.Require().NoError(err)
			suite.Require().Equal(expect, actual)
		})
	}
}

func (suite *test) TestAvailableSpace() {
	stockID := uuid.NewString()
	table := []struct {
		name   string
		desc   string
		setup  func() uint64
		input  string
		expect error
	}{
		{
			name:  "Success",
			desc:  "should return available space for a valid stock ID",
			input: stockID,
			setup: func() uint64 {
				stock := factory.NewStock(stockID, "Stock", 100)
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(stock, nil)
				return 100
			},
		},
		{
			name:   "InvalidID",
			desc:   "should return 400 when stock ID is not a valid UUID",
			input:  "abc",
			expect: ports.NewBadRequestError(nil),
			setup: func() uint64 {
				return 0
			},
		},
		{
			name:   "StockNotFound",
			desc:   "should return 404 when stock does not exist",
			input:  stockID,
			expect: ports.NewNotFoundError(nil),
			setup: func() uint64 {
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(nil, sql.ErrNoRows)

				return 0
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails on occupancy check",
			input:  stockID,
			expect: ports.NewInternalError(nil),
			setup: func() uint64 {
				stock := &domain.Stock{ID: stockID, Capacity: 100}
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(stock, ports.NewInternalError(nil))

				return 100
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expect := tt.setup()

			actual, err := suite.service.AvailableSpace(context.Background(), tt.input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}

			suite.Require().NoError(err)
			suite.Require().Equal(expect, actual)
		})
	}
}

func (suite *test) TestUpdate() {
	old := factory.NewStock(uuid.NewString(), "Old Stock", 50)
	input := factory.NewStock(uuid.NewString(), "Updated Stock", 100)

	table := []struct {
		name   string
		desc   string
		setup  func() *domain.Stock
		expect error
	}{
		{
			name: "Success",
			desc: "should update the stock and return nil",
			setup: func() *domain.Stock {
				suite.stock.EXPECT().ByID(gomock.Any(), input.ID).Times(1).Return(old, nil)
				suite.stock.EXPECT().UpdateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(nil)

				return input
			},
		},
		{
			name:   "InvalidStock",
			desc:   "should return 400 when stock fields are invalid",
			expect: ports.NewBadRequestError(nil),
			setup: func() *domain.Stock {
				return &domain.Stock{ID: "abc"}
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 when stock does not exist",
			expect: ports.NewNotFoundError(nil),
			setup: func() *domain.Stock {
				suite.stock.EXPECT().ByID(gomock.Any(), input.ID).Times(1).Return(nil, sql.ErrNoRows)

				return input
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails on update",
			expect: ports.NewInternalError(nil),
			setup: func() *domain.Stock {
				suite.stock.EXPECT().ByID(gomock.Any(), input.ID).Times(1).Return(input, nil)
				suite.stock.EXPECT().UpdateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(ports.NewInternalError(nil))

				return input
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			err := suite.service.Update(context.Background(), input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}
			suite.NoError(err, tt.desc)
		})
	}
}

func (suite *test) TestRemove() {
	table := []struct {
		name   string
		desc   string
		setup  func() string
		expect error
	}{
		{
			name: "Success",
			desc: "should delete the stock and return nil",
			setup: func() string {
				stockID := uuid.NewString()
				stock := factory.NewStock(stockID, "Stock", 100)

				suite.stock.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)

				suite.stock.EXPECT().
					DeleteAtomic(gomock.Any(), stockID).
					Times(1).
					Return(nil)
				return stockID
			},
		},
		{
			name:   "InvalidID",
			desc:   "should return 400 when stock ID is not a valid UUID",
			expect: ports.NewBadRequestError(nil),
			setup: func() string {
				return "abc"
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 when stock does not exist",
			expect: ports.NewNotFoundError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(nil, sql.ErrNoRows)
				return stockID
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails on remove",
			expect: ports.NewInternalError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				stock := factory.NewStock(stockID, "Stock", 100)

				suite.stock.EXPECT().ByID(gomock.Any(), stockID).Times(1).Return(stock, nil)

				suite.stock.EXPECT().DeleteAtomic(gomock.Any(), stockID).Times(1).Return(ports.NewInternalError(nil))
				return stockID
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			err := suite.service.Remove(context.Background(), input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}
			suite.Require().NoError(err, tt.desc)
		})
	}
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
