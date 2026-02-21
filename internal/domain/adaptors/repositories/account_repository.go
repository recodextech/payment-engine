package repositories

import (
	"context"
)

type AccountRepository interface {
	Exists(ctx context.Context, key string) (exists bool, err error)
}
