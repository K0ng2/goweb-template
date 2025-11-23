package postgres

import (
	"context"
	_ "embed"

	"github.com/go-jet/jet/v2/postgres"
)

//go:embed queries/init.sql
var initSchema string

func (r *Postgres) Init(ctx context.Context) error {
	stmt := postgres.RawStatement(initSchema)

	_, err := stmt.ExecContext(ctx, r.exec)
	if err != nil {
		return err
	}

	return nil
}
