---
title: 'Makefile Recipe: Build with Docker and Export To A Gzipped Tarball'
date: 2023-08-28T17:30:00+06:00
tags:
  - Makefile
  - Docker
  - 100DaysToOffload
---

I have been writing a lot of Makefiles lately. I find them simple and easy to like. And, as with all old-school things, I am starting to overlook its quirks.

I needed to write a Makefile that lets me build something with Docker and then export the entire contents of the image to a gzipped tarball.

The first part is easy:

``` Makefile
.PHONY: build
build:
  (cd $(APP); docker build -t thing-$(APP):latest .)
```

When run as `make build APP=something`, it will build the `Dockerfile` within the `something/` directory.

But the next part is where it gets a bit interesting.

To export the contents of the image, as if it were a root filesystem, into a tarball, we need to create a container first. We need to run a few `docker` commands back to back, so we need to know the container name (or ID).

In Makefile:

``` Makefile
.PHONY: tarball
tarball:
  TS=$$(date +'%s') && \
  docker create --name=tmp-$(APP)-$$TS tmp-$(APP):latest && \
  docker export tmp-$(APP)-$$TS | gzip > $(APP).tar.gz && \
  docker rm tmp-$(APP)-$$TS
```

This script will first set the current Unix timestamp to `TS` and then:

- Create a Docker container named "tmp-{APP}-{TS}", where `{APP}` and `{TS}` will be interpolated.
- Export the contents of the Docker container (which is a tarball), pipe it through `gzip` and then save it as a file named `{APP}.tar.gz`.
- Finally, it will remove the Docker container.

The script could be better. But it is a good start for what I needed. And, as with all things, it can be improved as and when necessary.
