---
title: Setup Multiple Passphrases for a LUKS Device
date: 2023-09-13T00:10:00+06:00
tags:
  - LUKS
  - Linux
  - 100DaysToOffload
---

Let's say you have a computer at home shared by multiple people. And, you want to encrypt your hard drive with LUKS but not have to use the same passphrase.

You can do that. LUKS has 8 key slots (LUKS1 does, LUKS2 can support more).

When you set up a LUKS encrypted device you are configuring the first key slot only.

But by running the following command you can set up an additional passphrase:

``` sh {linenos=false}
cryptsetup luksAddKey <device>
```

## Example

``` sh
# Create an empty image file that we will turn into a LUKS device. You will probably be using a real device.
truncate --size=512MB dummy.img

# Create a LUKS device. You will be setting up your first passphrase here.
cryptsetup luksFormat dummy.img

# Add a second passphrase.
cryptsetup luksAddKey dummy.img

# Test both passphrases.
cryptsetup open --test-passphrase dummy.img
echo $? # If you enter the correct passphrase, the `cryptsetup open` command will exit with status 0.
```
