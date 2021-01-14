---
title: "Using an external USB drive for the root file system of a Raspberry Pi"
date: 2014-01-14T00:00:00+06:00
---

By design, a Raspberry Pi always requires an SD card to boot from. But one can still have its root partition located on an external storage device. Be it for reasons involving speed improvement, or avoid challenging the [write endurance](http://en.wikipedia.org/wiki/Flash_memory#Write_endurance) of an SD card.

The details in the following steps may vary based on the distribution of Linux being used, but the fundamental idea should be similar anyway:

- Assuming a distribution of Linux is already installed on the SD card, use it to boot a Raspberry Pi up.

- Once booted up, execute the following command and ensure that the desired new root partition shows up in the list of available partitions. In this case, we are looking for `sda1`. Here `mmcblk0p5` is the current root partition denoted by the "/" at the end.

    ``` sh {linenos=false}
    lsblk
    ```
    ``` txt {linenos=false}
    NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    sda           8:0    0  1.8T  0 disk 
    |-sda1        8:1    0 23.3G  0 part 
    |-sda2        8:2    0    2G  0 part 
    `-sda3        8:3    0  1.8T  0 part 
    mmcblk0     179:0    0  3.8G  0 disk
    |-mmcblk0p1 179:1    0   90M  0 part /boot
    |-mmcblk0p2 179:2    0    1K  0 part
    `-mmcblk0p5 179:5    0  1.7G  0 part /
    ```

- Overwrite the new root partition with the contents of the current root partition. **Beware** that this will replace all existing contents (if any) on the new root partition!

    ``` sh {linenos=false}
    dd if=/dev/mmcblk0p5 of=/dev/sda1 conv=sync,noerror
    ```

- Once the execution of the command completes, check the new root partition file system for errors. The following command assumes an `ext4` partition. You may want to use the appropriate tool in this step, and the next one as well. Press "y" to fix errors (if found).

    ``` sh {linenos=false}
    e2fsck -f /dev/sda1
    ```

- Execute the following command to update the file system to use the entire available space of the partition.

    ``` sh {linenos=false}
    resize2fs /dev/sda1
    ```

- Change the value of the root parameter within the "/boot/cmdline.txt" file to "/dev/sda1", while keeping a backup of the original file.

    ``` sh {linenos=false}
    sed -i.bak 's/root\=[^ ]*/root\=\/dev\/sda1/' /boot/cmdline.txt
    ```

- Some Linux distributions require you to list the root partition in the `fstab` file. If the one you are using does, then you can fix it by editing the "/etc/fstab" file within the new root partition.

    ``` sh {linenos=false}
    mount /dev/sda1 /mnt
    nano /mnt/etc/fstab
    ```

    Change "/dev/mmcblk0p5" (the trailing number may vary depending on where your current root partition is located) to "/dev/sda1"

- Reboot the Raspberry Pi.

    ``` sh {linenos=false}
    reboot
    ```

At this point, the Raspberry Pi should boot from the SD card, but use the external USB drive for the root file system.