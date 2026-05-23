package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type LocalStorage struct{ basePath string }

func NewLocal(path string) (*LocalStorage, error) {
	if err := os.MkdirAll(path, 0o750); err != nil {
		return nil, err
	}
	return &LocalStorage{basePath: path}, nil
}

var _ Storage = (*LocalStorage)(nil)

func (l *LocalStorage) Upload(_ context.Context, key string, r io.Reader) error {
	dest := filepath.Join(l.basePath, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(dest), 0o750); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (l *LocalStorage) List(_ context.Context, prefix string) ([]BackupMeta, error) {
	var metas []BackupMeta
	root := filepath.Join(l.basePath, filepath.FromSlash(prefix))
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, _ := e.Info()
		metas = append(metas, BackupMeta{
			Key:       strings.TrimPrefix(filepath.ToSlash(filepath.Join(root, e.Name())), l.basePath+"/"),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].CreatedAt.Before(metas[j].CreatedAt)
	})
	return metas, nil
}

func (l *LocalStorage) Delete(_ context.Context, key string) error {
	dest := filepath.Join(l.basePath, filepath.FromSlash(key))
	return os.Remove(dest)
}

func (l *LocalStorage) Download(_ context.Context, key string) (io.ReadCloser, error) {
	dest := filepath.Join(l.basePath, filepath.FromSlash(key))
	return os.Open(dest)
}
