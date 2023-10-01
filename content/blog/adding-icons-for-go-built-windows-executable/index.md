---
title: 'Adding Icons for Go-built Windows Executable'
date: 2023-09-30T19:55:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

I have been using Windows for video games only for several years now. But that changed a little as I started working on [Printd](https://github.com/FurqanSoftware/toph-printd), Toph's print server daemon.

{{< image src="splash.png" alt="Adding icons for Go-built Windows executable" >}}

An executable file (`.exe`) on Windows can provide its icons. If you build a Go program for Windows you get the generic executable icon, which is fine. But sometimes fine isn't enough.

Especially when adding an icon to a Go-built Windows executable is easy.

## Step 1: Prepare The Icon

You need one or more icon files (`.ico` ).

If you have a PNG file, there are online tools that you can use to convert from the PNG icon to a `.ico` file.

If you have the ImageMagick tool handy, you can use it like so:

``` sh
convert icon_16.png icon_32.png icon_48.png icon_256.png -colors 256 icon.ico
```

Notice how you can include several icons in a single `.ico` file. It allows you to embed icons of different sizes.

Try including 16×16, 32×32, 48×48 and 256×256 icons to cover all your bases.

## Step 2: Generate `.syso` Files

This step requires a special tool: [github.com/akavel/rsrc](https://github.com/akavel/rsrc).

Install `rsrc`:

``` sh
go install github.com/akavel/rsrc@latest
```

And use it to generate `.syso` files

``` sh
$GOPATH/bin/rsrc -arch 386 -icon icon.ico
$GOPATH/bin/rsrc -arch amd64 -icon icon.ico
```

Run this tool multiple times, once for each architecture you are targetting.

It will create a new `.syso` file every time. Store these files inside your `main` package.

## Step 3: Build the Go Program

Yes. In step 3, you only need to build the Go program.

Go will automatically pick any relevant `.syso` file in the `main` package directory and include that in the executable built for Windows.
