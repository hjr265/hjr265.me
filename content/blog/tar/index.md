---
title: "tar"
date: 2014-08-26T00:00:00+06:00
tags:
  - Tar
  - Tools
---

[![xkcd: tar](tar.png)](http://xkcd.com/1168/)

Tar is a file format that allows you to archive files and directories while preserving flags and other file information. Tar is also the name of the Unix utility that manipulates these files and is also popular for being notoriously enigmatic.

Did you know? Tar, the name, is derived from "tape archive".

*But Ridwan, I don't store my files on tapes. Do I need `tar`?*

Well, if you need to ask that question (and don't care about helping Rob disarm the bomb), then the answer may just be "no". But if you spend enough time in Unix-like systems, you are bound to come across a ".tar" file eventually.

Just like most other commands, `tar` takes a bunch of options, followed by some arguments.

You can **create** a tar archive by executing

``` sh {linenos=false}
tar cf archive.tar mario.txt luigi.txt peach.txt
```

The "c" tells `tar` to create an archive with the files "mario.txt", "luigi.txt" and "peach.txt"; and the "f" tells `tar` to emit it as the file "archive.tar". Note that if you use the traditional usage pattern (like above), then "f" must appear at the end of the options list and the archive filename must be the following argument.

To **extract** files from the archive, you execute 

``` sh {linenos=false}
tar xf archive.tar
```

Here, the "x" tells `tar` to extract the archive, and "f" simply points out the file "archive.tar" - the archive. You will be using "f" in almost all the commands involving `tar`, unless of course you are piping the archive to/from the `tar` command. Because that's how most of the archives are stored: as files.

But beware! Tar will replace all existing files on conflict while extracting. So, you might want to **list** the files before extracting them by executing

``` sh {linenos=false}
tar tf archive.tar
```

It is just like "xf", only with a "t" instead of an "x". Yes! I like stating the obvious. Here, the "t" simply tells `tar` to list the contents of the archive.

You can also **append** files to an existing archive by executing

``` sh {linenos=false}
tar rf archive.tar yoshi.txt
```

Using the "r" option (short for "append" - don't ask how) will allow you to append files to the end of the archive. Using it multiple times with the same file will append it multiple times.

To append files that have only changed since last added to the archive, or files that are new, execute

``` sh {linenos=false}
tar uf archive.tar koopa.txt
```

The "u" tells `tar` to append the file only if it is new or was not appended to the archive before. Using this multiple times on the same file may also lead to multiple entries of it.

In case you want to **delete** a file from the archive, you have to execute

``` sh {linenos=false}
tar --delete -f archive.tar luigi.txt
```

Oddly enough, `tar` do not provide a shorter version of the "--delete" option. Although, that makes the command above pretty self-explanatory. Just be sure to notice that small dash before "f".

Even though tar is a very old format, older than I am, it is still being used extensively. A number of Linux distributions' package repositories are built on top of this tar format. Software and their source codes are being distributed in ".tar" files and in its compressed variations. Clearly this is an example of one of those things that are simple yet very powerful.