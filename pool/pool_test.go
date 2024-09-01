package pool

import (
	"context"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		g := New(context.Background())
		var completed atomic.Int64
		for i := 0; i < 100; i++ {
			g.Go(func(context.Context) error {
				time.Sleep(1 * time.Millisecond)
				completed.Add(1)
				return nil
			})
		}
		_ = g.Wait()
		require.Equal(t, completed.Load(), int64(100))
	})

	t.Run("panics on configuration after init", func(t *testing.T) {
		t.Run("before wait", func(t *testing.T) {
			t.Parallel()
			g := New(context.Background())
			g.Go(func(context.Context) error { return nil })
		})

		t.Run("after wait", func(t *testing.T) {
			t.Parallel()
			g := New(context.Background())
			g.Go(func(context.Context) error { return nil })
			_ = g.Wait()
		})
	})

	t.Run("limit", func(t *testing.T) {
		t.Parallel()
		for _, maxConcurrent := range []int{1, 10, 100} {
			t.Run(strconv.Itoa(maxConcurrent), func(t *testing.T) {
				g := New(context.Background(), WithMaxGoroutines(maxConcurrent))

				var currentConcurrent atomic.Int64
				var errCount atomic.Int64
				taskCount := maxConcurrent * 10
				for i := 0; i < taskCount; i++ {
					g.Go(func(context.Context) error {
						cur := currentConcurrent.Add(1)
						if cur > int64(maxConcurrent) {
							errCount.Add(1)
						}
						time.Sleep(time.Millisecond)
						currentConcurrent.Add(-1)
						return nil
					})
				}
				_ = g.Wait()
				require.Equal(t, int64(0), errCount.Load())
				require.Equal(t, int64(0), currentConcurrent.Load())
			})
		}
	})

	t.Run("panics on invalid WithMaxGoroutines", func(t *testing.T) {
		t.Parallel()
		require.Panics(t, func() { New(context.Background(), WithMaxGoroutines(0)) })
	})

	t.Run("is reusable", func(t *testing.T) {
		t.Parallel()
		var count atomic.Int64
		p := New(context.Background())
		for i := 0; i < 10; i++ {
			p.Go(func(context.Context) error {
				count.Add(1)
				return nil
			})
		}
		_ = p.Wait()
		require.Equal(t, int64(10), count.Load())
		for i := 0; i < 10; i++ {
			p.Go(func(context.Context) error {
				count.Add(1)
				return nil
			})
		}
		_ = p.Wait()
		require.Equal(t, int64(20), count.Load())
	})
}
