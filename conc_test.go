package conc

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWaitGroup(t *testing.T) {
	t.Parallel()

	t.Run("all spawned run", func(t *testing.T) {
		t.Parallel()
		var count atomic.Int64

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := New(ctx)
		for i := 0; i < 100; i++ {
			wg.Go(func(context.Context) error {
				count.Add(1)
				return nil
			})
		}
		_ = wg.Wait()
		require.Equal(t, count.Load(), int64(100))
	})

	t.Run("all spawned return err", func(t *testing.T) {
		t.Parallel()

		var count atomic.Int64
		wg := New(context.Background())
		for i := 0; i < 100; i++ {
			wg.Go(func(context.Context) error {
				return fmt.Errorf("error: %d", count.Add(1))
			})
		}
		err := wg.Wait()
		require.Equal(t, count.Load(), int64(100))
		require.Error(t, err)

		_ = wg
		runtime.GC()
	})
}
