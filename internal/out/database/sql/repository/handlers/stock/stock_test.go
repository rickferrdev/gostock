package stock

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/domain/factory"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository/handlers/product"
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository/schema"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type test struct {
	suite.Suite
	db      *bun.DB
	stock   ports.StockRepository
	product ports.ProductRepository
}

func (suite *test) SetupTest() {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file:gostock_stock_test?mode=memory&cache=shared")
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

	suite.stock, err = New(suite.db)
	suite.Require().NoError(err)

	suite.product, err = product.New(suite.db)
	suite.Require().NoError(err)
}

func (suite *test) TearDownTest() {
	suite.Require().NoError(suite.db.Close())
}

func (suite *test) TestAllWhenEmpty() {
	list, err := suite.stock.All(context.Background())

	suite.Require().NoError(err)
	suite.Empty(list)
	suite.T().Logf("stocks found: %d", len(list))
}

func (suite *test) TestAllWithData() {
	stock := factory.NewStockDefault()

	err := suite.stock.CreateAtomic(context.Background(), stock)
	suite.Require().NoError(err)

	list, err := suite.stock.All(context.Background())
	suite.Require().NoError(err)

	suite.Equal([]*domain.Stock{stock}, list)

	for _, item := range list {
		find, err := suite.stock.ByID(context.Background(), item.ID)
		suite.Require().NoError(err)
		suite.Equal(item, find)
	}
}

func (suite *test) TestByIDWhenNotFound() {
	stock, err := suite.stock.ByID(context.Background(), uuid.NewString())
	suite.Error(err)
	suite.Nil(stock)
}

func (suite *test) TestByIDWithData() {
	stock := factory.NewStockDefault()

	err := suite.stock.CreateAtomic(context.Background(), stock)
	suite.NoError(err)

	find, err := suite.stock.ByID(context.Background(), stock.ID)
	suite.NoError(err)
	suite.Equal(stock, find)
}

func (suite *test) TestCreate() {
	stock := factory.NewStockDefault()
	err := suite.stock.CreateAtomic(context.Background(), stock)
	suite.NoError(err)

	find, err := suite.stock.ByID(context.Background(), stock.ID)
	suite.NoError(err)
	suite.Equal(stock, find)
}

func (suite *test) TestUpdate() {
	stock := factory.NewStockDefault()

	err := suite.stock.CreateAtomic(context.Background(), stock)
	suite.NoError(err)

	stock.Name = "Updated Stock Name"
	err = suite.stock.UpdateAtomic(context.Background(), stock)
	suite.NoError(err)

	find, err := suite.stock.ByID(context.Background(), stock.ID)
	suite.NoError(err)
	suite.Equal(stock, find)
}

func (suite *test) TestUpdateWhenNotFound() {
	stock := factory.NewStockDefault()
	err := suite.stock.UpdateAtomic(context.Background(), stock)
	suite.Error(err)
}

func (suite *test) TestDelete() {
	stock := factory.NewStockDefault()
	err := suite.stock.CreateAtomic(context.Background(), stock)
	suite.NoError(err)

	err = suite.stock.DeleteAtomic(context.Background(), stock.ID)
	suite.NoError(err)

	find, err := suite.stock.ByID(context.Background(), stock.ID)
	suite.Error(err)
	suite.Nil(find)
}

func (suite *test) TestDeleteWhenNotFound() {
	err := suite.stock.DeleteAtomic(context.Background(), uuid.NewString())
	suite.Error(err)
}

func (suite *test) TestDeleteCascade() {
	stock, product := factory.NewStockWithOneProducts(1, 100)

	suite.Require().NoError(suite.stock.CreateAtomic(context.Background(), stock))

	suite.Require().NoError(suite.product.DeleteAtomic(context.Background(), product.ID))
	suite.Require().NoError(suite.stock.DeleteAtomic(context.Background(), stock.ID))

	find, err := suite.product.ByID(context.Background(), product.ID)
	suite.Require().Error(err)
	suite.Nil(find)
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
