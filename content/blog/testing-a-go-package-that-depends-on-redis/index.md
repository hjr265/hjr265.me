---
title: "Testing a Go Package That Depends on Redis"
date: 2023-11-13T13:45:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

Redsync, one of my open-source Go packages, implements a distributed lock using Redis. It is an implementation of the [Redlock algorithm](https://redis.io/docs/manual/patterns/distributed-locks/).

This Go package has tests that run against multiple real Redis servers. And it is an example of how you can use the `TestMain` function to customize your Go tests.

The [`TestMain` function](https://pkg.go.dev/testing#hdr-Main), if defined in your Go tests, will allow you to run the Go code before and after the actual tests, which is really ideal for setup and teardown work.

In the case of Redsync, the `TestMain` function is where multiple Redis servers spawn before running the tests. A Go package, [`tempredis`](https://github.com/stvp/tempredis), makes it easy to start Redis servers without having to muck around with `exec.Command`.

``` go
var (
  nServers = 10
  servers  []*tempredis.Server
)

func TestMain(m *testing.M) {
  // Spawn Redis servers.
  for i := 0; i < nServers; i++ {
    server, err := tempredis.Start(tempredis.Config{})
    if err != nil {
      panic(err)
    }
    servers = append(servers, server)
  }

  // Run the tests.
  code := m.Run()

  // Cleanup
  for _, server := range servers {
    server.Term()
  }

  // Exit with code reported by tests.
  os.Exit(code)
}
```

When any tests run in this package, `servers` will have instances of `*tempredis.Server`.

``` go
func TestPing() {
  for _, srv := range servers {
    rdb := redis.NewClient(&redis.Options{
      Network:  "unix",
      Addr:     srv.Socket(),
    })

    err := rdb.Ping(context.Background()).Err()
    if err != nil {
      panic(err)
    }
  }
}
```
