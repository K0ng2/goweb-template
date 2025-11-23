package sqlite

import (
	"context"
	"database/sql"

	"project/database"
)

type Sqlite struct {
	db   *database.Database
	exec database.Executor
}

func New(db *database.Database) *Sqlite {
	return &Sqlite{db: db, exec: db.Conn()}
}

func (r *Sqlite) WithTx(tx *sql.Tx) *Sqlite {
	return &Sqlite{
		db:   r.db,
		exec: tx,
	}
}

func (r *Sqlite) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
