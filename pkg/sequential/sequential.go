package sequential

import (
	"context"
	"fmt"
)

func MustDo[R any](ctx context.Context, n int, handleFunc func(context.Context, int) (R, error)) (r R, i int, err error) {
	if n == 0 {
		err = fmt.Errorf("#clients == 0")
		return
	}

	if n == 1 {
		for {
			r, err = handleFunc(ctx, 0)
			if err == nil {
				return
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}

	for j := 0; ; j++ {
		i = j % n
		r, err = handleFunc(ctx, i)
		if err == nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
