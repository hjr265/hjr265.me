---
title: "Tracking io.Copy Progress in Go"
date: 2023-11-12T09:45:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

If you are writing Go code for any period, you must have used the `io.Copy` function. It takes an `io.Writer` and an `io.Reader` and copies everything from the reader to the writer until it reaches the end of file (EOF).

The function returns the number of bytes copied and an error (if any, other than `io.EOF`).

But this function blocks until the copy completes. How do you track the progress of `io.Copy`?

You can write a custom `io.Writer` wrapper:

``` go
import (
  "io"
  "sync/atomic"
)

type ProgressWriter struct {
  w io.Writer
  n atomic.Int64
}

func NewProgressWriter(w io.Writer) *ProgressWriter {
  return ProgressWriter{w: w}
}

func (w *ProgressWriter) Write(b []byte) (n int64, err error) {
  n, err = w.Write(b)
  w.n.Add(n)
  return
}

func (w *ProgressWriter) N() int64 {
  return w.n.Load()
}
```

The `ProgressWriter` wraps any `io.Writer` and proxies all `Write` calls to it. As it proxies the call, it tracks the number of bytes written to the `io.Writer`.

The method `N` can be called on the `ProgressWriter` to read the total number of bytes written to the `io.Writer`. It is safe to call `N` from different or multiple goroutines since the `ProgressWriter` stores and reads the total number of bytes from an `atomic.Int64`.

``` go
// Error handling omitted for brevity.

var (
  w io.Writer
  r io.Reader
)

pw := NewProgressWriter(w)

go io.Copy(pw, r)

go func() {
  for {
    time.Sleep(1 * time.Second)
    fmt.Printf("Copied %d bytes\n", pw.N())
  }
}()
```

The second goroutine will print the total number of bytes copied to the `io.Writer` every second.
