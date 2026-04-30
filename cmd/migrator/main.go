package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rickferrdev/dotenv"
	"github.com/rickferrdev/gostock/internal/config/env"
	"github.com/rickferrdev/gostock/internal/database/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/migrate"
)

func main() {
	dotenv.Collect()
	var e env.Env

	if err := dotenv.Unmarshal(&e); err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(e.APP_DATABASE_URI)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("error: %s\n", err)
	}

	fl := fmt.Sprintf("file:%v?cache=shared&mode=rwc", e.APP_DATABASE_URI)
	sqldb, err := sql.Open(sqliteshim.ShimName, fl)
	if err != nil {
		log.Fatal(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	migrator := migrate.NewMigrator(db, migrations.Migrations)
	ctx := context.Background()
	cmd := "migrate"

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "init":
		if err := migrator.Init(ctx); err != nil {
			log.Fatal(err)
		}
		log.Println("migration control tables initialized")

	case "migrate":
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if group.IsZero() {
			log.Printf("no new migrations to run\n")
			return
		}

		log.Printf("success! Migrations run: %s\n", group)
	case "rollback":
		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if group.IsZero() {
			log.Printf("no migrations to roll back\n")
			return
		}

		log.Printf("migrations rolled back: %s\n", group)
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("enter the migration name. Ex: go run ./cmd/migrator/main.go create initial")
		}

		name := os.Args[2]
		files, err := migrator.CreateSQLMigrations(ctx, name)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			log.Printf("created: %s\n", f.Path)
		}
	default:
		log.Fatalf("unsupported command: %s", cmd)
	}

	if err := db.Close(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
