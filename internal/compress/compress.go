package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/klauspost/compress/zstd"
)

type Algorithm string

const (
	Gzip Algorithm = "gzip"
	Zstd Algorithm = "zstd"
	None Algorithm = "none"
)

// AlgorithmFromKey infers the compression algorithm from a backup key's extension.
// This is the correct way to pick a decompressor on restore — not from current config,
// which may have changed since the backup was made.
func AlgorithmFromKey(key string) Algorithm {
	switch {
	case strings.HasSuffix(key, ".gz"):
		return Gzip
	case strings.HasSuffix(key, ".zst"):
		return Zstd
	default:
		return None
	}
}

// Wrap returns a compressing WriteCloser around w.
func Wrap(w io.Writer, algo Algorithm) (io.WriteCloser, error) {
	switch algo {
	case Gzip:
		return gzip.NewWriterLevel(w, gzip.BestCompression)
	case Zstd:
		return zstd.NewWriter(w)
	case None:
		return nopWriteCloser{w}, nil
	default:
		return nil, fmt.Errorf("unexpected compression algorithm: %q", algo)
	}
}

// Unwrap returns a decompressing Reader around r.
func Unwrap(r io.Reader, algo Algorithm) (io.Reader, error) {
	switch algo {
	case Gzip:
		return gzip.NewReader(r)
	case Zstd:
		zr, err := zstd.NewReader(r)
		if err != nil {
			return nil, err
		}
		return zr, nil
	case None:
		return r, nil
	default:
		return nil, fmt.Errorf("unexpected compression algorithm: %q", algo)
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
