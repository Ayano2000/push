package storage

import (
	"context"
	"database/sql"
	"github.com/Ayano2000/push/internal/pkg/types"
	"io"
	"time"
)

// ObjectStoreHandler defines an interface for object storage operations
type ObjectStoreHandler interface {
	CreateBucket(ctx context.Context, webhook types.Webhook) error
	// PutObject uploads a file to object storage
	PutObject(ctx context.Context, bucketName, payload string) error
	// GetObject downloads a file from object storage
	GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	// GetObjects returns a string array of all bucket contents
	GetObjects(ctx context.Context, bucketName string) ([]string, error)
	// DeleteObject deletes a file from object storage
	DeleteObject(ctx context.Context, bucketName, objectName string) error
	// GetPresignedURL returns temporary URL for object
	GetPresignedURL(ctx context.Context, bucketName, objectName string, expires time.Duration) (string, error)
	// Close connection
	Close() error
}

// DatabaseHandler defines an interface for database storage operations
type DatabaseHandler interface {
	// ExecContext executes a query with parameters
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// QueryContext queries multiple rows
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// QueryRowContext queries a single row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	// BeginTx begins a transaction
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	// Close connection
	Close() error
	// Ping can be used to connection health
	Ping() error
}
