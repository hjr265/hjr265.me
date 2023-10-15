---
title: "Updating Creality CR-10 Smart Firmware"
date: 2021-10-21T17:25:00+06:00
tags:
  - 3DPrinting
  - Creality
  - CR10Smart
---

There seems to be a lot of frustration with this 3D printer. And one of those frustration stems from how confusing the firmware update process is. The documentation is lacking in terms of some of the details.

To keep things simple, here are the important bits:

Updating the firmware of this 3D printer involves two separate steps. The first one is to update the firmware of the hardware. The second one is to update the firmware of the screen.

This after-sale training video does an okay job at demonstrating the overall process: https://www.youtube.com/watch?v=qZHdCAixygg

But some additional things that you need to keep in mind are:

- The SD card is 16 GB or smaller.

- It has one partition of type W95 FAT32 (LBA). You can use cfdisk to ensure this.

- The partiion is formatted as FAT32 (with an allocation size of 4096 KB, sector size of 512 bytes).
  
  You can do this with `mkdosfs`:

  ``` sh {linenos=false}
  sudo mkdosfs -s 8 /dev/sda1
  ```

_It is worth noting that it is possible to use an SD card of larger capacity. But all you have to do is make sure the partition is of a smaller size. In my attempt, I was able to use a 2 GB partition (type: W95 FAT32) on a 32 GB micro SD card with the above configuration to flash the screen firmware._

---

{{< box >}}

**Update (2023-10-10)**

Scott Wood (commented below) faced an issue where the update didn't work, and the printer didn't stay on. That was nine months ago as of writing this update. I didn't have a solution for the problem until I faced it myself today.

I tried to flash the firmware with a different SD card and solved the issue.

In my case, I tried with the SD card that came with the printer. And I followed the same steps to prepare it as outlined above.

Since the printer wouldn't stay on during the update, I just held the power button as instructed in the original training video.

{{< /box >}}
