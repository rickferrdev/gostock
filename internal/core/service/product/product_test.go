package product

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

func (suite *test) TearDownSuite() {
	suite.ctrl.Finish()
}

func (suite *test) TestValidateProduct() {
	stockID, productID := uuid.NewString(), uuid.NewString()
	invalidProduct := factory.NewProduct("abc", "abc")
	validProduct := factory.NewProduct(productID, stockID)

	table := []struct {
		name         string
		desc         string
		inputProduct *domain.Product
		expectErr    error
	}{
		{
			name:         "Success",
			desc:         "should pass validation for a well-formed product",
			inputProduct: validProduct,
			expectErr:    nil,
		},
		{
			name:         "InvalidProduct",
			desc:         "should return 400 when product fields are not valid UUIDs",
			inputProduct: invalidProduct,
			expectErr:    ports.NewBadRequestError(nil),
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			err := suite.service.ValidateProduct(tt.inputProduct)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}

			suite.NoError(err, tt.desc)
		})
	}
}

func (suite *test) TestByStockID() {
	productID, stockID := uuid.NewString(), uuid.NewString()
	validProduct := factory.NewProductList(5, productID, stockID)

	table := []struct {
		name          string
		desc          string
		input         string
		mockProduct   []*domain.Product
		mockErr       error
		expectProduct []*domain.Product
		expectErr     error
		callsRepo     bool
	}{
		{
			name:          "Success",
			desc:          "should return all products for a valid stock ID",
			input:         stockID,
			mockProduct:   validProduct,
			expectProduct: validProduct,
			mockErr:       nil,
			expectErr:     nil,
			callsRepo:     true,
		},
		{
			name:      "InvalidStockID",
			desc:      "should return 400 when stock ID is not a valid UUID",
			input:     "ABC",
			mockErr:   ports.NewBadRequestError(nil),
			expectErr: ports.NewBadRequestError(nil),
			callsRepo: false,
		},
		{
			name:      "NotFound",
			desc:      "should return 404 when no products exist for the stock ID",
			input:     stockID,
			mockErr:   ports.NewNotFoundError(nil),
			expectErr: ports.NewNotFoundError(nil),
			callsRepo: true,
		},
		{
			name:      "InternalError",
			desc:      "should return 500 when repository fails unexpectedly",
			input:     stockID,
			mockErr:   ports.NewInternalError(nil),
			expectErr: ports.NewInternalError(nil),
			callsRepo: true,
		},
	}

	for _, tt := range table {
		tt := tt
		suite.Run(tt.name, func() {
			if tt.callsRepo {
				suite.product.EXPECT().
					ByStockID(gomock.Any(), tt.input).
					Times(1).
					Return(tt.mockProduct, tt.mockErr)
			}

			product, err := suite.service.ByStockID(context.Background(), tt.input)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}
			suite.NoError(err)
			suite.Equal(tt.expectProduct, product)
		})
	}
}

func (suite *test) TestByID() {
	productID, stockID := uuid.NewString(), uuid.NewString()
	validProduct := factory.NewProduct(productID, stockID)

	table := []struct {
		name           string
		desc           string
		inputID        string
		productMock    *domain.Product
		productMockErr error
		expectProduct  *domain.Product
		expectErr      error
		callsRepo      bool
	}{
		{
			name:           "Success",
			desc:           "should return the product for a valid ID",
			inputID:        productID,
			productMock:    validProduct,
			expectProduct:  validProduct,
			productMockErr: nil,
			expectErr:      nil,
			callsRepo:      true,
		},
		{
			name:           "InvalidProductID",
			desc:           "should return 400 when product ID is not a valid UUID",
			inputID:        "ABC",
			productMockErr: ports.NewBadRequestError(nil),
			expectErr:      ports.NewBadRequestError(nil),
			callsRepo:      false,
		},
		{
			name:           "NotFound",
			desc:           "should return 404 when product does not exist",
			inputID:        productID,
			productMockErr: ports.NewNotFoundError(nil),
			expectErr:      ports.NewNotFoundError(nil),
			callsRepo:      true,
		},
		{
			name:           "InternalError",
			desc:           "should return 500 when repository fails unexpectedly",
			inputID:        productID,
			productMockErr: ports.NewInternalError(nil),
			expectErr:      ports.NewInternalError(nil),
			callsRepo:      true,
		},
	}

	for _, tt := range table {
		tt := tt
		suite.Run(tt.name, func() {
			if tt.callsRepo {
				suite.product.EXPECT().
					ByID(gomock.Any(), tt.inputID).
					Times(1).
					Return(tt.productMock, tt.productMockErr)
			}

			product, err := suite.service.ByID(context.Background(), tt.inputID)

			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}
			suite.NoError(err)
			suite.Equal(tt.expectProduct, product)
		})
	}
}

func (suite *test) TestCreate() {
	validProduct := factory.NewProduct("", uuid.NewString())
	invalidProduct := factory.NewProduct("abc", uuid.NewString())

	table := []struct {
		name           string
		desc           string
		input          *domain.Product
		productMock    *domain.Product
		productMockErr error
		expectProduct  *domain.Product
		expectErr      error
		callsRepo      bool
	}{
		{
			name:          "Success",
			desc:          "should create the product and return it",
			input:         validProduct,
			productMock:   validProduct,
			expectProduct: validProduct,
			callsRepo:     true,
		},
		{
			name:      "InvalidProductID",
			desc:      "should return 409 when product ID already exists",
			input:     invalidProduct,
			expectErr: ports.NewConflictError(nil),
			callsRepo: false,
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			if tt.callsRepo {
				suite.product.EXPECT().
					CreateAtomic(gomock.Any(), gomock.Any()).
					Times(1).
					Return(tt.productMockErr)
			}

			_, err := suite.service.Create(context.Background(), tt.input)

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
		input     *domain.Product
		expectErr error
		setup     func() *domain.Product
	}{
		{
			name: "Success",
			desc: "should update the product and return nil",
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID)
				oldProduct := factory.NewProduct(productID, stockID)
				stock := &domain.Stock{ID: stockID, Capacity: 100}

				suite.product.EXPECT().
					ByID(gomock.Any(), productID).
					Times(1).
					Return(oldProduct, nil)

				suite.stock.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)

				suite.product.EXPECT().
					UpdateAtomic(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				return input
			},
		},
		{
			name:      "InvalidProduct",
			desc:      "should return 400 when product fields are invalid",
			expectErr: ports.NewBadRequestError(nil),
			setup: func() *domain.Product {
				return &domain.Product{}
			},
		},
		{
			name:      "ProductNotFound",
			desc:      "should return 404 when product does not exist",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID)

				suite.product.EXPECT().
					ByID(gomock.Any(), productID).
					Times(1).
					Return(nil, ports.NewNotFoundError(nil))

				return input
			},
		},
		{
			name:      "StockNotFound",
			desc:      "should return 404 when stock of the product does not exist",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID)
				oldProduct := factory.NewProduct(productID, stockID)

				suite.product.EXPECT().
					ByID(gomock.Any(), productID).
					Times(1).
					Return(oldProduct, nil)

				suite.stock.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(nil, ports.NewNotFoundError(nil))

				return input
			},
		},
		{
			name:      "CapacityExceeded",
			desc:      "should return capacity exceeded when stock cannot fit the new quantity",
			expectErr: ports.NewCapacityExceeded(nil),
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID)
				oldProduct := &domain.Product{ID: productID, StockID: stockID, Qtd: 21}
				stock := &domain.Stock{ID: stockID, Capacity: 20}

				suite.product.EXPECT().
					ByID(gomock.Any(), productID).
					Times(1).
					Return(oldProduct, nil)

				suite.stock.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)

				return input
			},
		},
		{
			name:      "InternalError",
			desc:      "should return 500 when repository fails on update",
			expectErr: ports.NewInternalError(nil),
			setup: func() *domain.Product {
				productID, stockID := uuid.NewString(), uuid.NewString()
				input := factory.NewProduct(productID, stockID)
				oldProduct := factory.NewProduct(productID, stockID)
				stock := &domain.Stock{ID: stockID, Capacity: 100}

				suite.product.EXPECT().
					ByID(gomock.Any(), productID).
					Times(1).
					Return(oldProduct, nil)

				suite.stock.EXPECT().
					ByID(gomock.Any(), stockID).
					Times(1).
					Return(stock, nil)

				suite.product.EXPECT().
					UpdateAtomic(gomock.Any(), gomock.Any()).
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

func (suite *test) TestDelete() {
	table := []struct {
		name      string
		desc      string
		setup     func() string
		expectErr error
	}{
		{
			name: "Success",
			desc: "",
			setup: func() string {
				productID := uuid.NewString()

				suite.product.EXPECT().
					Delete(gomock.Any(), productID).
					Times(1).
					Return(nil)

				return productID
			},
		},
		{
			name:      "InvalidID",
			expectErr: ports.NewBadRequestError(nil),
			setup: func() string {
				productID := "abc"

				suite.product.EXPECT().
					Delete(gomock.Any(), productID).
					Times(0).
					Return(nil)

				return productID
			},
		},
		{
			name:      "ProductNotFound",
			expectErr: ports.NewNotFoundError(nil),
			setup: func() string {
				productID := uuid.NewString()
				suite.product.EXPECT().
					Delete(gomock.Any(), productID).
					Times(1).
					Return(ports.NewNotFoundError(nil))

				return productID
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			input := tt.setup()

			err := suite.service.Delete(context.Background(), input)
			if tt.expectErr != nil {
				suite.ErrorAs(err, &tt.expectErr, tt.desc)
				return
			}

			suite.NoError(err, tt.desc)
		})
	}
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
