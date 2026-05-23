package engine

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/EmilioQ74/EmBackUp/internal/adapters"
	"github.com/EmilioQ74/EmBackUp/internal/compress"
	"github.com/EmilioQ74/EmBackUp/internal/storage"
)

type Engine struct {
	adapter   adapters.Adapter
	store     storage.Storage
	compress  compress.Algorithm
	retention int
	log       *slog.Logger
}

func New(a adapters.Adapter, s storage.Storage, algo compress.Algorithm, retention int, log *slog.Logger) *Engine {
	return &Engine{adapter: a, store: s, compress: algo, retention: retention, log: log}
}

func (e *Engine) Backup(ctx context.Context, cfg adapters.DBConfig) error {
	if err := e.adapter.Validate(cfg); err != nil {
		return err
	}

	start := time.Now()
	var ext string
	switch e.compress {
	case compress.Gzip:
		ext = ".dump.gz"
	case compress.Zstd:
		ext = ".dump.zst"
	case compress.None:
		ext = ".dump"
	default:
		return fmt.Errorf("unsupported compression algorithm: %q", e.compress)
	}

	key := fmt.Sprintf(
		"%s/%s/%s%s",
		e.adapter.Name(),
		cfg.Database,
		start.UTC().Format("2006-01-02T150405"),
		ext,
	)

	e.log.Info("backup started", "db", cfg.Database, "key", key)

	dumpR, err := e.adapter.Dump(ctx, cfg)
	if err != nil {
		return fmt.Errorf("dump failed: %w", err)
	}
	defer dumpR.Close()

	pr, pw := io.Pipe()
	compW, err := compress.Wrap(pw, e.compress)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		_, copyErr := io.Copy(compW, dumpR)
		if closeErr := compW.Close(); copyErr == nil {
			copyErr = closeErr
		}
		pw.CloseWithError(copyErr)
		errCh <- copyErr
	}()

	if err := e.store.Upload(ctx, key, pr); err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	if err := <-errCh; err != nil {
		return fmt.Errorf("compress/copy failed: %w", err)
	}

	e.log.Info("backup complete",
		"db", cfg.Database,
		"key", key,
		"duration", time.Since(start).Round(time.Millisecond),
	)

	if e.retention > 0 {
		_ = e.pruneOld(ctx, cfg.Database)
	}
	return nil
}

func (e *Engine) pruneOld(ctx context.Context, db string) error {
	items, err := e.store.List(ctx, e.adapter.Name()+"/"+db)
	if err != nil {
		return err
	}
	if len(items) <= e.retention {
		return nil
	}
	for _, old := range items[:len(items)-e.retention] {
		if err := e.store.Delete(ctx, old.Key); err != nil {
			e.log.Warn("prune failed", "key", old.Key, "err", err)
		}
	}
	return nil
}

func (e *Engine) Restore(ctx context.Context, cfg adapters.DBConfig, key string) error {
	e.log.Info("restore started", "key", key)

	r, err := e.store.Download(ctx, key)
	if err != nil {
		return fmt.Errorf("download %s: %w", key, err)
	}
	defer r.Close()

	algo := compress.AlgorithmFromKey(key)
	reader, err := compress.Unwrap(r, algo)
	if err != nil {
		return fmt.Errorf("decompression init (%s): %w", algo, err)
	}

	return e.adapter.Restore(ctx, cfg, reader)
}

func (e *Engine) List(ctx context.Context, db string) ([]storage.BackupMeta, error) {
	prefix := e.adapter.Name() + "/" + db
	return e.store.List(ctx, prefix)
}
