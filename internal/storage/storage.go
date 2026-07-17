package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) error
	Delete(ctx context.Context, key string) error
	PresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}
