---
title: 'Synchronization Constructs in the Go Standard Library'
date: 2023-10-26T15:30:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
toc: yes
---

Go provides `sync.Mutex` as its implementation of a mutual exclusion lock. However, it is not the only synchronization construct that is a part of the standard library.

This blog post will look at four synchronization constructs that we can use instead of a [`sync.Mutex`](https://pkg.go.dev/sync#Mutex).

## Counter

You may often see code using a `sync.Mutex` to synchronize access to a counter variable from multiple goroutines.

Like this:

``` go
var (
  n int
  m sync.Mutex
)
for i := 0; i < 7; i++ {
  go func() {
    m.Lock()
    defer m.Unlock()
    n++
  }()
}

// Elsewhere
m.Lock()
fmt.Println(n)
m.Unlock()
```

Instead of this, you can use an [`atomic.Int32`](https://pkg.go.dev/sync/atomic#Int32) or [`atomic.Int64`](https://pkg.go.dev/sync/atomic#Int64):

``` go
var n atomic.Int32
for i := 0; i < 7; i++ {
  go func() {
    n.Add(1)
  }()
}

// Elsewhere
fmt.Println(n.Load())
```

## Map

To synchronize a map, you will see code that uses a `sync.Mutex` or [`sync.RWMutex`](https://pkg.go.dev/sync#RWMutex).

You gain some benefits using a `sync.RWMutex` with read-heavy operations on the map.

``` go
var (
  s map[int]bool
  m sync.RWMutex
)

for i := 0; i < 7; i++ {
  go func(i int) {
    m.Lock()
    defer m.Unlock()
    s[i] = true
  }(i)
}

// Elsewhere
m.RLock()
for k, v := range s {
  fmt.Println(k, v)
}
m.RUnlock()
```

Instead of using a map-mutex pair, you can use [`sync.Map`](https://pkg.go.dev/sync#Map):

``` go
var s sync.Map

for i := 0; i < 7; i++ {
  go func(i int) {
    s.Store(i, true)
  }(i)
}

// Elsewhere
s.Range(func(k, v any) bool {
  fmt.Println(k.(int), v.(bool))
  return true
})
```

If you feel uneasy about the use of `any` in all `sync.Map` functions, you could define a generic wrapper:

``` go
type Map[K any, V any] struct {
  m sync.Map
}

func (m *Map[K, V]) Load(key K) (V, bool) {
  v, ok := m.m.Load(key)
  return v.(V), ok
}

func (m *Map[K, V]) Range(fn func(key K, value V) bool) {
  m.m.Range(func(key, value any) bool {
    return fn(key.(K), value.(V))
  })
}

func (m *Map[K, V]) Store(key K, value V) {
  m.m.Store(key, value)
}
```

And then use the wrapper instead:

``` go
var s Map[int, bool]

for i := 0; i < 7; i++ {
  go func(i int) {
    s.Store(i, true)
  }(i)
}

// Elsewhere
s.Range(func(k int, v bool) bool {
  fmt.Println(k, v)
  return true
})
```

One caveat is that the `Range` function is different from holding a lock around the `range` loop in the `sync.RWMutex` example. `Range` does not necessarily correspond to any consistent snapshot of the map's contents.

## Once Function

If you have a function called from multiple places in the code but you want it to run exactly once, you can do something like this:

``` go
var (
  called bool
  m sync.Mutex
)

func DoSomethingOnlyOnce() {
  m.Lock()
  defer m.Unlock()
  if called {
    return
  }

  called = true

  // ...
}
```

Or, better, you can use [`sync.Once`](https://pkg.go.dev/sync#Once):

``` go
var once sync.Once

func DoSomethingOnlyOnce() {
  once.Do(func() {
    // ...
  })
}
```

If your function is meant to return one or two values, you can use [`sync.OnceValue`](https://pkg.go.dev/sync#OnceValue) or [`sync.OnceValues`](https://pkg.go.dev/sync#OnceValues).

## Condition

Let's say you want to make several goroutines wait for a condition to be met. You can, but shouldn't, use a `sync.Mutex` for it:

``` go
var (
  condition bool
  m sync.Mutex
)

for i := 0; i < 7; i++ {
  go func() {
    for {
      m.Lock()
      if condition {
        m.Unlock()
        break
      }
      m.Unlock()
    }
    // ...
  }()
}

// Elsewhere
m.Lock()
condition = true
m.Unlock()
```

This code is very taxing on the CPU.

The standard library has [`sync.Cond`](https://pkg.go.dev/sync#Cond) for this:

``` go
var (
  condition bool
  c = sync.NewCond(sync.Mutex{})
)

for i := 0; i < 7; i++ {
  go func() {
    c.L.Lock()
    for !condition {
      c.Wait()
    }
    c.L.Unlock()
    // ...
  }()
}

// Elsewhere
c.L.Lock()
condition = true
c.L.Unlock()
c.Broadcast()
```

Unlike the example using `sync.Mutex`, this code is not taxing on the CPU. `Wait` returns for all goroutines only when `Broadcast` is called.

However, in a simple case like this, where you will signal only when the condition has been met, you can close a channel to signal to multiple goroutine:

``` go
var c chan struct{}

for i := 0; i < 7; i++ {
  go func() {
    <-c
    // ...
  }()
}

// Elsewhere, only after the condition has been met
close(c)
```

Note that in the case of using a channel to signal the goroutines you cannot reopen the channels. This approach works if you expect the condition to meet and not be changed anymore.
