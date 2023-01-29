package ziface

import "context"

type IDatabase interface {
	Connect(ctx context.Context) (bool, error)
	Close(ctx context.Context) error
}
