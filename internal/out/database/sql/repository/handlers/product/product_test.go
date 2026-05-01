package product

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain/factory"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository/handlers/stock"
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository/schema"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type test struct {
	suite.Suite
	db      *bun.DB
	product ports.ProductRepository
	stock   ports.StockRepository
}

func (suite *test) SetupTest() {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file:gostock_product_test?mode=memory&cache=shared")
	suite.Require().NoError(err)
	sqldb.SetMaxOpenConns(1)

	suite.db = bun.NewDB(sqldb, sqlitedialect.New())
	suite.db.RegisterModel((*schema.Stock)(nil), (*schema.Product)(nil))

	_, err = suite.db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON")
	suite.Require().NoError(err)

	_, err = suite.db.NewCreateTable().
		Model((*schema.Stock)(nil)).
		IfNotExists().
		Exec(context.Background())
	suite.Require().NoError(err)

	_, err = suite.db.NewCreateTable().
		Model((*schema.Product)(nil)).
		WithForeignKeys().
		IfNotExists().
		Exec(context.Background())
	suite.Require().NoError(err)

	suite.product, err = New(suite.db)
	suite.Require().NoError(err)
	suite.stock, err = stock.New(suite.db)
	suite.Require().NoError(err)
}

func (suite *test) TearDownTest() {
	suite.Require().NoError(suite.db.Close())
}

func (suite *test) TestByStockIDWithProductNotFound() {
	products, err := suite.product.ByStockID(context.Background(), uuid.NewString())
	suite.Require().NoError(err)
	suite.Require().Empty(products)
}

func (suite *test) TestByStockIDWithProductFound() {
	stock, product := factory.NewStockWithOneProducts(5, 10)

	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))

	products, err := suite.product.ByStockID(context.Background(), stock.ID)
	suite.Require().NoError(err)
	suite.Require().Len(products, 1)
	suite.Equal(product, products[0])
}

func (suite *test) TestByIDWithProductNotFound() {
	product, err := suite.product.ByID(context.Background(), uuid.NewString())
	suite.Require().Error(err)
	suite.Require().Nil(product)
}

func (suite *test) TestByIDWithProductFound() {
	stock, product := factory.NewStockWithOneProducts(5, 10)
	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))

	find, err := suite.product.ByID(context.Background(), product.ID)
	suite.Require().NoError(err)
	suite.Equal(product, find)
}

func (suite *test) TestCreateAtomic() {
	stock := factory.NewStockDefault()
	product := factory.NewProduct(uuid.NewString(), stock.ID, "Product", 5, 100)
	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))

	findStock, err := suite.stock.ByID(context.Background(), stock.ID)
	suite.Require().NoError(err)
	suite.Equal(stock, findStock)

	suite.Require().NoError(suite.product.CreateAtomic(context.Background(), product))
	findProduct, err := suite.product.ByID(context.Background(), product.ID)

	suite.Require().NoError(err)
	suite.Equal(product, findProduct)
}

func (suite *test) TestCreateAtomicWithStockNotFound() {
	product := factory.NewProduct(uuid.NewString(), uuid.NewString(), "Product", 1, 100)
	err := suite.product.CreateAtomic(context.Background(), product)
	suite.Require().Error(err)
}

func (suite *test) TestCreateAtomicWithStockCapacityExceeded() {
	stock := factory.NewStockDefault()
	product := factory.NewProduct(uuid.NewString(), stock.ID, "Product", 2, 100)

	stock.Capacity = 1
	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))
	suite.Require().Error(suite.product.CreateAtomic(context.Background(), product))
}

func (suite *test) TestUpdateAtomic() {
	stock := factory.NewStockDefault()
	product := factory.NewProduct(uuid.NewString(), stock.ID, "Product", 10, 100)

	stock.Capacity = 30

	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))
	suite.Require().NoError(suite.product.CreateAtomic(context.Background(), product))

	product.Qtd = 20
	suite.Require().NoError(suite.product.UpdateAtomic(context.Background(), product))

	found, err := suite.product.ByID(context.Background(), product.ID)
	suite.Require().NoError(err)
	suite.Equal(product, found)
}

func (suite *test) TestUpdateAtomicWithStockNotFound() {
	product := factory.NewProduct(uuid.NewString(), uuid.NewString(), "Product", 1, 100)
	err := suite.product.UpdateAtomic(context.Background(), product)
	suite.Require().Error(err)
}

func (suite *test) TestUpdateAtomicWithStockCapacityExceeded() {
	stock := factory.NewStockDefault()
	product := factory.NewProduct(uuid.NewString(), stock.ID, "Product", 2, 100)

	stock.Capacity = 1
	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))

	suite.Require().Error(suite.product.UpdateAtomic(context.Background(), product))
}

func (suite *test) TestDelete() {
	stock := factory.NewStockDefault()
	product := factory.NewProduct(uuid.NewString(), stock.ID, "Product", 5, 100)

	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))
	suite.Require().NoError(suite.product.CreateAtomic(context.Background(), product))

	suite.Require().NoError(suite.product.DeleteAtomic(context.Background(), product.ID))

	_, err := suite.product.ByID(context.Background(), product.ID)
	suite.Require().Error(err)
}

func (suite *test) TestDeleteWithProductNotFound() {
	err := suite.product.DeleteAtomic(context.Background(), uuid.NewString())
	suite.Require().Error(err)
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
