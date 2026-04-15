package ports

import "context"

type TextExtractor interface {
	ExtractText(ctx context.Context, path string) (string, error)
}
