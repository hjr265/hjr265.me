---
title: 'Building Multi-platform Docker Image for Go Applications'
date: 2023-09-10T21:20:00+06:00
tags:
  - Docker
  - Buildx
  - 100DaysToOffload
  - Go
---

I began the weekend writing a silly program: MGHSIAC. It's the elegant abbreviation of "My GitHub Status Is A Clock". It turns my GitHub status into a working clock. You can read more about it [here](/blog/my-github-status-is-a-clock/).

But as silly as it is, I am now committed to keep it running.

I have an always-on Raspberry Pi with Portainer running on it already. If I could make a Docker image and host it on Docker Hub, I could easily pull it to that Raspberry Pi and have it continuously update my GitHub status with clock emojis and messages.

## Dockerfile

I wrote a two-stage Dockerfile. The first stage builds the Go application and the second stage runs it.

``` dockerfile
# Build Stage
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.0-alpine3.18 AS builder

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Build the Go application.
WORKDIR /mghsiac
ADD . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" .

# Run Stage
FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.18.3

RUN apk add --no-cache tzdata

# Copy the built Go application.
WORKDIR /mghsiac
COPY --from=builder /mghsiac/mghsiac /mghsiac/mghsiac

ENTRYPOINT ["./mghsiac"]
```

I just had to add a few lines to it to make it work across multiple platforms:

- The four `ARG` lines define the arguments that `docker buildx` passes.
- We use the `BUILDPLATFORM` argument to choose the correct build image. If I am building on Linux on AMD64, I want to pull the AMD64 Alpine.
- We also pass the `TARGETOS` and `TARGETARCH` values as `GOOS` and `GOARCH` respectively. This way when building for Linux on ARM64, we can tell the Go compiler to do that.
- We use the `TARGETPLATFORM` argument to choose the correct image for the run stage.

And that's pretty much all there is to it.

You can now build Docker images for different platforms like so:

- `docker buildx build --platform=linux/amd64 .`
- `docker buildx build --platform=linux/arm64 .`
- `docker buildx build --platform=win/amd64 .`
- `docker buildx build --platform=darwin/amd64 .`
- `docker buildx build --platform=darwin/arm64 .`
- etc.

### Example Build

Let's work through an example.

If we run `docker buildx build --platform=linux/arm64 .` on a `linux/amd64` (notice `arm64` vs `amd64`) machine, we are telling Docker to build the Go program with `TARGETPLATFORM` set to `linux/arm64` and `TARGETOS` set to `linux` and `TARGETARCH` set to `arm64`.

Since we are running on `linux/amd64`, the argument `BUILDPLATFORM` is `linux/amd64`.

At this point, the build stage will use an `amd64` variant of the Alpine image. This is because we set `--platform=${BUILDPLATFORM:-linux/amd64}` with the with the first `FROM`. (The extra bit of `:-linux/amd64` within the curly braces is just a way to provide a fallback value for this variable).

The build stage will pass the `GOOS` and `GOARCH` environment variables set to `linux` and `arm64`. That's cross-compiling that the Go compiler can already deal with.

Unlike the build stage, the run stage here will use the `arm64` variant of the Alpine image.

## Further Reading

The `docker buildx` command can do more than just build multi-platform images. And it is worth reading up on those details:

- https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/
- https://docs.docker.com/build/building/multi-platform/
