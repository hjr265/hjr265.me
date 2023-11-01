---
title: "Uploading Files Over SSH in Go"
date: 2023-10-29T14:00:00+06:00
tags:
  - Go
  - SSH
  - 100DaysToOffload
  - Tidbit
---

If you access servers remotely over SSH connections, you are bound to have come across `scp`. It is what you use to upload files to these remote servers.

If you want to programmatically upload files like `scp` over an SSH connection to a remote server using Go, then you can use an [`ssh.Client`](https://pkg.go.dev/golang.org/x/crypto/ssh):

``` go
// Error handling omitted for brevity.

func uploadFile(sshClient *ssh.Client, filename string, mode os.FileMode, size int64, r io.Reader) {
  // Set up a new SSH session.
  sess, _ := sshClient.NewSession()
  defer sess.Close()

  // Write the file's metadata and contents to the stdin pipe.
  w, _ := sess.StdinPipe()
  go func() {
    defer w.Close()
    // Write "C{mode} {size} {filename}\n"
    fmt.Fprintf(w, "C%#o %d %s\n", mode, size, path.Base(filename))
    // Write the file's contents.
    io.Copy(w, r)
    // End with a null byte.
    fmt.Fprint(w, "\x00")
  }()

  sess.Stdout = os.Stdout
  sess.Stderr = os.Stderr
  // Run `scp -t {filename}` on the server-side.
  sess.Run(fmt.Sprintf("scp -t %s", filename))
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

// Upload the file.
f, _ := os.Open("hello.txt")
fi, _ := f.Stat()
uploadFile(c, fi.Name(), 0644, fi.Size(), f)
```
