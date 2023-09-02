---
title: "Configuring Docker Health Check for Go Web Applications"
date: 2023-09-02T12:45:00+06:00
tags:
  - Docker
  - Go
  - 100DaysToOffload
  - HealthCheck
---

Docker has been providing a health check mechanism for quite some time. It is useful in identifying issues with programs that can fail in ways other than just outright crashing.

And it is easy to set up.

Docker health checks work periodically running a program within the container and observing its exit status. If it exits with a 0, the container is considered healthy. If it exits with a 1, the container is considered to be unhealthy.

You can configure health check within the Dockerfile or when starting/creating a container using `docker run`/`docker create`.

Let's begin with a two-staged build+run Dockerfile for a simple Go web application:

``` Dockerfile
# Build Stage
FROM golang:1.21.0-alpine3.18 AS builder

# Build the Go web application.
WORKDIR /hellohealth
COPY . .
RUN go build .

# Run Stage
FROM alpine:3.18.3

# We need curl to hit the health check endpoint.
RUN apk --no-cache add curl

# Copy the built Go web application.
WORKDIR /hellohealth
COPY --from=builder /hellohealth /hellohealth

# Use curl to hit the healthcheck endpoint. If the hit suceeds exit with 0, 1 otherwise.
HEALTHCHECK --interval=30s --timeout=3s CMD curl -f http://127.0.0.1:8080/api/health_checks/ready || exit 1

EXPOSE 8080
ENTRYPOINT ["./hellohealth"]
```

Within your Go program, handle the `/api/health_checks/ready` HTTP path:

``` go
http.HandleFunc("/api/health_checks/ready", func(w http.ResponseWriter, r *http.Request) {
  // Check if everything is okay.
  // - Ping databases
  // - Ping message brokers
  // - etc.

  // If everything is okay...
  io.WriteString(w, "READY")
})
```

You can now build and run the Docker container:

``` sh
docker build -t hellohealth:latest .
docker run -P hellohealth
```

Alternatively, you can configure health check when running the container (instead of in the Dockerfile):

``` sh {linenos=false}
docker run -P --health-interval=30s --health-timeout=3s --health-cmd curl -f http://127.0.0.1:8080/api/health_checks/ready || exit 1 hellohealth
```

Either way, checking the status of the container you will notice that it is now reporting its health as "starting".

``` txt {linenos=false}
» docker ps
CONTAINER ID   IMAGE         COMMAND           CREATED         STATUS                            PORTS                                         NAMES
08b127b45803   hellohealth   "./hellohealth"   4 seconds ago   Up 2 seconds (health: starting)   0.0.0.0:32769->8080/tcp, :::32769->8080/tcp   elegant_colden
```

But after a few seconds it will change to "healthy".

``` txt {linenos=false}
» docker ps
CONTAINER ID   IMAGE         COMMAND           CREATED          STATUS                    PORTS                                         NAMES
08b127b45803   hellohealth   "./hellohealth"   37 seconds ago   Up 35 seconds (healthy)   0.0.0.0:32769->8080/tcp, :::32769->8080/tcp   elegant_colden
```

A somewhat complete example of this tutorial is available in [this GitHub repository](https://github.com/hjr265/hellohealth).
