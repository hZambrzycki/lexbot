package ports

import "context"

type FileHasher interface {
	HashFile(ctx context.Context, path string) (string, error)
}
