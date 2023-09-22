---
title: 'Go Tidbit: Waiting on Go Routines'
date: 2023-09-17T21:15:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
toc: yes
toch2only: yes
---

Concurrency is one of the central features of Go. And, to build concurrent programs in Go, you need goroutines.

A goroutine is like a thread, but lighter. Much lighter. And, like any other built-in feature of Go, using it is dead simple:

``` go
package main

func main() {
  go func() {
    println("Hello World") // Print "Hello World" from a different goroutine.
  }()
}
```

Wait. That didn't print anything.

This is because when the `main` function of a Go program returns, it will abort all goroutines. Go will not wait until these goroutines finish running. But you can wait for goroutines to finish running.

In this blog post, we will take a look different ways of waiting for one or more goroutines to finish running.

## Using a Channel

If you want to wait on one goroutine, you could go the primitive route: use a channel.

Create a channel, of any type. From within the goroutine, at the very end of the function, send a value to the channel. And where you want to wait for the goroutine to finish, read from the channel.

The read will block until there is something in the channel to read.

``` go
package main

func main() {
  donech := make(chan int)

  go func() {
    println("Hello World")
    donech <- 1
  }()

  <- donech // The program will wait until it can read from this channel.
}
```

You could close the channel `donech` instead of sending a value over it to signal multiple readers that the goroutine has finished running.

If you want to wait on multiple goroutines, then read from the channel as many times as there are goroutines.

``` go
package main

func main() {
  const n = 7

  donech := make(chan int)

  for i := 0; i < n; i++ { 
    go func() {
      println("Hello World")
      donech <- 1
    }()
  }

  for i := 0; i < n; i++ { 
    <- donech // The program will wait until it can read from this channel.
  }
}
```

In the example above, you are reading from the channel `n` times. The loop, reading from the channel, can exit only when there are `n` values to read from the channel.

The standard library has a neat abstraction for this: `sync.WaitGroup`.

## Using `sync.WaitGroup`

The `sync.WaitGroup` type provides three methods:

- `Add(delta int)`: Add counts the number of goroutines we are waiting for.
- `Done()`: Done decrements the count by 1.
- `Wait()`: Wait blocks until the counter is zero.

Using a `sync.WaitGroup`, you will `Add` the number of goroutines to wait for. From within each goroutine, call `Done` right before the function returns.

Finally, call `Wait` from where you want to wait on the goroutines.

``` go
package main

import (
  "sync"
)

func main() {
  const n = 7

  wg := sync.WaitGroup{}

  wg.Add(n)
  for i := 0; i < n; i++ { 
    go func() {
      defer wg.Done()
      println("Hello World")
    }()
  }

  wg.Wait()
}
```

## Until One of Several Goroutines Fails (Using `errgroup.Group`)

If you have several goroutines and you want to stop as soon as any one of them experiences an error, then you could use Go's _almost standard_ package `golang.org/x/sync/errgroup`. 

This package provides the handy `errgroup.Group` type that can run multiple goroutines and return the first non-nil error, if any.

``` go
package main

import (
  "golang.org/x/sync/errgroup"
)

func main() {
  g := errgroup.Group{}
  for i := 0; i < 7; i++ {
    g.Go(func() error {
      err := mightReturnErr()
      return err
    })
  }
  err := g.Wait()
  if err == nil {
    println("All goroutines finished without error")
  }
}
```

### Aborting the Rest

If your program doesn't exit, the remaining goroutines will continue to run and consume resources. Depending on what these goroutines are doing you could use a cancellable `context.Context` to signal the remaining goroutines to be aborted.

## Until a Timeout

You can do wonderful things in Go with channels.

The `time.NewTimer` function returns a `time.Timer` with a channel. The channel is closed when the timer expires. By calling `time.NewTimer` with a duration, the channel will close after the given duration has elapsed.

By combining multiple channel operations using a `select` statement, you can wait on different conditions. For example, either wait on all goroutines to send values over a channel, or wait for a timer to expire.

``` go
package main

func main() {
  const n = 7

  donech := make(chan int)

  for i := 0; i < n; i++ { 
    go func() {
      println("Hello World")
      donech <- 1
    }()
  }

  t := time.After(10 * time.Second)
  for i := 0; i < n; i++ { 
    select {
    case <-donech: // The program will wait until it can read from this channel.
    case <-t.C:     // Or, until the timer expires.
    }
  }
}
```

In the code above, the second loop will exit when all goroutines have sent messages over `donech`. Or, the timer `t` has expired.

Similarly, you could wait until the user presses `Ctrl+C` (i.e. sends an interrupt from the terminal) by [using a signal channel](https://hjr265.me/blog/go-tidbit-handling-signals-exitting-gracefully/).

## Wrap Up

In Go, channels are the primitive for synchronization.

If you want to wait on a `goroutine` you have to use a channel; whether it is as a primitive or through a package.
