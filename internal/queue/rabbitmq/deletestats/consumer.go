package deletestats

import "context"

func WrapClickhouseDelete(fn func(shortCode string) error) func(ctx context.Context, shortCode string) error {
	return func(ctx context.Context, shortCode string) error {
		return fn(shortCode)
	}
}
