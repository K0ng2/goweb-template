package postgres

import (
	"context"
	"database/sql"

	"project/database"
)

type Postgres struct {
	db   *database.Database
	exec database.Executor
}

func New(db *database.Database) *Postgres {
	return &Postgres{db: db, exec: db.Conn()}
}

func (r *Postgres) WithTx(tx *sql.Tx) *Postgres {
	return &Postgres{
		db:   r.db,
		exec: tx,
	}
}

func (r *Postgres) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
