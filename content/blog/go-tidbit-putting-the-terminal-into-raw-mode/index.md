---
title: 'Go Tidbit: Putting The Terminal Into Raw Mode'
date: 2023-09-20T00:15:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
toc: yes
---

I learned something new today. It helped solve a long-standing bug in [Bullet](https://github.com/FurqanSoftware/bullet).

Bullet is an application deployment tool that I wrote several years ago. It is a simple tool that SSHs into a server and uses Docker to run applications.

I use it in production for some of my projects.

## About the Bug

Bullet can SSH into a remote server, spin up a one-off container and attach the terminal to it. You could run commands in that container as if you were SSH'ed directly into that environment.

Except control characters didn't work.

You couldn't press `Ctrl+C` to interrupt the currently running program. You couldn't press `Ctrl+D` to signal EOF.

In fact, for example, pressing `Ctrl+C` would kill Bullet, ending the SSH session.

And why didn't control characters work? Because I didn't put the terminal into raw mode before connecting to the container over SSH.

## SSH Client in Go

Here is a simplified Go code that connects to a remote server over SSH and acts like an SSH client:

``` go
package main

import (
  "os"

  "github.com/mattn/go-tty"
  "golang.org/x/crypto/ssh"
)

func main() {
  // Connect to the remote server over SSH.
  c, _ := ssh.Dial("tcp", "carrot.local:22", &ssh.ClientConfig{
    User: "hjr265",
    Auth: []ssh.AuthMethod{
      ssh.PasswordCallback(func() (string, error) { return "KeyboardCat", nil }),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  })

  // Set up a new SSH session.
  sess, _ := c.NewSession()
  defer sess.Close()

  // Set up PTY with the SSH session on the remote host.
  tty, _ := tty.Open()
  defer tty.Close()
  w, h, _ := tty.Size()
  sess.RequestPty("xterm", h, w, ssh.TerminalModes{
    ssh.ECHO:          1,
    ssh.TTY_OP_ISPEED: 14400,
    ssh.TTY_OP_OSPEED: 14400,
  })

  // Use current stdin, stdout and stderr with the SSH session.
  sess.Stdin = os.Stdin
  sess.Stdout = os.Stdout
  sess.Stderr = os.Stderr
  sess.Run("/bin/bash")
}
```

It is almost like running `ssh carrot.local`.

In fact, after running the Go code, I was greeted with a familiar prompt.

``` txt {linenos=false}
[hjr265@Potato gossh]$ go run .
[hjr265@Carrot ~]$ 
```

But, unlike the `ssh` program, pressing `Ctrl+C` here terminated the Go program instead of sending a `SIGINT` to the remote server.

``` txt {linenos=false}
[hjr265@Potato gossh]$ go run .
[hjr265@Carrot ~]$ signal: interrupt
[hjr265@Potato gossh]$ 
```

## The Fix: Put The Terminal Into Raw Mode

The solution turned out to be fairly simple. But it is something that I only came to learn about today.

Terminals can run in raw mode or cooked mode.

>  In cooked mode data is preprocessed before being given to a program, while raw mode passes the data as-is to the program without interpreting any of the special characters. [\[...\]](https://en.wikipedia.org/wiki/Terminal_mode)

What cooked mode means is dependent on the operating system. However, in cooked mode, control characters are handled by the operating system.

When `Ctrl+C` is pressed, the operating system sends `SIGINT` to the currently running program. In most cases, the program aborts immediately.

But what we want here is to handle that control character ourselves.

And we can do that by putting the terminal into raw mode:

``` go
oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
if err != nil {
  panic(err)
}
defer term.Restore(int(os.Stdin.Fd()), oldState)
```

The function `term.MakeRaw` comes with the [`golang.org/x/term`](https://pkg.go.dev/golang.org/x/term) package.

A call to this function will put the terminal into raw mode and return a copy of the old state. This old state can be used to restore the terminal to its previous mode with the `term.Restore` function.

## Working SSH Client in Go

``` go
package main

import (
  "os"

  "github.com/mattn/go-tty"
  "golang.org/x/crypto/ssh"
  "golang.org/x/term"
)

func main() {
  // Connect to the remote server over SSH.
  c, _ := ssh.Dial("tcp", "carrot.local:22", &ssh.ClientConfig{
    User: "hjr265",
    Auth: []ssh.AuthMethod{
      ssh.PasswordCallback(func() (string, error) { return "KeyboardCat", nil }),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  })

  // Set up a new SSH session.
  sess, _ := c.NewSession()
  defer sess.Close()

  // Put the terminal into raw mode.
  oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
  defer term.Restore(int(os.Stdin.Fd()), oldState)

  // Set up PTY with the SSH session on the remote host.
  tty, _ := tty.Open()
  defer tty.Close()
  w, h, _ := tty.Size()
  sess.RequestPty("xterm", h, w, ssh.TerminalModes{
    ssh.ECHO:          1,
    ssh.TTY_OP_ISPEED: 14400,
    ssh.TTY_OP_OSPEED: 14400,
  })

  // Use current stdin, stdout and stderr with the SSH session.
  sess.Stdin = os.Stdin
  sess.Stdout = os.Stdout
  sess.Stderr = os.Stderr
  sess.Run("/bin/bash")
}
```

Now I can send `SIGINT`s to my heart's content:

``` txt {linenos=false}
[hjr265@Potato gossh]$ go run .
[hjr265@Carrot ~]$ sleep 5
^C
[hjr265@Carrot ~]$
```
