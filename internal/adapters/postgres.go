package adapters

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type PostgresAdapter struct{}

func (p *PostgresAdapter) Name() string { return "postgres" }

func (p *PostgresAdapter) Validate(cfg DBConfig) error {
	if cfg.Database == "" {
		return fmt.Errorf("postgress: database name is required")
	}
	return nil
}

func (p *PostgresAdapter) Dump(ctx context.Context, cfg DBConfig) (io.ReadCloser, error) {
	args := []string{
		"-h", cfg.Host,
		"-p", fmt.Sprint(cfg.Port),
		"-U", cfg.User,
		"-d", cfg.Database,
		"--no-password",
		"-F", "c",
		"--verbose",
	}

	cmd := exec.CommandContext(ctx, "pg_dump", args...)
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("postgress dump pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("postgress dump start: %w", err)
	}

	return &cmdReadCloser{ReadCloser: stdout, cmd: cmd}, nil

}

func (p *PostgresAdapter) Restore(ctx context.Context, cfg DBConfig, src io.Reader) error {
	args := []string{
		"-h", cfg.Host, "-p", fmt.Sprint(cfg.Port), "-U", cfg.User, "-d", cfg.Database, "--no-password", "--clean",
		"--if-exists",
		"--single-transaction",
		"--verbose",
	}

	cmd := exec.CommandContext(ctx, "pg_restore", args...)
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("PGPASSWORD=%s", cfg.Password))
	cmd.Stdin = src

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_restore: %w - %s", err, out)
	}
	return nil

}
