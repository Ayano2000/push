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

func (pgsql *Pgsql) CreateWebhook(ctx context.Context, webhook types.Webhook) error {
	_, err := pgsql.Client.Exec(ctx, `
		INSERT INTO webhooks (name, path, method, description, jq_filter, forward_to, preserve_payload) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		webhook.Name,
		webhook.Path,
		webhook.Method,
		webhook.Description,
		webhook.JQFilter,
		webhook.ForwardTo,
		webhook.PreservePayload,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (pgsql *Pgsql) GetWebhooks(ctx context.Context) ([]types.Webhook, error) {
	rows, err := pgsql.Client.Query(ctx, `SELECT * FROM webhooks`)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var webhooks []types.Webhook
	for rows.Next() {
		var webhook types.Webhook
		err = rows.Scan(
			&webhook.Name,
			&webhook.Path,
			&webhook.Method,
			&webhook.Description,
			&webhook.JQFilter,
			&webhook.ForwardTo,
			&webhook.PreservePayload)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		webhooks = append(webhooks, webhook)
	}

	if rows.Err() != nil {
		return nil, errors.WithStack(rows.Err())
	}

	return webhooks, nil
}

func (pgsql *Pgsql) GetWebhookByName(ctx context.Context, name string) (types.Webhook, error) {
	rows, err := pgsql.Client.Query(ctx,
		`SELECT * FROM webhooks WHERE name = $1 limit 1`,
		name,
	)
	if err != nil {
		return types.Webhook{}, errors.WithStack(err)
	}
	defer rows.Close()

	var webhooks []types.Webhook
	for rows.Next() {
		var webhook types.Webhook
		err = rows.Scan(
			&webhook.Name,
			&webhook.Path,
			&webhook.Method,
			&webhook.Description,
			&webhook.JQFilter,
			&webhook.ForwardTo,
			&webhook.PreservePayload)
		if err != nil {
			return types.Webhook{}, errors.WithStack(err)
		}
		webhooks = append(webhooks, webhook)
	}

	if rows.Err() != nil {
		return types.Webhook{}, errors.WithStack(rows.Err())
	}

	return webhooks[0], nil
}

func (pgsql *Pgsql) Close() {
	defer pgsql.Client.Close()
}
