package storage

import (
	"context"
	"io"
)

type Storage interface {
	Download(ctx context.Context, key *string) (io.Reader, error)
	PreSign(ctx context.Context, key, contentType *string) (*string, error)
	ListObjects(ctx context.Context) ([]*File, error)
	Upload(ctx context.Context, key, contentType *string, data io.Reader) error
}
