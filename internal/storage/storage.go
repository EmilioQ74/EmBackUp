package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, key string, r io.Reader) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	List(ctx context.Context, prefix string) ([]BackupMeta, error)
	Delete(ctx context.Context, key string) error
}

type BackupMeta struct {
	Key       string
	Size      int64
	CreatedAt time.Time
}
