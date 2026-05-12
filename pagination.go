package hhru

import "context"

func PagesUntil(ctx context.Context, startPage int, fn func(ctx context.Context, page int) (continueNext bool, err error)) error {
	for page := startPage; ; page++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		more, err := fn(ctx, page)
		if err != nil {
			return err
		}
		if !more {
			return nil
		}
	}
}
