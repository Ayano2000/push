package pgsql

import (
	"context"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Pgsql struct {
	Config *config.Config
	Client *pgxpool.Pool
}

func NewPgsql(ctx context.Context, config *config.Config) (*Pgsql, error) {
	pgsqlConn, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return &Pgsql{
		Config: config,
		Client: pgsqlConn,
	}, err
}

func (pgsql *Pgsql) CreateBucket(ctx context.Context, bucket types.Bucket) error {
	_, err := pgsql.Client.Exec(ctx, `
		INSERT INTO buckets (name, path, method, description, jq_filter, forward_to, preserve_payload) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		bucket.Name,
		bucket.Path,
		bucket.Method,
		bucket.Description,
		bucket.JQFilter,
		bucket.ForwardTo,
		bucket.PreservePayload,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (pgsql *Pgsql) GetBuckets(ctx context.Context) ([]types.Bucket, error) {
	rows, err := pgsql.Client.Query(ctx, `SELECT * FROM buckets`)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var buckets []types.Bucket
	for rows.Next() {
		var bucket types.Bucket
		err = rows.Scan(
			&bucket.Name,
			&bucket.Path,
			&bucket.Method,
			&bucket.Description,
			&bucket.JQFilter,
			&bucket.ForwardTo,
			&bucket.PreservePayload)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		buckets = append(buckets, bucket)
	}

	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return buckets, nil
}

func (pgsql *Pgsql) GetBucketByName(ctx context.Context, name string) (types.Bucket, error) {
	rows, err := pgsql.Client.Query(ctx,
		`SELECT * FROM buckets WHERE name = $1 limit 1`,
		name,
	)
	if err != nil {
		return types.Bucket{}, errors.WithStack(err)
	}
	defer rows.Close()

	var buckets []types.Bucket
	for rows.Next() {
		var bucket types.Bucket
		err = rows.Scan(
			&bucket.Name,
			&bucket.Path,
			&bucket.Method,
			&bucket.Description,
			&bucket.JQFilter,
			&bucket.ForwardTo,
			&bucket.PreservePayload)
		if err != nil {
			return types.Bucket{}, errors.WithStack(err)
		}
		buckets = append(buckets, bucket)
	}

	if rows.Err() != nil {
		return types.Bucket{}, errors.WithStack(rows.Err())
	}

	return buckets[0], nil
}

func (pgsql *Pgsql) Close() {
	defer pgsql.Client.Close()
}
