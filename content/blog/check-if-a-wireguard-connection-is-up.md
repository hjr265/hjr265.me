---
title: "Check If a WireGuard Connection Is Up"
date: 2023-11-21T15:55:00+06:00
tags:
  - WireGuard
  - Linux
  - 100DaysToOffload
---

I have several scripts and automation on my primary computer at home that can run when connected to the local area network of my workspace through a WireGuard connection.

These scripts are for routine tasks for my servers at my workspace, like backing them up to remote storage.

When the WireGuard connection is not running, the scripts fail at different points.

I wanted the scripts to fail right when they start if the WireGuard connection is not running.

For that, I prefixed the scripts with a few lines of code.

Running `wg show` shows the status of all WireGuard connections. But this command requires root. It is not an ideal option since none of the scripts requires root access on the local computer for any other reason.

Instead, I added the following code:

``` sh
WG_INTERFACE_NAME=wgCarrot

ip l | grep "$WG_INTERFACE_NAME" > /dev/null
WG_INTERFACE_CONNECTED=$?

if [ $WG_INTERFACE_CONNECTED -ne 0 ]; then
  echo $WG_INTERFACE_NAME is not connected.
  exit 1
fi
```

I am using `ip l` and `grep` to determine if the WireGuard interface exists. If it does, I can assume that the WireGuard connection is active.

If not, the script exits with an error message.
