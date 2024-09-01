# Go Library Template

[![Build Badge]][build status]
[![Go Reference]][godoc]
[![License Badge]][license]

A template repository for Go Library.

## Installing

```shell
$ go get github.com/lovelysunlight/conc
```

## Usage

Use `WaitGroup` to run goroutine.
```golang
import (
	"context"
	"fmt"
	"time"

	"github.com/lovelysunlight/conc/pool"
)

func main() {
    wg := conc.New(context.Background())
    for i := 0; i < 10; i++ {
        wg.Go(func(ctx context.Context) error {
            if i >= 0 && i <= 8 {
                time.Sleep(3 * time.Second)
            }
            fmt.Printf("hello world from %d \n", i)
            return nil
        })
    }
    _ = wg.Wait()
}
```

Use `pool.New` to create a goroutine pool.
```golang
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
```

## Contributing

You can contribute in one of three ways:

1. File bug reports using the [issue tracker](https://github.com/lovelysunlight/conc/issues).
2. Answer questions or fix bugs on the [issue tracker](https://github.com/lovelysunlight/conc/issues).
3. Contribute new features or update the wiki.

## License

MIT

[build badge]: https://github.com/lovelysunlight/conc/actions/workflows/ci.yaml/badge.svg
[build status]: https://github.com/lovelysunlight/conc/actions/workflows/ci.yaml
[go reference]: https://pkg.go.dev/badge/github.com/lovelysunlight/conc?status.svg
[godoc]: https://pkg.go.dev/github.com/lovelysunlight/conc?tab=doc
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[license]: https://raw.githubusercontent.com/lovelysunlight/conc/master/LICENSE