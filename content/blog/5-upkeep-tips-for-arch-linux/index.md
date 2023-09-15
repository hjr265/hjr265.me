---
title: 5 Upkeep Tips for Arch Linux
date: 2023-09-13T00:16:00+06:00
tags:
  - ArchLinux
  - Linux
  - 100DaysToOffload
  - Upkeep
toc: yes
toch2only: yes
---

Yes, I am one of those guys: "I use Arch btw".

After distro-hopping for almost my entire time at the university, I found Arch Linux. I did momentarily switch to a MacBook after that, but I don't think I have tried any other Linux distribution for over a decade.

Arch Linux doesn't do everything the way I like on a Linux device. But it doesn't make me feel like switching to another distribution.

In this blog post, I will share 5 general tips for maintaining an Arch Linux installation.

## 1. Merge `.pacnew` files

Pacman. One of the first great things that come to mind when I hear Arch Linux. 

You probably are updating your system regularly. But did you notice that your Arch Linux installation has been accumulating `.pacnew` files? These files are created when you have made modifications to a file that Pacman manages.

That `/etc/mkinitcpio.conf` for example.

You have tuned it to your system, to your liking. You run `pacman -Syu` and Pacman sees a new version of this file. Instead of mucking around with the file you have modified, Pacman will store the new version as `/etc/mkinitcpio.conf.pacnew`.

It is now up to you to decide what you want to do about this file.

### Using `pacdiff` to Manage `.pacnew` Files

A really useful tool for managing `.pacnew` files is [`pacdiff`](https://wiki.archlinux.org/title/Pacman/Pacnew_and_Pacsave#pacdiff). This tool comes as a part of the `pacman-contrib` package. You can install it with:

``` txt {linenos=false}
# pacman -S pacman-contrib
```

Once you have installed the package, run the tool with `pacdiff`. You will probably want to use `sudo` as this tool will help you clean up places in the system that require root access.

By default, `pacdiff` will use `vimdiff` to display differences between the file you have and the `.pacnew` version. This will do the job.

But if you haven't made an oath to never touch a mouse or a trackpad, you can use a graphical diff program like `meld` like so: `DIFFPROG=meld pacdiff`.

To combine that with `sudo`, run:

``` txt {linenos=false}
$ sudo DIFFPROG=meld pacdiff
```

## 2. Cleaning Pacman Cache

Over years of running an Arch Linux installation, you will be accumulating a lot of cached packages that Pacman downloads. It is easy to [clean them up](https://wiki.archlinux.org/title/Pacman#Cleaning_the_package_cache).

You can run `pacman -Sc` to clean up cached packages that are not installed locally.

``` txt
~ » sudo pacman -Sc                                       
Packages to keep:
  All locally installed packages

Cache directory: /var/cache/pacman/pkg/
:: Do you want to remove all other packages from cache? [Y/n] 
removing old packages from cache...

Database directory: /var/lib/pacman/
:: Do you want to remove unused repositories? [Y/n] 
removing unused sync repositories...
```

You can run `pacman -Scc` to clean up all cached packages.

### Run `paccache` on Schedule

Even better is `paccache` which comes with the `pacman-contrib` package. By running `paccache -r` you can remove all cached packages and keep only the latest 3. To keep the latest 1 only, you can run `paccache -rk1`.

It comes with a Systemd timer that you can enable to run this cleanup weekly, automatically.

### Removing Partially Downloaded Cached Packages

Sometimes you end up with partially downloaded packages in the cache directory. `pacman -Sc` will warn you about this.

If you want to remove these files, the following command can take care of it:

``` txt {linenos=false}
# find /var/cache/pacman/pkg -name '*.part' -delete
```

Just make sure Pacman isn't running and these partially downloaded files are packages being downloaded right now.

## 3. Removing Unneeded/Unused Packages

There are two very useful `pacman -Q` commands for listing packages that you may no longer need to keep installed in your system.

### Orphan Packages

First, [`pacman -Qtd`](https://wiki.archlinux.org/title/Pacman/Tips_and_tricks#Removing_unused_packages_(orphans)). This command lists all the packages that were installed as a dependency for some other package at some point in the past, but no package depends on it anymore.

You will probably want to go through this list and uninstall any package that you think is no longer needed.

Alternatively, if you think you need any of these packages, you can mark them as explicitly installed by running `pacman -D --asdeps <package>`.

### Dropped Packages

Second, `pacman -Qm` will list all packages that are no longer a part of any of the configured remote repositories.

While going through this list you will want to ignore any package that you have manually installed or installed from AUR.

Whether you need these packages depends on you. Keep in mind that these packages are no longer been updated by Pacman.

## 4. Clean `~`, `~/.config`, `~/.cache` and `~/.local/share`

This one is a little harder.

When you uninstall an application using `pacman -R` it doesn't quite clean up things it created in your home directory.

Well-behaved applications will follow the [XDG Base Directory specification](https://wiki.archlinux.org/title/XDG_Base_Directory) and create directories under `~/.config`, `~/.cache` and `~/.local/share` that you can easily identify by the name. When you uninstall that application you can just check these directories and remove anything that you think is no longer needed.

Let's say you have uninstalled VLC. If you do not plan to reinstall VLC or care about any configuration files it may have left in your home directory, you can simply delete `~/.config/vlc`.

But then some applications deserve the dunce hat. These applications will leave their mark directly at the base of your home directory.

List all the hidden files in your home directory and remove anything that you know is no longer needed.

``` txt {linenos=false}
$ ls -lA ~
```

## 5. Review Systemd Services and Timers

Start by reviewing if you have any failing Systemd services:

``` txt {linenos=false}
$ systemctl --failed
```

You can check the relevant logs of the failing service using the `journalctl -u <service>` command. You can learn more about filtering `journalctl` output [here](https://wiki.archlinux.org/title/Systemd/Journal#Filtering_output).

You can also list all the services by running the following command:

``` txt {linenos=false}
$ systemctl list-units --type=service
```

In addition to services, review the Systemd timers that are on schedule:

``` txt {linenos=false}
$ systemctl list-units --type=timer
```

Disable services and timers that you do not need 

You can also identify the services that have been negatively impacting your system's boot time by running the following command:

``` txt {linenos=false}
$ systemd-analyze blame
```

This will show a list of services ordered by decreasing impact on startup time.

## Wrap Up

System maintenance is important to keep your operating system performing optimally. Doing it regularly makes it possible to keep using the same installation over several years.

For this blog post, I picked 5 of many possible upkeep tips.

If you ask me what else I like about Arch Linux other than Pacman. My answer would be the Wiki. It documents Arch Linux and its components so well that even users of other Linux distributions find it useful from time to time.

Never hesitate to look into the [ArchWiki](https://wiki.archlinux.org/) for additional [system maintenance](https://wiki.archlinux.org/title/system_maintenance) tips.
