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
wp := wpool.New(200, 200)
```

The pool accepts functions with this signature: **func() error**. To add a new task to the pool, you can use the Add function and add a wpool.Job like this:

```go
fn := func() error {
	// same expensive task
}

wp.Add(wpool.Job{Task: fn, Retry: 5})
```

When an error occurs in the Task, wpool will try to execute again until the limit determined in Retry option.

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

	for i := 0; i < 10; i++ {
		a := i
		wp.Add(wpool.Job{Retry: 5, Task: func() error {
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("Running task:  %v\n", a+1)
			return nil
		}})
	}
	
	wp.Wait()
}
```

You still need to think about race conditions. If your task access or modify the same variable, you can have a problem.

## License

Copyright (c) 2017-present [José Carlos Ustra Júnior](https://github.com/ustrajunior)

Licensed under [MIT License](https://github.com/ustrajunior/wpool/blob/master/LICENSE)
