<div align="center">

# 📦 GoStock

<p align="center">
  <em>A RESTful API for managing products and stocks</em><br/>
  <em>Built with Go · Hexagonal Architecture · Production-ready</em>
</p>

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)
[![CI](https://img.shields.io/github/actions/workflow/status/rickferrdev/gostock/ci.yml?branch=main&style=for-the-badge&logo=github-actions&logoColor=white&label=CI)](https://github.com/rickferrdev/gostock/actions)

[![Fiber](https://img.shields.io/badge/Fiber-v3-00ACD7?style=flat-square&logo=go&logoColor=white)](https://github.com/gofiber/fiber)
[![Bun ORM](https://img.shields.io/badge/Bun-ORM-FF6B6B?style=flat-square&logo=go&logoColor=white)](https://github.com/uptrace/bun)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat-square&logo=sqlite&logoColor=white)](https://www.sqlite.org/)
[![Uber FX](https://img.shields.io/badge/Uber_FX-DI-1D2127?style=flat-square&logo=uber&logoColor=white)](https://github.com/uber-go/fx)
[![Swagger](https://img.shields.io/badge/Swagger-UI-85EA2D?style=flat-square&logo=swagger&logoColor=black)](http://localhost:8080/api/v1/docs)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://www.docker.com/)

[![Go Report Card](https://goreportcard.com/badge/github.com/rickferrdev/gostock?style=flat-square)](https://goreportcard.com/report/github.com/rickferrdev/gostock)
[![golangci-lint](https://img.shields.io/badge/golangci--lint-passing-4BC51D?style=flat-square&logo=go)](https://golangci-lint.run/)

</div>

---

## 🚀 Tech Stack

| Layer                | Technology                                                                  |
| -------------------- | --------------------------------------------------------------------------- |
| Language             | Go 1.25+                                                                    |
| HTTP Framework       | [Fiber v3](https://github.com/gofiber/fiber)                                |
| ORM / Query Builder  | [Bun](https://github.com/uptrace/bun)                                       |
| Database             | SQLite (via `sqliteshim`)                                                   |
| Dependency Injection | [Uber FX](https://github.com/uber-go/fx)                                    |
| Validation           | [go-playground/validator](https://github.com/go-playground/validator)       |
| API Documentation    | [Swaggo](https://github.com/swaggo/swag) + Swagger UI                       |
| Logging              | `log/slog` + [slog-betterstack](https://github.com/samber/slog-betterstack) |
| Testing              | `testify` + [uber/mock](https://github.com/uber-go/mock)                    |
| Linting              | [golangci-lint](https://golangci-lint.run/)                                 |
| Containerization     | Docker (multi-stage build)                                                  |

## ✨ Features

- **Product management** — Create, read, update, and delete products
- **Stock management** — Create, read, update, and delete stocks with capacity control
- **Capacity enforcement** — Products can only be added to a stock if space is available
- **Health check** endpoint with service status reporting
- **Swagger UI** served at `/api/v1/docs`
- **Rate limiting**, request ID, CORS, and recovery middlewares out of the box
- **Structured logging** with Betterstack integration
- **Database migrations** via a standalone migrator binary

## 📋 Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (optional, for containerized run)
- [Make](https://www.gnu.org/software/make/)

## ⚙️ Environment Variables

> [!TIP]
> Create a `.env` file at the root of the project before running.

```env
APP_SERVER_PORT=8080
APP_DATABASE_URI=./database/app.db
APP_SERVER_ORIGIN_FRONT=http://localhost:3000

# Optional – Betterstack log ingestion
APP_BETTERSTACK_TOKEN=your_token_here
APP_BETTERSTACK_INGESTING_HOST=https://in.logs.betterstack.com
```

## 🛠️ Running Locally

### Install dependencies

```bash
go mod download
```

### Run database migrations

```bash
make migrate-up
```

### Start the server

```bash
make run
```

> [!NOTE]
> The API will be available at `http://localhost:8080/api/v1`.
> Swagger UI is available at `http://localhost:8080/api/v1/docs`.

---

## 🐳 Running with Docker

```bash
# Build and start
make docker-up

# Stop
make docker-down
```

> [!NOTE]
> The Docker setup uses a **multi-stage build**:
>
> 1. **Builder** — compiles the API binary and the migrator, generates Swagger docs
> 2. **Runtime** — minimal Alpine image with only the compiled binaries

---

## 🧪 Testing

```bash
# Run all tests with race detector
make test

# Generate HTML coverage report
make cover
```

## 📐 API Endpoints

<details>
<summary><b>🏥 Health</b></summary>

| Method | Path             | Description                   |
| ------ | ---------------- | ----------------------------- |
| `GET`  | `/api/v1/health` | Returns service health status |

</details>

<details>
<summary><b>📦 Stocks</b></summary>

| Method   | Path                 | Description        |
| -------- | -------------------- | ------------------ |
| `GET`    | `/api/v1/stocks`     | List all stocks    |
| `POST`   | `/api/v1/stocks`     | Create a new stock |
| `GET`    | `/api/v1/stocks/:id` | Get a stock by ID  |
| `PUT`    | `/api/v1/stocks/:id` | Update a stock     |
| `DELETE` | `/api/v1/stocks/:id` | Delete a stock     |

</details>

<details>
<summary><b>🛍️ Products</b></summary>

| Method   | Path                               | Description            |
| -------- | ---------------------------------- | ---------------------- |
| `GET`    | `/api/v1/products`                 | List all products      |
| `POST`   | `/api/v1/products`                 | Create a new product   |
| `GET`    | `/api/v1/products/:id`             | Get a product by ID    |
| `GET`    | `/api/v1/products/stock/:stock_id` | List products by stock |
| `PUT`    | `/api/v1/products/:id`             | Update a product       |
| `DELETE` | `/api/v1/products/:id`             | Delete a product       |

</details>

> Full request/response schemas are available in the **Swagger UI** at `/api/v1/docs`.

## 🗄️ Database Migrations

The project ships a standalone migrator binary with the following commands:

```bash
make migrate-init    # Initialize the migrations table
make migrate-up      # Apply pending migrations
make migrate-down    # Roll back the last migration
make migrate-create  # Scaffold a new migration file
```

## 🧹 Development Utilities

```bash
make fmt       # Format all Go files
make vet       # Run go vet
make lint      # Install and run golangci-lint
make tidy      # Tidy go.mod / go.sum
make mock      # Regenerate mocks with mockgen
make swagger   # Regenerate Swagger docs
make clean     # Remove build artifacts and coverage files
```

## 🔄 CI/CD

GitHub Actions runs on every push and pull request to `main`:

1. **golangci-lint** — static analysis
2. **tests** — `make test` (with `-race` flag)

Dependabot is configured to keep **DevContainer** features up to date weekly.

## 🧑‍💻 Dev Container

A fully configured Dev Container is included for VS Code and GitHub Codespaces, based on the official `go:2-1.26-trixie` image and featuring:

- **Fish shell**
- **Docker-outside-of-Docker**
- **`act`** for running GitHub Actions locally

## 📄 License

This project is licensed under the [MIT License](LICENSE).

<div align="center">
  <sub>Built with ❤️ using Go & Hexagonal Architecture</sub>
</div>
