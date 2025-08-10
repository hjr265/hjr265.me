---
title: "Virtual Private Arcade: KVM, IOMMU, VFIO, GPU Passthrough, and Other Cool Stuff"
date: 2021-04-12T00:00:00+06:00
draft: true
tags:
  - Linux
  - VFIO
  - KVM
---

I have a desktop computer that I affectionately named XT. I named it after the first computer I ever used: an IBM XT.

Currently XT is powered by an Intel i7 4790 CPU, 32 GB of RAM, and an Nvidia RTX 2080 Ti GPU. For the longest time I have been using my computer with 3 operating systems: Arch Linux, Windows, and macOS (Hackintosh). But with dual/triple boot solution, you start to hate the dreaded reboots. Since I do most of my work on Linux, if I wanted to take a quick break and play something that runs on Windows only, I would have to reboot to Windows. 

So, here is what I did: I kept Arch, and I installed Windows in a Virtual Machine with the GPU passed through. Now I can boot and switch to Windows with a single command, and not have to go through a reboot:

```
sudo ./start.sh
```

And, if I want to play remotely (using [Steam Link](https://store.steampowered.com/steamlink/about/) or [Moonlight](https://moonlight-stream.org/)) from my phone or my TV, I can always keep running Windows with Steam or Nvidia Game Stream in the background.

## Script with Qemu Command

Here is a breakdown of the start.sh script.

``` sh
#!/bin/bash

set -e

if [[ $EUID -ne 0 ]]; then
   echo "Please run as root" 
   exit 1
fi
```

``` sh
echo Preparing host

echo .. KVM
echo 1 > /sys/module/kvm/parameters/ignore_msrs
```

``` sh
echo .. CPU
cpupower frequency-set -g performance > /dev/null
```

``` sh
echo .. GPU
modprobe vfio-pci >/dev/null
[ -f /sys/bus/pci/devices/0000:01:00.0/driver/unbind ] && echo '0000:01:00.0' > /sys/bus/pci/devices/0000:01:00.0/driver/unbind
echo '10de 1e07' > /sys/bus/pci/drivers/vfio-pci/new_id
[ -f /sys/bus/pci/devices/0000:01:00.1/driver/unbind ] && echo '0000:01:00.1' > /sys/bus/pci/devices/0000:01:00.1/driver/unbind
echo '10de 10f7' > /sys/bus/pci/drivers/vfio-pci/new_id
[ -f /sys/bus/pci/devices/0000:01:00.2/driver/unbind ] && echo '0000:01:00.2' > /sys/bus/pci/devices/0000:01:00.2/driver/unbind
echo '10de 1ad6' > /sys/bus/pci/drivers/vfio-pci/new_id
[ -f /sys/bus/pci/devices/0000:01:00.3/driver/unbind ] && echo '0000:01:00.3' > /sys/bus/pci/devices/0000:01:00.3/driver/unbind
echo '10de 1ad7' > /sys/bus/pci/drivers/vfio-pci/new_id
```

``` sh
echo .. Hugepages
echo 3 > /proc/sys/vm/drop_caches
echo 1 > /proc/sys/vm/compact_memory
echo 10 > /sys/kernel/mm/hugepages/hugepages-1048576kB/nr_hugepages
sysctl vm.hugetlb_shm_group=992 > /dev/null
mkdir /dev/hugepages1G > /dev/null || true
mount -t hugetlbfs -o pagesize=1G none /dev/hugepages1G > /dev/null || true
```

``` sh
echo Starting QEMU
qemu-system-x86_64 \
    -machine pc-q35-5.2,accel=kvm,usb=off,vmport=off,dump-guest-core=off,kernel_irqchip=on \
    -cpu host,migratable=on,hv-time,hv-relaxed,hv-vapic,hv-vpindex,hv-synic,hv-stimer,hv-vendor-id=26051990,kvm=off \
    -smp 6,sockets=1,dies=1,cores=3,threads=2 \
    -global kvm-pit.lost_tick_policy=delay \
    -global ICH9-LPC.disable_s3=1 \
    -global ICH9-LPC.disable_s4=1 \
    -boot strict=on \
    -name arcade \
    -uuid 3e96588a-f0a2-437a-b994-9f513b486ca5 \
    -m 9216 \
    -mem-path /dev/hugepages1G \
    -audiodev pa,id=hostaudio0,server=unix:/tmp/pulse-socket \
    -device pcie-root-port,id=pcieroot.1,port=1,chassis=1,bus=pcie.0,addr=1c.0,multifunction=on \
    -device pcie-root-port,id=pcieroot.2,port=2,chassis=2,bus=pcie.0,addr=02.0 \
    -device pcie-root-port,id=pcieroot.3,port=3,chassis=3,bus=pcie.0,addr=03.0 \
    -device pcie-root-port,id=pcieroot.4,port=4,chassis=4,bus=pcie.0,addr=04.0 \
    -device pcie-root-port,id=pcieroot.5,port=5,chassis=5,bus=pcie.0,addr=05.0 \
    -device pcie-root-port,id=pcieroot.6,port=6,chassis=6,bus=pcie.0,addr=06.0 \
    -device vfio-pci,host=0000:01:00.0,bus=pcieroot.1,addr=00.0,multifunction=on,x-vga=on \
    -device vfio-pci,host=0000:01:00.1,bus=pcieroot.1,addr=00.1 \
    -device vfio-pci,host=0000:01:00.2,bus=pcieroot.1,addr=00.2 \
    -device vfio-pci,host=0000:01:00.3,bus=pcieroot.1,addr=00.3 \
    -device virtio-blk-pci,drive=drive0,bus=pcieroot.2,addr=0x0,bootindex=1,iothread=iothread1,write-cache=on \
    -device virtio-blk-pci,drive=drive1,bus=pcieroot.3,addr=0x0,iothread=iothread1,write-cache=on \
    -device ide-cd,drive=cdrom0,bus=ide.1 \
    -device virtio-net,netdev=netdev0,mac=de:ad:be:ef:aa:b4 \
    -device virtio-mouse-pci,bus=pcieroot.4,addr=0x0 \
    -device virtio-keyboard-pci,bus=pcieroot.5,addr=0x0 \
    -device ich9-intel-hda,bus=pcieroot.6,addr=0x1b \
    -device hda-micro,audiodev=hostaudio0 \
    -drive id=pflash0,file=/usr/share/edk2-ovmf/x64/OVMF_CODE.fd,if=pflash,readonly,format=raw \
    -drive id=pflash1,file=./nvram_vars.fd,if=pflash,format=raw \
    -drive id=drive0,file=./drive0.qcow2,if=none,media=disk,cache=none,format=qcow2 \
    -drive id=drive1,file=/dev/sdb1,if=none,media=disk,cache=none,aio=native,format=raw \
    -drive id=cdrom0,if=none,media=cdrom \
    -display none \
    -vga none \
    -no-hpet \
    -netdev tap,id=netdev0,ifname=tap0,script=no,downscript=no \
    -serial none \
    -parallel none \
    -monitor tcp:127.0.0.1:52431,server,nowait \
    -overcommit mem-lock=on \
    -enable-kvm \
    -daemonize \
    -pidfile ./qemu.pid \
    -rtc base=localtime,driftfix=slew \
    -nodefaults \
    -sandbox on,obsolete=deny,elevateprivileges=deny,resourcecontrol=deny \
    -no-user-config \
    -msg timestamp=on \
    -object iothread,id=iothread1 \
    -object input-linux,id=mouse1,evdev=/dev/input/by-id/usb-Corsair_CORSAIR_IRONCLAW_RGB_WIRELESS_Gaming_Dongle_84FA3C9F45063637-event-mouse \
    -object input-linux,id=kbd1,evdev=/dev/input/by-id/usb-CORSAIR_CORSAIR_K63_Wireless_USB_Receiver_1103903EAF4C20E15CCD8AF9F5001C01-event-kbd,grab_all=on,repeat=on \
```

``` sh
echo Switching display
ddcutil setvcp 0x60 0x10
```