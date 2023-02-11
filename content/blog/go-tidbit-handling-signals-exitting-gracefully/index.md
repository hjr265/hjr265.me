---
title: 'Go Tidbit: Handling Signals, Exiting Gracefully'
date: 2023-02-09T20:25:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
toc: yes
---

Signals are standardized messages that an operating system can send your programs.

Take `Ctrl+C` for example. When running a program from the terminal and you hit `Ctrl+C`, you expect the program to end immediately.

How does that work, though? `Ctrl+C` is a _shortcut_ for the POSIX signal `SIGINT`. By default, this signal causes your program to be terminated.

But this is one of those signals you can handle: You can intercept it and do whatever you please.

## Handling Signals in Go

The following snippet does the bare minimum in Go to handle a signal.

``` go
package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("Waiting for signal.")

	// Make a buffered channel.
	sigch := make(chan os.Signal, 1)
	// Register the signals that you want to handle.
	signal.Notify(sigch, os.Interrupt)

	// Wait for the signal.
	<-sigch

	fmt.Println("Received interrupt. Exiting.")
}
```

If you run this program, you will see the message "Waiting for signal.". You can then press `Ctrl+C` to trigger a SIGINT. The program will print "Received interrupt. Exiting." and exit.

## Exitting Go Programs Gracefully

This signal and channel craft now open your Go programs to many possibilities.

### Flushing Writes Before Exitting

Imagine your Go program is working with large files, and you want to flush all writes before the program exits.

``` go
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

func main() {
	// This program creates a numbers.txt file and populates it with a lot of
	// numbers, one per line. Errors omitted for brevity.

	// Create a file and wrap it in a buffered writer.
	f, _ := os.Create("numbers.txt")
	bw := bufio.NewWriter(f)

	// Close this channel to stop the number generator.
	stopch := make(chan struct{})

	// Use a WaitGroup to wait for the number generator to stop.
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Run the number generator in a separate Go routine.
	go func() {
		defer wg.Done()
	L:
		for i := 0; i < 1000000000; i++ {
			fmt.Fprintf(bw, "%d\n", i)
			select {
			case <-stopch:
				// The stopch channel has been closed. Break the loop.
				break L
			default:
			}
		}
	}()

	// Make a signal channel. Register SIGINT.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	// In a separate Go routine, wait for the signal. On signal, close the
	// stopch channel so that the number generator stops.
	go func() {
		<-sigch
		fmt.Println("Interrupted.")
		close(stopch)
	}()

	// Wait for the number generator to stop.
	wg.Wait()

	fmt.Println("Flushing.")

	// Flush buffered writer. Close file.
	bw.Flush()
	f.Close()

	fmt.Println("Exiting.")
}
```

In this program, instead of allowing the default `SIGINT` behaviour, it stops the generator. The main function then flushes the buffered writer and closes the file.

### Shutting Down Go HTTP Server Gracefully

Another example is where you can shut down a Go HTTP server gracefully.

``` go
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Errors omitted for brevity.

	// Make an HTTP server.
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate a slow HTTP response.
			time.Sleep(10 * time.Second)
			io.WriteString(w, "Hello")
		}),
	}

	// Start the HTTP server in a separate Go routine.
	go func() {
		fmt.Println("Listening for HTTP connections.")
		server.ListenAndServe()
	}()

	// Make a signal channel. Register SIGINT.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	// Wait for signal.
	<-sigch

	fmt.Println("Interrupted. Exiting.")

	// Trigger a shutdown and allow 13 seconds to drain connections. Ignoring
	// CancelFunc for brevity.
	ctx, _ := context.WithTimeout(context.Background(), 13*time.Second)
	server.Shutdown(ctx)
}
```

In this example, the HTTP server simulates a slow response. It takes 10 seconds before sending the response.

To try out this example, start this program and navigate to https://localhost:8080 on your web browser. Then switch back to the terminal immediately and press `Ctrl+C`. You will see that the program waits until your response is served before exiting.

In case there were no pending requests, the program would exit immediately.

## Forcing Exit on Second Interrupt

When writing Go programs for other people to run on their computers, allow the program to exit immediately after receiving a second `SIGINT`.

This gives users control when they want to skip the graceful exit behaviour.

You can do this by starting up a Go routine right after receiving the first interrupt. In this Go routine, you can wait on the signal channel again and exit the program after receiving a signal.

``` go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Make a signal channel. Register SIGINT.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	// Wait for the signal.
	<-sigch
	go func() {
		fmt.Println("Interrupt again to force exit.")
		// Wait for a second signal.
		<-sigch
		os.Exit(1)
	}()

	fmt.Println("Interrupted. Exiting.")

	// Long clean-up code goes here.
	time.Sleep(5 * time.Second)
}
```

## Simplifying With Context

The `signal.Notify` also comes in the context variation as `signal.NotifyContext`.

When the functions of your Go program support context, you can simplify signal handling by copying your context with `signal.NotifyContext`. This copy of the context is marked as done as soon as one of the signals you register arrives.

The following Go code logically brings the same behaviour as the previous example under [Forcing Exit on Second Interrupt](#forcing-exit-on-second-interrupt) but with fewer lines of code.

``` go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Make a signal-based context. The stop function, when called, unregisters
	// the signals and restores the default signal behaviour.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Wait for the signal.
	<-ctx.Done()
	stop() // After calling stop, another SIGINT will terminate the program.

	fmt.Println("Interrupted. Exiting.")

	// Long clean-up code goes here.
	time.Sleep(5 * time.Second)
}
```

## Signals You Cannot Handle

There are signals, like `SIGKILL` and `SIGSTOP`, that you cannot handle. `SIGKILL`, for example, will terminate your program, no questions asked.

## Further Reading

Go's default behaviour for various signals is documented here: https://pkg.go.dev/os/signal

And you can learn more about all the POSIX signals here: https://man7.org/linux/man-pages/man7/signal.7.html

<br>

_This post is 17th of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
