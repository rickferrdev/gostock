package stock

import (
	"context"
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
	repo    *mocks.MockStockRepository
}

func (suite *test) SetupTest() {
	var errService error
	suite.ctrl = gomock.NewController(suite.T())
	suite.repo = mocks.NewMockStockRepository(suite.ctrl)
	suite.service, errService = New(ServiceParams{Stock: suite.repo})
	suite.NoError(errService)
}

func (suite *test) TearDownSuite() {
	suite.ctrl.Finish()
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
				return factory.NewStock(uuid.NewString())
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
				suite.ErrorAs(err, &tt.expect, tt.desc)
				return
			}

			suite.NoError(err, tt.desc)
		})
	}
}

func (suite *test) TestAll() {
	stocks := factory.NewStockList(5, uuid.NewString())
	table := []struct {
		name        string
		desc        string
		setup       func()
		expectStock []*domain.Stock
		expectErr   error
	}{
		{
			name:        "Success",
			desc:        "",
			expectStock: stocks,
			setup: func() {
				suite.repo.EXPECT().All(gomock.Any()).Times(1).Return(stocks, nil)
			},
		},
		{
			name:      "StocksNotFound",
			desc:      "",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() {
				suite.repo.EXPECT().All(gomock.Any()).Times(1).Return(nil, ports.NewNotFoundError(nil))
			},
		},
		{
			name:      "InternalError",
			desc:      "",
			expectErr: ports.NewInternalError(nil),
			setup: func() {
				suite.repo.EXPECT().All(gomock.Any()).Times(1).Return(nil, ports.NewInternalError(nil))
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()

			actual, err := suite.service.All(context.Background())
			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}

			suite.NoError(err)
			suite.Equal(stocks, actual)
		})
	}
}

func (suite *test) TestByID() {
	table := []struct {
		name        string
		desc        string
		setup       func() string
		expectStock *domain.Stock
		expectErr   error
	}{
		{
			name:        "Success",
			desc:        "should return the stock for a valid ID",
			expectStock: factory.NewStock(uuid.NewString()),
			setup: func() string {
				stockID := uuid.NewString()
				stock := factory.NewStock(stockID)
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)
				return stockID
			},
		},
		{
			name:      "InvalidID",
			desc:      "should return 400 when stock ID is not a valid UUID",
			expectErr: ports.NewBadRequestError(nil),
			setup: func() string {
				return "abc"
			},
		},
		{
			name:      "NotFound",
			desc:      "should return 404 when stock does not exist",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewNotFoundError(nil))
				return stockID
			},
		},
		{
			name:      "InternalError",
			desc:      "should return 500 when repository fails unexpectedly",
			expectErr: ports.NewInternalError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewInternalError(nil))
				return stockID
			},
		},
	}

	for _, tt := range table {
		tt := tt
		suite.Run(tt.name, func() {
			input := tt.setup()

			actual, err := suite.service.ByID(context.Background(), input)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}
			suite.NoError(err)
			suite.NotNil(actual)
		})
	}
}

func (suite *test) TestCreate() {
	table := []struct {
		name      string
		desc      string
		setup     func() *domain.Stock
		expectErr error
	}{
		{
			name: "Success",
			desc: "should create the stock and return it",
			setup: func() *domain.Stock {
				stock := factory.NewStock(uuid.NewString())
				suite.repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
				return stock
			},
		},
		{
			name:      "InvalidStock",
			desc:      "should return 400 when stock fields are invalid",
			expectErr: ports.NewBadRequestError(nil),
			setup: func() *domain.Stock {
				return &domain.Stock{ID: "abc"}
			},
		},
		{
			name:      "InternalError",
			desc:      "should return 500 when repository fails on create",
			expectErr: ports.NewInternalError(nil),
			setup: func() *domain.Stock {
				stock := factory.NewStock(uuid.NewString())
				suite.repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(ports.NewInternalError(nil))
				return stock
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			_, err := suite.service.Create(context.Background(), input)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}
			suite.NoError(err)
		})
	}
}

func (suite *test) TestOccupancy() {
	table := []struct {
		name            string
		desc            string
		setup           func() string
		expectErr       error
		expectNonNil    bool
		expectOccupancy uint64
	}{
		{
			name:            "Success",
			desc:            "should return occupancy for a valid stock ID",
			expectOccupancy: 500,
			expectNonNil:    true,
			setup: func() string {
				stockID := uuid.NewString()
				productID := uuid.NewString()
				stock := &domain.Stock{ID: stockID, Name: "Stock", Capacity: 100, Items: factory.NewProductList(5, productID, stockID)}
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)
				return stockID
			},
		},
		{
			name:      "InvalidID",
			desc:      "should return 400 when stock ID is not a valid UUID",
			expectErr: ports.NewBadRequestError(nil),
			setup: func() string {
				return "abc"
			},
		},
		{
			name:      "NotFound",
			desc:      "should return 404 when stock does not exist",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewNotFoundError(nil))
				return stockID
			},
		},
		{
			name:      "InternalError",
			desc:      "should return 500 when repository fails unexpectedly",
			expectErr: ports.NewInternalError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewInternalError(nil))
				return stockID
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			actual, err := suite.service.Occupancy(context.Background(), input)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}
			suite.NoError(err)
			suite.NotNil(actual)
			suite.Equal(tt.expectOccupancy, actual)
		})
	}
}

func (suite *test) TestAvailableSpace() {
	table := []struct {
		name      string
		desc      string
		setup     func() string
		expectErr error
	}{
		{
			name: "Success",
			desc: "should return available space for a valid stock ID",
			setup: func() string {
				stockID := uuid.NewString()
				stock := factory.NewStock(stockID)
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)
				return stockID
			},
		},
		{
			name:      "InvalidID",
			desc:      "should return 400 when stock ID is not a valid UUID",
			expectErr: ports.NewBadRequestError(nil),
			setup: func() string {
				return "abc"
			},
		},
		{
			name:      "StockNotFound",
			desc:      "should return 404 when stock does not exist",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewNotFoundError(nil))
				return stockID
			},
		},
		{
			name:      "InternalError",
			desc:      "should return 500 when repository fails on occupancy check",
			expectErr: ports.NewInternalError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				stock := &domain.Stock{ID: stockID, Capacity: 100}
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, ports.NewInternalError(nil))
				return stockID
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			_, err := suite.service.AvailableSpace(context.Background(), input)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}
			suite.NoError(err)
		})
	}
}

func (suite *test) TestUpdate() {
	table := []struct {
		name      string
		desc      string
		setup     func() *domain.Stock
		expectErr error
	}{
		{
			name: "Success",
			desc: "should update the stock and return nil",
			setup: func() *domain.Stock {
				stockID := uuid.NewString()
				input := factory.NewStock(stockID)
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(input, nil)
				suite.repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
				return input
			},
		},
		{
			name:      "InvalidStock",
			desc:      "should return 400 when stock fields are invalid",
			expectErr: ports.NewBadRequestError(nil),
			setup: func() *domain.Stock {
				return &domain.Stock{ID: "abc"}
			},
		},
		{
			name:      "NotFound",
			desc:      "should return 404 when stock does not exist",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() *domain.Stock {
				stockID := uuid.NewString()
				input := factory.NewStock(stockID)
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewNotFoundError(nil))
				return input
			},
		},
		{
			name:      "InternalError",
			desc:      "should return 500 when repository fails on update",
			expectErr: ports.NewInternalError(nil),
			setup: func() *domain.Stock {
				stockID := uuid.NewString()
				input := factory.NewStock(stockID)
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(input, nil)
				suite.repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(1).
					Return(ports.NewInternalError(nil))
				return input
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			err := suite.service.Update(context.Background(), input)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
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
				stock := factory.NewStock(stockID)

				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)

				suite.repo.EXPECT().
					Delete(gomock.Any(), stockID).
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
				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewNotFoundError(nil))
				return stockID
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails on remove",
			expect: ports.NewInternalError(nil),
			setup: func() string {
				stockID := uuid.NewString()
				stock := factory.NewStock(stockID)

				suite.repo.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)

				suite.repo.EXPECT().
					Delete(gomock.Any(), stockID).
					Times(1).
					Return(ports.NewInternalError(nil))
				return stockID
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			err := suite.service.Remove(context.Background(), input)

			if tt.expect != nil {
				suite.ErrorAs(err, &tt.expect, tt.desc)
				return
			}
			suite.NoError(err, tt.desc)
		})
	}
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
