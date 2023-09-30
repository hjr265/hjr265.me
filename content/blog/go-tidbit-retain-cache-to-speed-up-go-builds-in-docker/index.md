---
title: 'Go Tidbit: Retain Cache to Speed up Go Builds in Docker'
date: 2023-09-29T20:20:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

If you are building Go binaries inside a Docker container, you can speed up the builds by retaining the cache that `go build` creates between builds.

{{< image src="splash.png" alt="Speed up Go Builds in Docker" >}}

This tidbit is especially effective for large projects with several packages.

Let's take one of my Go projects as an example.

Toph Platform, the Go project that powers the web application at [toph.co](https://toph.co/) has 55 packages inside the repository. If I include dependencies (first and third-party), the total number of packages comes to 813.

Building this project on my laptop takes more than a few seconds.

``` sh
time docker run --rm -t -e GOOS -e GOARCH -v `pwd`:/toph/platform -w /toph/platform golang:1.21.1 go build -v ./cmd/tophd
```

With the `-v` flag, `go build` lists all the packages it has to build.

``` txt
internal/goarch
internal/race
internal/goos
[... 784 lines]
[...]/toph/platform/booklet
[...]/toph/platform/sendy
[...]/toph/platform/api
[...]/toph/platform/ui
[...]/toph/platform/cmd/tophd
0.01s user 0.02s system 0% cpu 30.643 total
```

That's 30 seconds for a build.

Go is building each package, including the dependencies, every time. That isn't ideal.

But by retaining the cache directory created during a `go build`, you can speed up subsequent builds significantly. With the cache present, Go will rebuild a package only if it changes.

And, let's face it, it is unlikely that I will be changing all 55 or 813 packages between builds while I am working on the project.

To retain the cache, bind a host directory to the cache directory within the container.

The path to the cache directory is known to the `go` command through the `GOCACHE` environment variable. Run `go env` in the Docker container to know where Go will be caching its build artifacts.

``` sh
docker run --rm -t golang:1.21.1 go env
```

``` txt
...
GOCACHE='/root/.cache/go-build'
...
```

With that, you can now run the Docker container with an additional flag: `-v /tmp/toph-platform-go-build:/root/.cache/go-build`

The first part, `/tmp/toph-platform-go-build`, is the directory where you want to retain the contents of the cache. The second part, `/root/.cache/go-build`, is where Go will cache build artifacts within the container.

The first build will take time as the cache will be empty. But each subsequent build will be so much faster.

``` sh
time docker run --rm -t -e GOOS -e GOARCH -v `pwd`:/toph/platform -v /tmp/toph-platform-go-build:/root/.cache/go-build -w /toph/platform golang:1.21.1 go build -v ./cmd/tophd
```

``` txt
[...]/toph/platform/ui
[...]/toph/platform/cmd/tophd
0.01s user 0.02s system 0% cpu 3.677 total
```

You can leverage the same technique in your CI/CD configuration to speed up your Go build pipelines.
