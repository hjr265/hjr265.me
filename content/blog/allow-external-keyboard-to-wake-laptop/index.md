---
title: Allow External Keyboard to Wake Laptop
date: 2023-10-04T13:15:00+06:00
tags:
  - Linux
  - Udev
  - 100DaysToOffload
aliases:
  - /blog/allow-external-keyboard-to-wake-up-laptop/
---

My setup at my workspace isn't much: a display and a few peripherals connected to a USB-C hub. This way, I can come in, plug in the hub to my laptop, and start working.

But if my laptop goes to sleep, I can wake it up only by pressing a key on the internal keyboard or the trackpad. Since I keep the lid of my laptop closed, I have to take it out of the stand first, then open the lid (which causes the laptop to wake up anyway).

However, it is easy to let a USB device wake the laptop. The configuration is per USB port.

``` sh {linenos=false}
echo enabled > /sys/bus/usb/devices/usb3/power/wakeup
```

The easy route is to allow all the ports to wake the laptop up. But that doesn't sit right with me.

Instead, I added a small _udev_ rule to allow just one port to wake the laptop - the port where I connect the USB C hub. Which, in turn, my keyboard connects to.

``` txt
ACTION=="add", SUBSYSTEM=="usb", ATTRS{idVendor}=="2516", ATTRS{idProduct}=="0067" RUN+="/bin/sh -c 'echo enabled > /sys/bus/usb/devices/usb3/power/wakeup'"
ACTION=="remove", SUBSYSTEM=="usb", ATTRS{idVendor}=="2516", ATTRS{idProduct}=="0067" RUN+="/bin/sh -c 'echo disabled > /sys/bus/usb/devices/usb3/power/wakeup'"
```

The ArchWiki has a well-written [outline around creating and managing _udev_ rules](https://wiki.archlinux.org/title/udev#About_udev_rules).

The vendor ID, the product ID, and the USB bus number can be known using the `lsusb` command:

``` txt
             Vendor ID ↓
Bus 003 Device 022: ID 2516:0067 Cooler Master Co., Ltd. MK750
      ↑ USB Bus             ↑ Product ID
```

This way, when I connect my keyboard, _udev_ will enable wake-up on USB activity for that specific port. And it will be disabled when I disconnect the keyboard.

The only caveat is that I need to connect or disconnect the keyboard while the laptop is awake.