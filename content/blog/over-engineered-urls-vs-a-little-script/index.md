---
title: Over-engineered URLs vs. A Little Script
date: 2023-10-04T08:15:00+06:00
tags:
  - Bash
  - Curl
  - 100DaysToOffload
---

One of the local courier services where I live will send you this helpful notification SMS with a tracking URL whenever they pick up a parcel for you. It is useful. At least compared to the hundreds of spam text messages we get here daily.

{{< image src="sms.png" alt="Screenshot of SMS from the local courier service" captionMD="It reads: [{Courier}] We have collected your parcel {ID} from {Vendor}. Track {MaskedURL}." >}}

The tracking URL is masked, much like a shortened URL, but not exactly short. The URL will redirect you to the real tracking URL only if you are not accessing it from a mobile device.

I will try not to make this sound like a rant. But these over-engineered URLs make me 🤮.

What does it do when you access it from a mobile device? Send you to the Google Play Store or App Store, of course.

Am I wrong to find it infuriating?

I took matters into my own hands and wrote a Bash script to figure out the real tracking URL.

``` bash
#!/bin/bash

curl -s -I -o /dev/null -w '%header{location}' $1
echo
```

Curl command breakdown:

- `-s`: Silent
- `-I`: HEAD-only request
- `-o /dev/null`: Discard output
- `-w %header{location}`: Output the value of the `location` header to stdout
- `$1`: The first argument passed to this script

Since `curl` here will not print a new line at the end, the `echo` pushes the next prompt to a new line.

With a terminal emulator like [Termux](https://termux.dev/en/) on my phone, I can quickly run a masked tracking URL through this script to figure out the real tracking URL.

``` txt {linenos=false}
$ ./unmaskurl.sh {MaskedURL}
{TrackingURL}
```
