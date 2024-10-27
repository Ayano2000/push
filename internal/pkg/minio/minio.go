package minio

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/types"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type Minio struct {
	Config *config.Config
	Client *minio.Client
}

func NewMinio(config *config.Config) (*Minio, error) {
	minioConn, err := minio.New(config.MinioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Minio{
		Config: config,
		Client: minioConn,
	}, err
}

// CreateBucket will throw an error if the requested name is already in use
func (m *Minio) CreateBucket(ctx context.Context, webhook types.Webhook) error {
	exists, err := m.Client.BucketExists(ctx, webhook.Name)
	if err != nil {
		return errors.WithStack(err)
	}
	if exists {
		return errors.WithStack(errors.New("webhook already exists"))
	}

	err = m.Client.MakeBucket(ctx, webhook.Name, minio.MakeBucketOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// PutObject is a wrapper for the actual minio PutObject call
func (m *Minio) PutObject(ctx context.Context, webhook types.Webhook, payload string) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = m.Client.PutObject(
		ctx,
		webhook.Name,
		fmt.Sprintf("%s.json", uid.String()),
		io.NopCloser(strings.NewReader(payload)),
		int64(len(payload)),
		minio.PutObjectOptions{
			ContentType: "application/json",
			// todo: UserTags, so we know if this is a pre or post transform value
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// GetObjects should as much as possible be run in a separate goroutine
// the get object calls run one at a time since there is no bulk operation
func (m *Minio) GetObjects(ctx context.Context, webhook types.Webhook) ([]string, error) {
	exists, err := m.Client.BucketExists(ctx, webhook.Name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !exists {
		return nil, errors.WithStack(errors.New("webhook does not exist"))
	}

	var objects []string
	objectChannel := m.Client.ListObjects(ctx, webhook.Name, minio.ListObjectsOptions{})
	for object := range objectChannel {
		object, err := m.Client.GetObject(ctx, webhook.Name, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		payload, err := io.ReadAll(object)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		objects = append(objects, string(payload))
	}

	return objects, nil
}
