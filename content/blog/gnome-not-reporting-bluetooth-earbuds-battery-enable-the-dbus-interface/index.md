---
title: 'Gnome Not Reporting Bluetooth Earbuds Battery? Enable the D-Bus Interface'
date: 2023-10-17T10:00:00+06:00
tags:
  - Linux
  - Bluetooth
  - 100DaysToOffload
---

I use one of these true wireless Edifier earbuds.

I noticed how Android reports the battery level of these earphones. But Gnome doesn't.

Turns out, you need to enable the experimental D-Bus interface in the Bluetooth daemon on Linux for Gnome to know the battery level of the connected wireless earbuds. 

On Arch Linux (which I use, btw), I had to modify `/etc/bluetooth/main.conf` and set:

``` txt {linenos=false}
Experimental = true
```

You will probably find a `# Experimental = false` line in that file that you can uncomment and change from `false` to `true`.

Then, restart the Bluetooth daemon:

``` sh {linenos=false}
sudo systemctl restart bluetooth.service
```

And now, Gnome reports the battery level of the earbuds under the Power tab in Gnome Settings.  

{{< image src="screen.png" alt="Screenshot of the Power tab in Gnome Settings" caption="Screenshot of the Power tab in Gnome Settings" >}}
