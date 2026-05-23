package adapters

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type SQLiteAdapter struct{}

func (s *SQLiteAdapter) Name() string { return "sqlite" }

func (s *SQLiteAdapter) Validate(cfg DBConfig) error {
	if cfg.Database == "" {
		return fmt.Errorf("sqlite: database file path is required")
	}
	if _, err := os.Stat(cfg.Database); os.IsNotExist(err) {
		return fmt.Errorf("sqlite: database file not found: %s", cfg.Database)
	}
	return nil
}

func (s *SQLiteAdapter) Dump(ctx context.Context, cfg DBConfig) (io.ReadCloser, error) {
	// sqlite3 <file> .dump outputs full SQL dump to stdout
	cmd := exec.CommandContext(ctx, "sqlite3", cfg.Database, ".dump")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("sqlite dump pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("sqlite dump start: %w", err)
	}

	return &cmdReadCloser{ReadCloser: stdout, cmd: cmd}, nil
}

func (s *SQLiteAdapter) Restore(ctx context.Context, cfg DBConfig, src io.Reader) error {
	// delete existing db file and restore from SQL dump
	if err := os.Remove(cfg.Database); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("sqlite: remove existing db: %w", err)
	}

	cmd := exec.CommandContext(ctx, "sqlite3", cfg.Database)
	cmd.Stdin = src

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sqlite restore: %w - %s", err, out)
	}
	return nil
}
