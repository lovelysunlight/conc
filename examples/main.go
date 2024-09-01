package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lovelysunlight/conc/pool"
)

func main() {
	p := pool.New(context.Background(), pool.WithMaxGoroutines(5))
	for i := 0; i < 10; i++ {
		p.Go(func(context.Context) error {
			if i >= 0 && i <= 8 {
				time.Sleep(3 * time.Second)
			}
			fmt.Printf("hello world from %d \n", i)
			return nil
		})
	}

	_ = p.Wait()
}
