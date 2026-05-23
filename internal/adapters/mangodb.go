package adapters

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type MongoAdapter struct{}

func (m *MongoAdapter) Name() string { return "mongodb" }

func (m *MongoAdapter) Validate(cfg DBConfig) error {
	if cfg.Database == "" {
		return fmt.Errorf("mongodb: database name is required")
	}
	return nil
}

func (m *MongoAdapter) Dump(ctx context.Context, cfg DBConfig) (io.ReadCloser, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	args := []string{
		"--uri", uri,
		"--db", cfg.Database,
		"--archive", // stream to stdout instead of creating files
		"--gzip",    // mongodb handles its own compression here
	}

	cmd := exec.CommandContext(ctx, "mongodump", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("mongodb dump pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("mongodb dump start: %w", err)
	}

	return &cmdReadCloser{ReadCloser: stdout, cmd: cmd}, nil
}

func (m *MongoAdapter) Restore(ctx context.Context, cfg DBConfig, src io.Reader) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	args := []string{
		"--uri", uri,
		"--db", cfg.Database,
		"--archive", // read from stdin
		"--gzip",    // decompress mongodb's own compression
		"--drop",    // drop collections before restoring (like --clean)
	}

	cmd := exec.CommandContext(ctx, "mongorestore", args...)
	cmd.Stdin = src

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mongorestore: %w - %s", err, out)
	}
	return nil
}
