# WPool

[![Code Climate](https://codeclimate.com/github/ustrajunior/wpool/badges/gpa.svg)](https://codeclimate.com/github/ustrajunior/wpool)
[![Build Status](https://travis-ci.org/ustrajunior/wpool.svg?branch=master)](https://travis-ci.org/ustrajunior/wpool)


WPool it's a simple background worker pool for Go. You can send tasks to be executed concurrently.

## Instalation

Just run

```
go get -u github.com/ustrajunior/wpool
```

## How to use

You have to start the pool and choose how many workers you want.

```go
wp := wpool.New(200)
```

The pool accepts functions with this signature: **func() error**. To add a new task to the pool, you can use the Add function like this

```go
wp.Add(func() error {
	// same expensive task
})
```

Now, we need to wait for all tasks to finish, so:

```go
wp.Wait()
```

With this simple 3 steps, you can run functions concurrently.


### "Full" example

```go
package main

import (
	"fmt"
	"time"
	"github.com/ustrajunior/wpool"
)

func main() {
	wp := wpool.New(200, 100)

	for i := 0; i < 1000; i++ {
		a := i
		wp.Add(func() error {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Running task:  %v\n", a+1)
			return nil
		})
	}
	
	wp.Wait()
}
```

## License

Copyright (c) 2017-present [José Carlos Ustra Júnior](https://github.com/ustrajunior)

Licensed under [MIT License](https://github.com/ustrajunior/wpool/blob/master/LICENSE)
