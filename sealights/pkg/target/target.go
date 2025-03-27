package target

import (
	"context"
)

type Target interface {
	Do(ctx context.Context, request TargetRequest) (TargetResponse, error)
	Close() error
}
