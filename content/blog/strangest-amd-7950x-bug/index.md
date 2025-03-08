---
title: Strangest AMD Ryzen 7950x Bug
date: 2025-03-08T19:20:00+06:00
tags:
  - AMD
  - Bug
  - Workaround
---

I upgraded my primary computer to the AMD AM5 platform sometime in mid-2023. Before the upgrade, I had been using an Intel Core i7 4790 with 32 GB of RAM for about a decade. For the upgrade, I went with an AMD Ryzen 7950x with 128 GB of RAM. Given that my work lately has started to involve a lot of virtual machines, this upgrade was worth every bit.

However, the upgrade came with a few annoyances. One of them, and the worst one, in my opinion, was random reboots.

The reboots were abrupt and did not involve kernel panics or blue/black screens. Randomly, the computer rebooted like someone had pressed the "reset" button on its front panel.

## Months of Tinkering

At first, I thought it had something to do with memory stability. I tried running MemTest86 overnight several times, and it never crashed. I also tried running other stress tests involving CPU and memory. It never crashed.

I tried disabling the XMP profile, tuning the memory frequency, and other measures, but I could never stop the random reboots.

I also tried different Linux kernel versions and flavours, whether the Zen variant or the LTS one. Nothing made any difference. The reboots were random and always happened at the worst moments.

3D rendering jobs that ran for several days worked fine. But then, randomly, a 3D rendering job that was running for a few hours and was 97% done was interrupted by a reboot.

After several months of suffering, I noticed a pattern to this.

The computer seemed to reboot randomly after doing certain things, but it stayed stable until that point.

For example, a warm reboot (including after a random reboot) would make the computer more likely to undergo another random reboot. The computer would almost never reboot randomly after a cold boot.

But then, if I started a virtual machine with CPU host-passthrough, it left the computer in a ready-to-have-a-random-reboot state. Even if I stopped the virtual machine, a random reboot was guaranteed shortly.

I also tried running different BIOS versions. Nothing seemed to help.

I was already running the 7950x in 105W Eco mode, but running it in different modes didn't help.

The last thing I tried was disabling Global C-State Control. That also didn't help either.

## CPU Frequency

The clues I had started to note down didn't answer the mystery, but it all seemed to point to the CPU.

I recalled reading a Reddit comment that suggested, with no evidence, that AMD AM5 was probably not boosting their CPUs right. It mentioned something about AMD's Precision Boost Overdrive (PBO) pushing the CPU beyond its limit.

What if I capped the maximum possible CPU frequency?

I first applied a maximum CPU frequency limit of 5.45 GHz using `cpufreq`. And what do you know? My computer became entirely stable.

I kept using my computer with the limit in place, and then I did everything that would result in a random reboot later on. And nothing seemed to result in a crash.

After ensuring stability for a couple of weeks, I applied the same limit from the BIOS, and it worked like a charm.

## Wrap Up

Well, modern hardware is not built like it used to be. Finding modern hardware that isn't buggy today is as challenging as finding good software through random ads on the Internet.

Somehow, quality control is no longer a thing.
