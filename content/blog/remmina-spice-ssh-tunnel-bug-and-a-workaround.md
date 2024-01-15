---
title: Remmina SPICE SSH Tunnel Bug and a Workaround
date: 2024-01-15T19:45:00+06:00
tags:
  - Remmina
  - Bug
  - Workaround
---

I seem to come across the strangest of bugs.

## Bug

Remmina SPICE over an SSH tunnel fails to handle keyboard-mouse interactions.

I have set up several virtual machines with desktop operating systems using [Libvirt](https://libvirt.org/) on my primary computer. I can access these virtual machines over the network using Libvirt on my laptop. However, I wanted to use [Remmina](https://www.remmina.org/) as the SPICE client since it is more configurable.

All worked well until I tried to access a Manjaro Gnome virtual machine with the SSH tunnel feature on Remmina.

Strangely, the same setup works fine with a Windows virtual machine on the same host.

It also works fine if I open a port on the virtual machine host and access the Manjaro Gnome virtual machine without tunnelling over SSH. However, as soon as I tried to use Remmina's SSH tunnel feature, I could see the desktop in the virtual machine but not interact with it.

## Workaround

Fortunately, Remmina allows you to configure connection profiles to run commands before making a connection and after disconnecting it. It is possible to setup and teardown an SSH tunnel using these options.

I turned off Remmina's SSH tunnel feature for that problematic profile, and based on the [clue I found on Remmina's wiki](https://gitlab.com/Remmina/Remmina/-/wikis/Usage/Remmina-ssh-wizardry), I came up with the following commands.

Before:

``` sh {linenos=false}
/bin/sh -c 'ssh -o ControlMaster=auto -o ControlPath=$HOME/.ssh/remmina-tunnel-%r@%h:%p -4fnN -L 17890:localhost:17890 carrot.local'
```

After:

``` sh {linenos=false}
/bin/sh -c 'ssh -o ControlMaster=auto -o ControlPath=$HOME/.ssh/remmina-tunnel-%r@%h:%p -4fnN -O exit carrot.local'
```

In these commands, `17890` is the SPICE port number on the virtual machine and `carrot.local` is its hostname.

With this workaround, I no longer face the issue where Remmina SPICE over an SSH tunnel fails to handle keyboard-mouse interactions.

Although the workaround solves the problem I was facing, it is anticlimactic not knowing what causes this issue in the first place. Let me know please if you have any ideas.
