package adapters

import (
	"context"
	"io"
	"os/exec"
)

type Adapter interface {
	Dump(ctx context.Context, cfg DBConfig) (io.ReadCloser, error)
	Restore(ctx context.Context, cfg DBConfig, src io.Reader) error
	Validate(cfg DBConfig) error
	Name() string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string

	Options map[string]string
}

type cmdReadCloser struct {
	io.ReadCloser
	cmd *exec.Cmd
}

func (r *cmdReadCloser) Close() error {
	err := r.ReadCloser.Close()
	_ = r.cmd.Wait()
	return err
}
