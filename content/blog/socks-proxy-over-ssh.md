---
title: 'SOCKS Proxy Over SSH'
date: 2023-10-31T13:00:00+06:00
tags:
  - SSH
  - Networking
  - 100DaysToOffload
---

To test some of Toph's IP-based access control features, I needed to access it from a few different IP addresses than mine.

I thought I finally needed to get one of those VPN subscriptions YouTube content creators keep rambling about.

Fortunately, I remembered an easier way to do this.

You see, it is possible to run a SOCKS proxy that tunnels your connection over SSH. And it is built right into the `ssh` command:

``` txt {linenos=false}
ssh -D 9050 user@hostname
```

This command will create an SSH connection to the remote server identified by `hostname` and log in as `user`. At the same time, it will open port 9050 on localhost, which will be a SOCKS proxy.

I was able to access Toph from different IP addresses by using the servers that I already own.

You can test it out quickly with `curl`:

``` txt {linenos=false}
curl --proxy socks5://localhost:9050 https://ifconfig.me
```

It should output the IP address of the server you have connected to over SSH.
