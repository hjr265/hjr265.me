---
title: 'Switch Monitor Input from Linux Command Line'
date: 2023-10-24T10:00:00+06:00
tags:
  - Linux
  - 100DaysToOffload
---

For the longest time, I have looked at computer monitors as these _dumb_ devices. All they do is turn video signals into colours on the screen.

I was out of touch with the progress.

Most modern computer monitors come with what is known as DDC/CI. It may be disabled by default, so you need to enable it using the monitor's on-screen display (OSD) settings.

Then, use a tool like `ddcutil` to switch the current input of the monitor.

For example, on the monitor (Asus PA328Q) I am currently using with my primary computer, I can run the following command to change my input source to DisplayPort:

``` txt {linenos=false}
# ddcutil setvcp 0x60 0x0f
```

Here `x60` implies "input source", the feature we are changing. And, `0x0f` implies "DisplayPort-1"

``` txt {linenos=false}
# ddcutil vcpinfo 0x60
VCP code 60: Input Source
   Selects active video source
   MCCS versions: 2.0, 2.1, 3.0, 2.2
   ddcutil feature subsets: 
   Attributes (v2.0): Read Write, Non-Continuous (simple)
   Attributes (v2.1): Read Write, Non-Continuous (simple)
   Attributes (v3.0): Read Write, Table (normal)
   Attributes (v2.2): Read Write, Non-Continuous (simple)
```

Again, I can run the following command to switch to Mini DisplayPort:

``` txt {linenos=false}
# ddcutil setvcp 0x60 0x10
```

Here `0x10` implies "DisplayPort-2", which is physically the Mini DisplayPort on the monitor.

To identify what values you can use with `0x60`, you can run the following command:

``` txt {linenos=false}
# ddcutil capabilities
```

Unfortunately, with the monitor I am using this command fails like so:

``` txt {linenos=false}
Unparsed capabilities string: (prot(monitor) type(LCD)model LCDPA328 cmds(01 02 03 07 0C F3) vcp(02 04 05 08 0B 0C 10 12 14(04 05 08) 16 18 1A 60(11 12 13 0F 10) 62 6C 6E 70 8D(01 02) A8 AC AE B6 C6 C8 C9 D6(01 04) DF) mccs_ver(2.1)asset_eep(32)mpu(01)mswhql(1))
Errors parsing capabilities string:
   Missing parenthesized value for segment model  at offset 30
Model: Not specified
MCCS version: Not specified
VCP Features:
```

However, I could still identify the possible values for feature `0x60` by looking at the unparsed capabilities string.

``` txt {linenos=false}
60(11 12 13 0F 10)
```

This part of the unparsed capabilities string means that the 3 HDMI ports and the 2 DisplayPort ports are identified by `0x11`, `0x12`, `0x13`, `0x0F` and `0x10`. With a bit of trial and error, I was able to map which value implies which physical port.
