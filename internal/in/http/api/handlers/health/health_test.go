package health

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gostock/internal/config/env"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/rickferrdev/gostock/internal/infra/server"
	"github.com/rickferrdev/gostock/internal/platform/validator"
	"github.com/rickferrdev/gostock/internal/tests/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type test struct {
	suite.Suite
	ctrl    *gomock.Controller
	service *mocks.MockHealthService
	app     *fiber.App
	router  fiber.Router
}

func (suite *test) SetupTest() {
	env := env.Env{
		APP_SERVER_ORIGIN_FRONT: "*",
	}

	suite.app, suite.router = server.New(server.ServerParams{
		Validator: validator.New(),
		Env:       &env,
	})

	suite.Require().NotNil(suite.app, "O servidor Fiber não deveria ser nil")
	suite.Require().NotNil(suite.router, "O router não deveria ser nil")

	suite.ctrl = gomock.NewController(suite.T())
	suite.service = mocks.NewMockHealthService(suite.ctrl)
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

func (suite *test) TestHealth() {
	table := []struct {
		name   string
		desc   string
		expect int
		setup  func()
	}{
		{
			name:   "SuccesW",
			expect: fiber.StatusOK,
			setup: func() {
				suite.service.EXPECT().Check(gomock.Any()).Times(1).Return(&ports.HealthResponse{Status: "UP", Services: map[string]any{
					"database": "UP",
					"redis":    "UP",
				}})
			},
		},
		{
			name:   "Unavailable",
			expect: fiber.StatusServiceUnavailable,
			setup: func() {
				suite.service.EXPECT().Check(gomock.Any()).Times(1).Return(&ports.HealthResponse{Status: "DOWN", Services: map[string]any{
					"database": "DOWN",
					"redis":    "DOWN",
				}})
			},
		},
	}

	for _, tt := range table {
		suite.Run(tt.name, func() {
			tt.setup()

			request := httptest.NewRequest(fiber.MethodGet, "/api/v1/health", nil)
			response, err := suite.app.Test(request)
			suite.NoError(err)

			suite.Equal(tt.expect, response.StatusCode)
		})
	}

}

func TestRunner(t *testing.T) {
	suite.Run(t, new(test))
}
