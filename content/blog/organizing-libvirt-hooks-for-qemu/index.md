---
title: 'Organizing Libvirt Hooks for Qemu'
date: 2023-07-16T16:30:00+06:00
tags:
  - Libvirt
  - Hooks
  - 100DaysToOffload
  - Qemu
---

Something I have been meaning to write about for a while is my KVM/VFIO-based gaming setup. Yes, I run Linux. And I run Windows, on Linux, in a virtual machine (VM). And it works! I game on it almost every day.

But the longer post has to wait for another day.

Today I am just sharing the short script that lets me keep my Libvirt hooks for Qemu a tad more organized:

``` bash
#!/bin/bash

set -e

SCRIPT="./${1}-${2}-${3}"
if [ -f "$SCRIPT" ]; then
    . "$SCRIPT"
fi
```

You can save the script as `/etc/libvirt/hooks/qemu`, and then you can put all your VM-specific scripts under `/etc/libvirt/hooks/` as separate files:

- `{vm}-prepare-begin`
- `{vm}-start-begin`
- `{vm}-started-begin`
- `{vm}-stopped-end`
- `{vm}-release-end`

For example, if you have a `win10-prepare-begin` script inside the `hooks/` directory, it will only be run when you start your `win10` VM at the prepare stage.

You can learn more about Libvirt hooks [here](https://libvirt.org/hooks.html).

But why not put these scripts inside a `qemu.d/` directory? Since v6.5.0, Libvirt lets you put multiple Qemu-specific hooks inside the `qemu.d/` directory. Libvirt still runs all the scripts for each of your VM like it runs `/etc/libvirt/hooks/qemu`.
