package sqlite

import (
	"context"
	_ "embed"

	"github.com/go-jet/jet/v2/sqlite"
)

//go:embed queries/init.sql
var initSchema string

func (r *Sqlite) Init(ctx context.Context) error {
	stmt := sqlite.RawStatement(initSchema)

	_, err := stmt.ExecContext(ctx, r.exec)
	if err != nil {
		return err
	}

	return nil
}
