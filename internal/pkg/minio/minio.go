package minio

import (
	"context"
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
func (m *Minio) CreateBucket(ctx context.Context, bucket types.Bucket) error {
	exists, err := m.Client.BucketExists(ctx, bucket.Name)
	if err != nil {
		return errors.WithStack(err)
	}
	if exists {
		return errors.WithStack(errors.New("bucket already exists"))
	}

	err = m.Client.MakeBucket(ctx, bucket.Name, minio.MakeBucketOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Minio) PutObject(ctx context.Context, bucket types.Bucket, payload string) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = m.Client.PutObject(
		ctx,
		bucket.Name,
		uid.String(),
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
