---
title: "Forwarding a Port Over SSH in Go"
date: 2023-11-16T12:50:00+06:00
tags:
  - Go
  - SSH
  - 100DaysToOffload
  - Tidbit
---

In the day and age where Kubernetes is the go-to tool for orchestrating your applications on the cloud, I have been spending time building [Bullet](https://github.com/FurqanSoftware/bullet). I enjoy working on this tool, building out features for it little by little. And it also allows me to learn so many details.

For example I just added the ability to forward ports from the remote server to the local over SSH. Bullet, being built using Go, I had to figure out how to forward a port over SSH in Go.

In the command line world, you would use the `-L` flag that the `ssh` command provides for something like this.

In Go, you need to do a bit more:

``` go
// Error handling omitted for brevity.

import (
  "io"
  "net"

  "golang.org/x/crypto/ssh"
)

func forwardPort(sshClient *ssh.Client, local, remote int) {
  l, _ := net.Listen("tcp", fmt.Sprintf(":%d", local))
  for {
    conn, _ := l.Accept()
    go connect(conn, remote)
  }
}

func connect(conn net.Conn, remote int) {
    sess, _ := c.Client.Dial("tcp", fmt.Sprintf("localhost:%d", remote))
    defer sess.Close()
    go io.Copy(conn, sess)
    io.Copy(sess, conn) // End connection when local TCP is closed.
}
```

You can use this function like so:

``` go
// Connect to the remote server over SSH.
c, _ := ssh.Dial("tcp", "carrot.local:22", &ssh.ClientConfig{
  User: "hjr265",
  Auth: []ssh.AuthMethod{
    ssh.PasswordCallback(func() (string, error) { return "KeyboardCat", nil }),
  },
  HostKeyCallback: ssh.InsecureIgnoreHostKey(),
})

// Forward port 80 on remote to local 8080.
forwardPort(c, 8080, 80)
```

And access it:

``` sh {linenos=false}
curl https://localhost:8080
```
