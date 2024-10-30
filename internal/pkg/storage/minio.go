package storage

import (
	"context"
	"fmt"
	"github.com/Ayano2000/push/internal/config"
	"github.com/Ayano2000/push/internal/pkg/types"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"io"
	"strings"
	"time"
)

type MinIOStorage struct {
	client *minio.Client
}

// NewMinIOStorage creates new MinIO storage instance
func NewMinIOStorage(config *config.Config) (*MinIOStorage, error) {
	minioConn, err := minio.New(config.MinioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOStorage{
		client: minioConn,
	}, nil
}

func (m *MinIOStorage) CreateBucket(ctx context.Context, webhook types.Webhook) error {
	exists, err := m.client.BucketExists(ctx, webhook.Name)
	if err != nil {
		return errors.WithStack(err)
	}
	if exists {
		return errors.WithStack(errors.New("webhook already exists"))
	}

	err = m.client.MakeBucket(ctx, webhook.Name, minio.MakeBucketOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *MinIOStorage) PutObject(ctx context.Context, bucketName, payload string) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = m.client.PutObject(
		ctx,
		bucketName,
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

func (m *MinIOStorage) GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	return m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (m *MinIOStorage) GetObjects(ctx context.Context, bucket string) ([]string, error) {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !exists {
		return nil, errors.WithStack(errors.New("webhook does not exist"))
	}

	var objects []string
	objectChannel := m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{})
	for object := range objectChannel {
		object, err := m.client.GetObject(ctx, bucket, object.Key, minio.GetObjectOptions{})
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

func (m *MinIOStorage) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m *MinIOStorage) GetPresignedURL(ctx context.Context, bucketName, objectName string, expires time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (m *MinIOStorage) Close() error {
	// MinIO client doesn't require explicit closing
	return nil
}

var _ ObjectStoreHandler = (*MinIOStorage)(nil)
