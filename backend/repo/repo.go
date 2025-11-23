package repo

import (
	"context"

	"github.com/gofiber/fiber/v3/log"

	"project/database"
	"project/repo/postgres"
	"project/repo/sqlite"
)

type Repository interface {
	Ping(ctx context.Context) error
	Init(ctx context.Context) error
}

func NewRepo(db *database.Database) Repository {
	var r Repository
	switch db.Driver() {
	case "sqlite":
		r = sqlite.New(db)
	case "postgres":
		return postgres.New(db)
	default:
		log.Fatal("unsupported database driver")
	}

	log.Info("Initializing database schema...")
	err := r.Init(context.Background())
	if err != nil {
		log.Fatalf("failed to initialize database schema: %v", err)
	}

	return r
}
