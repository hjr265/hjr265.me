---
title: 'Steam Link Black Screen on a Retro Pie Raspberry Pi 4'
date: 2023-10-27T08:30:00+06:00
tags:
  - RaspberryPi
  - RetroPie
  - 100DaysToOffload
---

Last night, I encountered a strange issue setting up Steam Link on a Raspberry Pi running Retro Pie.

Here is my attempt at a proper description of what I was seeing.

After installing Steam Link using Retro Pie's package manager and rebooting the Raspberry Pi, I could see "Ports" show up on the Emulation Station. Inside was "Steam Link".

I could start Steam Link just fine, have it test my network, and then pair it with Steam running on my gaming virtual machine.

However, after clicking on the "Start Playing" button, the TV showed a black screen with the mouse cursor visible on it. It also showed the tips bar at the bottom when starting up. I could also see the remote session menu by holding the escape button. And I could hear sounds made by Steam by moving the cursor around and clicking on things randomly.

What was the fix for this issue?

Two things:

- Increase the video RAM of the Raspberry Pi to 512 MB.
- Disable overscan.

You can configure both of these using the `raspi-config` program.

You can change Video RAM allocation from the "Performance Options" menu (shown as "GPU Memory"). And you can disable overscan from the "Display Options" menu (shown as "Underscan").

Alternatively, you can make these changes to the `/boot/config.txt`:

- Set `disable_overscan=1`
- Set `gpu_mem=512` under the `[all]` group

And after a reboot of the Raspberry Pi, you should see the "black screen" issue is gone.

I hope this blog post helps anyone facing the same issue with the same setup.
