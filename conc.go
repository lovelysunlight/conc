package conc

import (
	"context"
	"errors"
	"runtime"
	"sync"
)

type WaitGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	errch  chan error
	errors []error
}

type Task func(context.Context) error

func (g *WaitGroup) Go(f Task) {
	g.wg.Add(1)
	defer g.wg.Done()

	if err := f(g.ctx); err != nil {
		g.errch <- err
	}
}

func (g *WaitGroup) Wait() error {
	g.wg.Wait()

	return errors.Join(g.errors...)
}

func (g *WaitGroup) background() {
	for {
		select {
		case err := <-g.errch:
			g.errors = append(g.errors, err)
		case <-g.ctx.Done():
			return
		}
	}
}

func New(ctx context.Context) *WaitGroup {
	ctx, cancel := context.WithCancel(ctx)
	wg := &WaitGroup{
		ctx:    ctx,
		cancel: cancel,
		errch:  make(chan error),
		errors: make([]error, 0),
	}
	go wg.background()
	runtime.SetFinalizer(wg, func(wg *WaitGroup) {
		wg.cancel()
	})

	return wg
}
