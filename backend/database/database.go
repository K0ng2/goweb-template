package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type Database struct {
	conn   *sql.DB
	driver string
}

type Executor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func New(dsn string) (*Database, error) {
	var driver string
	switch {
	case strings.HasPrefix(dsn, "sqlite://"):
		driver = "sqlite"
		dsn = strings.TrimPrefix(dsn, "sqlite://file:")

	case strings.HasPrefix(dsn, "postgres://"):
		driver = "postgres"
		dsn = strings.TrimPrefix(dsn, "postgres://")

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dsn)
	}

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &Database{
		conn:   conn,
		driver: driver,
	}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

func (db *Database) Conn() *sql.DB {
	return db.conn
}

func (db *Database) Driver() string {
	return db.driver
}

func (db *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.conn.BeginTx(ctx, opts)
}

func (db *Database) PingContext(ctx context.Context) error {
	return db.conn.PingContext(ctx)
}
