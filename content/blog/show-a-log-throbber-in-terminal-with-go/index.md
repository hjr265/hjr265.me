---
title: 'Show a Log Throbber in Terminal with Go'
date: 2023-10-19T10:00:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
toc: yes
---

Show a Log Throbber in the Terminal with Go

A long-running program made to run in a terminal window should indicate what it is doing. Judicious logging is the first step.

While developing [Printd](https://github.com/FurqanSoftware/toph-printd) for Toph, we needed a way to indicate the program status without outputting loglines repeatedly.

Printd, a print server daemon, waits for print requests from Toph and prints out the contents of the request to a connected printer. Until these requests arrive, Printd is mostly sitting idle.

Printd uses a throbbing indicator at the end of the log lines to indicate its status. For example, while waiting for new requests, it shows `[~] Ready` with the `~` flashing slowly.

{{< video src="printd.webm" >}}

Printd uses [github.com/FurqanSoftware/pog](https://github.com/FurqanSoftware/pog) to implement multiple log throbbers in the terminal.

In this blog post, I only explain how Pog uses control escape sequences to implement these throbbers.

## Control Escape Sequence

Terminals support control escape sequences that allow you to move the cursor back to the start of a line and clear out its contents.

- `^[[2K`: Clear entire line
- `\r`: Move cursor to the beginning of the line

By combining these two, you can implement a throbber:

``` go
package main

// Imports omitted for brevity.

func main() {
  tick := 0
  for {
    time.Sleep(1 * time.Second)
    if tick == 0 {
      fmt.Print("\33[2K\r", "[ ] Waiting")
    } else {
      fmt.Print("\33[2K\r", "[~] Waiting")
    }
    tick = 1 - tick
  }
}
```

It will look something like this:

{{< video src="basic.webm" muted="true" >}}

But this is being output to `stdout`. To output this to `stderr`, you need to use `fmt.Fprint` with `os.Stderr` instead:

``` go
package main

// Imports omitted for brevity.

func main() {
  tick := 0
  for {
    time.Sleep(1 * time.Second)
    if tick == 0 {
      fmt.Fprint(os.Stderr, "\33[2K\r", "[ ] Waiting")
    } else {
      fmt.Fprint(os.Stderr, "\33[2K\r", "[~] Waiting")
    }
    tick = 1 - tick
  }
}
```

However, this now poses a new issue. It is interfering with log lines output to `stderr`. To fix that, you can set your logging library to clear the current line before writing the logline:

``` go
package main

// Imports omitted for brevity.

func main() {
  log.SetPrefix("\33[2K\r")

  tick := 0
  for {
    time.Sleep(1 * time.Second)
    if tick == 0 {
      fmt.Fprint(os.Stderr, "\33[2K\r", "[ ] Waiting")
    } else {
      fmt.Fprint(os.Stderr, "\33[2K\r", "[~] Waiting")
    }
    tick = 1 - tick
  }
}
```

At this point, you will notice every time the program outputs a log line, the throbber may not appear for a fraction of a second. It is because the throbber is being output once every second.

We need to make the throbber output more frequent.

## Complete Example

Below is a slightly more complete code example:

``` go
package main

import (
  "crypto/rand"
  "fmt"
  "log"
  "os"
  "time"
)

func main() {
  log.SetPrefix("\033[2K\r") // Clear current line before each log line.

  go func() {
    tick := -10
    for {
      time.Sleep(100 * time.Millisecond) // Output throbber more frequently, 10 times per second. But change state once every 10 ticks (1 second).
      if tick >= 0 {
        fmt.Fprint(os.Stderr, "\033[2K\r", "[ ] Waiting")
      } else {
        fmt.Fprint(os.Stderr, "\033[2K\r", "[~] Waiting")
      }
      tick++
      if tick == 10 {
        tick = 0 - tick
      }
    }
  }()

  for {
    time.Sleep(3 * time.Second)
    log.Println(rand.Int())
  }
}
```

The result looks like this:

{{< video src="example.webm" muted="true" >}}
