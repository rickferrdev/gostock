# ── Stage 1: Builder ───────────────────────────────────────────────────────────
FROM golang:alpine AS builder

RUN apk add --no-cache make git

WORKDIR /app

COPY go.mod go.sum Makefile ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    make swagger

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build && make migrator

# ── Stage 2: Runtime ───────────────────────────────────────────────────────────
FROM alpine:3.21 AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    mkdir -p /app/database

COPY --from=builder /app/build/gostock ./app
COPY --from=builder /app/build/migrator ./migrator

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "-c", "./migrator init && ./migrator migrate && ./app"]
