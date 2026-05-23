package adapters

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type MySQLAdapter struct{}

func (m *MySQLAdapter) Name() string { return "mysql" }

func (m *MySQLAdapter) Validate(cfg DBConfig) error {
	if cfg.Database == "" {
		return fmt.Errorf("mysql: database name is required")
	}
	return nil
}

func (m *MySQLAdapter) Dump(ctx context.Context, cfg DBConfig) (io.ReadCloser, error) {
	args := []string{
		"-h", cfg.Host,
		"-P", fmt.Sprint(cfg.Port),
		"-u", cfg.User,
		fmt.Sprintf("-p%s", cfg.Password),
		"--single-transaction",
		"--routines",
		"--triggers",
		"--events",
		"--set-gtid-purged=OFF",
		"--no-tablespaces",
		cfg.Database,
	}

	cmd := exec.CommandContext(ctx, "mysqldump", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("mysql dump pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("mysql dump start: %w", err)
	}

	return &cmdReadCloser{ReadCloser: stdout, cmd: cmd}, nil
}
func (m *MySQLAdapter) Restore(ctx context.Context, cfg DBConfig, src io.Reader) error {
	args := []string{
		"-h", cfg.Host,
		"-P", fmt.Sprint(cfg.Port),
		"-u", cfg.User,
		fmt.Sprintf("-p%s", cfg.Password),
		cfg.Database,
	}

	cmd := exec.CommandContext(ctx, "mysql", args...)
	cmd.Stdin = src

	stderr := &strings.Builder{}
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mysql restore: %w - %s", err, stderr.String())
	}
	return nil
}
