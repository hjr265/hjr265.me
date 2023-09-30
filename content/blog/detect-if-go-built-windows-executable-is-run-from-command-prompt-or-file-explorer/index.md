---
title: 'Detect If Go Built Windows Executable Is Run From Command Prompt or File Explorer'
date: 2023-09-30T11:00:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

In Windows, unlike the Unix-like POSIX-compatible operating systems, there is this notion of an application subsystem: `console` vs. `windows`.

If you build a Go program for Windows, it will, by default, use the `console` subsystem. When you start this program from File Explorer (e.g. by double-clicking its icon), Windows will show a console (like the Command Prompt window) and run the program inside the console.

When running a `console` subsystem program that finishes quickly, you may notice the console window appears and disappears quickly. You may not even see it flash.

Take this Go program as an example:

``` go
package main

import (
  "fmt"
  "time"
)

func main() {
  fmt.Println(time.Now())
}
```

It prints the current time.

If you run it from the Command Prompt (`cmd.exe`) you will see the time printed. But if you double-click on the executable from File Explorer, the console window will appear and disappear so quickly that you may not be able to see it at all.

One way to prevent the console window from disappearing so quickly is to make the program wait for a key press before exiting. But then this would require a key press to exit the program even when run from the Command Prompt.

The good news is that, on Windows, you can tell if a program is run from Command Prompt or not.

One way to do this is obvious: check if `cmd.exe` is the parent process.

## `PROMPT` Environment Variable

But the other approach is easier to implement: check if the  environment variable `PROMPT` is set to a non-empty value.

``` go
package main

import (
  "fmt"
  "os"
  "time"
)

func main() {
  if os.Getenv("PROMPT") != "" {
    defer stay()
  }

  fmt.Println(time.Now())
}

func stay() {
  fmt.Println(os.Stderr, "Press enter to exit.")
  fmt.Scanln()
}
```

This environment variable is set by `cmd.exe`. A program can check its presence to tell if it was run from the Command Prompt or by other means (e.g. double-clicking on the executable in the File Explorer).

## Real World Use

I am using this technique in [Printd](https://github.com/FurqanSoftware/toph-printd) Toph's print server daemon, to determine if the program should wait for a keypress before exiting.

This way, if Printd fails to start (e.g. due to errors in the configuration file), then Printd can still print helpful error messages and not worry about the messages being missed by the user due to the console window disappearing quickly.

{{< image src="explorer.png" alt="Screenshot of Printd Running in Console" caption="Run From File Explorer" >}}

<br>

{{< image src="cmd.png" alt="Screenshot of Printd Running in Command Prompt" caption="Run From Command Prompt" >}}
