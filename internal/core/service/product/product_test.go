package product

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
	product *mocks.MockProductRepository
	stock   *mocks.MockStockRepository
	ctrl    *gomock.Controller
	service ports.ProductService
}

func (suite *test) SetupTest() {
	var errService error
	suite.ctrl = gomock.NewController(suite.T())
	suite.product = mocks.NewMockProductRepository(suite.ctrl)
	suite.stock = mocks.NewMockStockRepository(suite.ctrl)
	suite.service, errService = New(ServiceParams{Product: suite.product, Stock: suite.stock})
	suite.NoError(errService)
}

func (suite *test) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *test) RequirePortError(err error, expected *ports.Error) {
	var got *ports.Error

	suite.Require().ErrorAs(err, &got)
	suite.Require().Equal(expected.Code, got.Code)
	suite.Require().Equal(expected.Status, got.Status)
	suite.Require().Equal(expected.Message, got.Message)
}

func (suite *test) TestValidateProduct() {
	stockID, productID := uuid.NewString(), uuid.NewString()

	table := []struct {
		name   string
		desc   string
		expect error
		setup  func() *domain.Product
	}{
		{
			name: "Success",
			desc: "should pass validation for a well-formed product",
			setup: func() *domain.Product {
				return factory.NewProduct(productID, stockID, "Product", 10, 50)
			},
			expect: nil,
		},
		{
			name: "InvalidProduct",
			desc: "should return 400 when product fields are not valid UUIDs",
			setup: func() *domain.Product {
				return factory.NewProduct("", stockID, "Product", 10, 50)
			},
			expect: ports.NewBadRequestError(nil),
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			product := tt.setup()
			err := suite.service.ValidateProduct(product)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}

			suite.Require().NoError(err, tt.desc)
		})
	}
}

func (suite *test) TestByStockID() {
	productID, stockID := uuid.NewString(), uuid.NewString()

	table := []struct {
		name   string
		desc   string
		input  string
		setup  func() []*domain.Product
		expect error
	}{
		{
			name:  "Success",
			desc:  "should return all products for a valid stock ID",
			input: stockID,
			setup: func() []*domain.Product {
				products := factory.NewProductSlice(5, productID, stockID)

				suite.product.EXPECT().
					ByStockID(gomock.Any(), stockID).Times(1).Return(products, nil)
				return products
			},
			expect: nil,
		},
		{
			name:   "InvalidStockID",
			desc:   "should return 400 when stock ID is not a valid UUID",
			input:  "ABC",
			expect: ports.NewBadRequestError(nil),
			setup: func() []*domain.Product {
				suite.product.EXPECT().ByStockID(gomock.Any(), "ABC").Times(0).Return(nil, ports.NewBadRequestError(nil))
				return nil
			},
		},
		{
			name:   "NotFound",
			desc:   "should return 404 when no products exist for the stock ID",
			input:  stockID,
			expect: ports.NewNotFoundError(nil),
			setup: func() []*domain.Product {
				suite.product.EXPECT().ByStockID(gomock.Any(), stockID).Times(1).Return(nil, sql.ErrNoRows)
				return nil
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails unexpectedly",
			input:  stockID,
			expect: ports.NewInternalError(nil),
			setup: func() []*domain.Product {
				suite.product.EXPECT().ByStockID(gomock.Any(), stockID).Times(1).Return(nil, ports.NewInternalError(nil))
				return nil
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			expectProducts := tt.setup()

			find, err := suite.service.ByStockID(context.Background(), tt.input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}
			suite.Require().NoError(err, tt.desc)
			suite.Require().Equal(expectProducts, find)
		})
	}
}

func (suite *test) TestByID() {
	productID, stockID := uuid.NewString(), uuid.NewString()

	table := []struct {
		name   string
		desc   string
		input  string
		setup  func() *domain.Product
		expect error
	}{
		{
			name:  "Success",
			desc:  "should return the product for a valid ID",
			input: productID,
			setup: func() *domain.Product {
				product := factory.NewProduct(productID, stockID, "Product", 10, 50)
				suite.product.EXPECT().ByID(gomock.Any(), productID).Times(1).Return(product, nil)
				return product
			},
			expect: nil,
		},
		{
			name:  "InvalidProductID",
			desc:  "should return 400 when product ID is not a valid UUID",
			input: "ABC",
			setup: func() *domain.Product {
				return nil
			},
			expect: ports.NewBadRequestError(nil),
		},
		{
			name:  "NotFound",
			desc:  "should return 404 when product does not exist",
			input: productID,
			setup: func() *domain.Product {
				suite.product.EXPECT().ByID(gomock.Any(), productID).Times(1).Return(nil, sql.ErrNoRows)
				return nil
			},
			expect: ports.NewNotFoundError(nil),
		},
		{
			name:  "InternalError",
			desc:  "should return 500 when repository fails unexpectedly",
			input: productID,
			setup: func() *domain.Product {
				suite.product.EXPECT().ByID(gomock.Any(), productID).Times(1).Return(nil, ports.NewInternalError(nil))
				return nil
			},
			expect: ports.NewInternalError(nil),
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			product := tt.setup()

			find, err := suite.service.ByID(context.Background(), tt.input)

			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}
			suite.Require().NoError(err, tt.desc)
			suite.Require().Equal(product, find)
		})
	}
}

func (suite *test) TestCreate() {
	table := []struct {
		name   string
		desc   string
		expect error
		setup  func() *domain.Product
	}{
		{
			name: "Success",
			desc: "should create the product and return nil",
			setup: func() *domain.Product {
				input := factory.NewProduct("", uuid.NewString(), "Product", 10, 100)
				suite.product.EXPECT().CreateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				return input
			},
		},
		{
			name:   "InvalidStockID",
			desc:   "should return 400 when stock ID is not a valid UUID",
			expect: ports.NewBadRequestError(nil),
			setup: func() *domain.Product {
				return factory.NewProduct("", "ABC", "Product", 10, 100)
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails unexpectedly",
			expect: ports.NewInternalError(nil),
			setup: func() *domain.Product {
				suite.product.EXPECT().CreateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(ports.NewInternalError(nil))
				return factory.NewProduct("", uuid.NewString(), "Product", 10, 100)
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			created, err := suite.service.Create(context.Background(), input)
			if tt.expect != nil {
				suite.RequirePortError(err, tt.expect.(*ports.Error))
				return
			}

			suite.Require().NoError(err, tt.desc)
			suite.Require().NotNil(created)
		})
	}
}

func (suite *test) TestUpdate() {
	table := []struct {
		name   string
		desc   string
		expect error
		setup  func() *domain.Product
	}{
		{
			name: "Success",
			desc: "should update the product and return nil",
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID, "Updated Product", 15, 150)
				old := factory.NewProduct(productID, stockID, "Product", 10, 100)
				stock := factory.NewStock(stockID, "Stock", 100)
				stock.Push(old)

				suite.product.EXPECT().ByID(gomock.Any(), input.ID).Times(1).Return(old, nil)
				suite.stock.EXPECT().ByID(gomock.Any(), input.StockID).Times(1).Return(stock, nil)
				suite.product.EXPECT().UpdateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(nil)

				return input
			},
		},
		{
			name:   "InvalidProduct",
			desc:   "should return 400 when product fields are invalid",
			expect: ports.NewBadRequestError(nil),
			setup: func() *domain.Product {
				return factory.NewProduct("ABC", "DEF", "Product", 10, 100)
			},
		},
		{
			name:   "ProductNotFound",
			desc:   "should return 404 when product does not exist",
			expect: ports.NewNotFoundError(nil),
			setup: func() *domain.Product {
				input := factory.NewProduct(uuid.NewString(), uuid.NewString(), "Product", 10, 100)

				suite.product.EXPECT().ByID(gomock.Any(), input.ID).Times(1).Return(nil, sql.ErrNoRows)
				return input
			},
		},
		{
			name:   "CapacityExceeded",
			desc:   "should return capacity exceeded when stock cannot fit the new quantity",
			expect: ports.NewCapacityExceeded(nil),
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID, "Product", 30, 100)
				old := factory.NewProduct(productID, stockID, "Product", 1, 100)
				stock := factory.NewStock(stockID, "Stock", 10)
				stock.Push(old)

				suite.product.EXPECT().ByID(gomock.Any(), input.ID).Times(1).Return(old, nil)
				suite.stock.EXPECT().ByID(gomock.Any(), input.StockID).Times(1).Return(stock, nil)

				return input
			},
		},
		{
			name:   "InternalError",
			desc:   "should return 500 when repository fails on update",
			expect: ports.NewInternalError(nil),
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID, "Product", 30, 100)
				old := factory.NewProduct(productID, stockID, "Product", 10, 100)
				stock := factory.NewStock(stockID, "Stock", 100)
				stock.Push(old)

				suite.product.EXPECT().ByID(gomock.Any(), input.ID).Times(1).Return(old, nil)
				suite.stock.EXPECT().ByID(gomock.Any(), input.StockID).Times(1).Return(stock, nil)
				suite.product.EXPECT().UpdateAtomic(gomock.Any(), gomock.Any()).Times(1).Return(ports.NewInternalError(nil))
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

			suite.Require().NoError(err, tt.desc)
		})
	}
}

func (suite *test) TestDelete() {
	table := []struct {
		name   string
		desc   string
		setup  func() string
		expect error
	}{
		{
			name: "Success",
			desc: "should delete the product and return nil",
			setup: func() string {
				productID := uuid.NewString()
				suite.product.EXPECT().DeleteAtomic(gomock.Any(), productID).Times(1).Return(nil)
				return productID
			},
		},
		{
			name:   "InvalidID",
			desc:   "should return 400 when product ID is not a valid UUID",
			expect: ports.NewBadRequestError(nil),
			setup: func() string {
				productID := "abc"
				suite.product.EXPECT().DeleteAtomic(gomock.Any(), productID).Times(0).Return(nil)

				return productID
			},
		},
		{
			name:   "ProductNotFound",
			desc:   "should return 404 when product does not exist",
			expect: ports.NewNotFoundError(nil),
			setup: func() string {
				productID := uuid.NewString()
				suite.product.EXPECT().DeleteAtomic(gomock.Any(), productID).Times(1).Return(sql.ErrNoRows)
				return productID
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			err := suite.service.Delete(context.Background(), input)
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
