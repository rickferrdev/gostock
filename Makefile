GO_CMD			= go
GO_MAIN			= ./cmd/api/main.go

APP_NAME		= gostock
APP_MODULE		= github.com/rickferrdev/gostock

MIGRATOR_MAIN	= ./cmd/migrator/main.go
MIGRATOR_BUILD	= ./build/migrator

APP_BUILD		= ./build
LDFLAGS			= -ldflags="-s -w"

DOCKER_CMD		= docker
DOCKER_COMPOSE	= docker/compose.yml
DOCKER_FILE		= docker/Dockerfile

SWAGGER_CMD 	= swag
SWAGGER_DIR		= ./docs/swagger
SWAGGER_GEN		= ./internal/infra/server/server.go
SWAGGER_FLAGS	= -g $(SWAGGER_GEN) --output $(SWAGGER_DIR) --parseDependency --parseInternal

.DEFAULT_GOAL 	:= build

.PHONY: fmt build run test cover lint tidy swagger vet \
        mock docker-build docker-up docker-down clean verify download vul

fmt:
	$(GO_CMD) fmt ./...

run:
	$(GO_CMD) run $(GO_MAIN)

tidy:
	$(GO_CMD) mod tidy

test:
	$(GO_CMD) test -v -count=1 -parallel=4 ./...

cover:
	$(GO_CMD) test -coverprofile=coverage.out ./...
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html

mock:
	$(GO_CMD) install go.uber.org/mock/mockgen@latest
	$(GO_CMD) generate ./...

lint:
	$(GO_CMD) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
	golangci-lint run ./...

swagger:
	$(GO_CMD) install github.com/swaggo/swag/cmd/swag@latest
	$(SWAGGER_CMD) init $(SWAGGER_FLAGS)

docker-build:
	$(DOCKER_CMD) build -f $(DOCKER_FILE) -t $(APP_NAME):$(VERSION) .

docker-up:
	$(DOCKER_CMD) compose -f $(DOCKER_COMPOSE) up -d --build

docker-down:
	$(DOCKER_CMD) compose -f $(DOCKER_COMPOSE) down

clean:
	rm -rf $(APP_BUILD) coverage.out coverage.html

vet:
	$(GO_CMD) vet ./...

migrate-init:
	$(GO_CMD) run ./cmd/migrator/main.go init

migrate-up:
	$(GO_CMD) run ./cmd/migrator/main.go migrate

migrate-down:
	$(GO_CMD) run ./cmd/migrator/main.go rollback

migrate-create:
	@read -p "migration name: " name; \
	$(GO_CMD) run ./cmd/migrator/main.go create $$name

verify:
	$(GO_CMD) mod verify

build:
	@mkdir -p $(APP_BUILD)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build $(LDFLAGS) -o $(APP_BUILD)/$(APP_NAME) $(GO_MAIN)

download:
	$(GO_CMD) mod download

migrator:
	@mkdir -p $(APP_BUILD)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_CMD) build $(LDFLAGS) -trimpath -o $(MIGRATOR_BUILD) $(MIGRATOR_MAIN)

vul:
	$(GO_CMD) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck -C  ./...
