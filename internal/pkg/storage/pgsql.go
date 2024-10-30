package storage

import (
	"context"
	"database/sql"
	"github.com/Ayano2000/push/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresDB implements Database interface
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresDB creates new PostgreSQL database instance
func NewPostgresDB(config *config.Config) (*PostgresDB, error) {
	db, err := sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

func (p *PostgresDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

func (p *PostgresDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

func (p *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, opts)
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) Ping() error {
	return p.db.Ping()
}

var _ DatabaseHandler = (*PostgresDB)(nil)
