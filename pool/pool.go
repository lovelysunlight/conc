package pool

import (
	"context"

	"github.com/lovelysunlight/conc"
)

type pool struct {
	limiter limiter
	wg      *conc.WaitGroup
}

func (p *pool) Go(f conc.Task) {
	if p.limiter != nil {
		p.limiter <- struct{}{}
	}
	p.addTask(f)
}

func (p *pool) addTask(f func(context.Context) error) {
	p.wg.Go(func(ctx context.Context) error {
		defer p.limiter.release()
		return f(ctx)
	})
}

func (p *pool) Wait() error {
	return p.wg.Wait()
}

type poolOption func(*pool)

func New(ctx context.Context, opts ...poolOption) *pool {
	p := &pool{
		wg: conc.New(ctx),
	}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithMaxGoroutines(n int) poolOption {
	return poolOption(func(p *pool) {
		if n < 1 {
			panic("max goroutines in a pool must be greater than zero")
		}

		p.limiter = make(limiter, n)
	})
}

type limiter chan struct{}

func (l limiter) release() {
	if l != nil {
		<-l
	}
}
