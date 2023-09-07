---
title: Making a Gnome Background Slide Show
date: 2023-09-07T11:50:00+06:00
tags:
  - Gnome
  - OpenSource
  - 100DaysToOffload
---

"Background Slide Show". That is what [Gnome](https://www.gnome.org/) uses in code to refer to its time-based dynamic backgrounds feature.

Instead of having a static background, a Gnome Background Slide Show allows you to have a set of images that Gnome selects from based on the current time of the day. You can also configure the dynamic background to use transitions when changing from one picture to the next.

There are tools that you can install to build dynamic backgrounds. But I wanted to read some documentation and figure out how it works.

I recalled seeing something about writing an XML file to describe the slide show. 

## Making a Gnome Background Slide Show

<!-- > If you don't want the nitty-gritty details and want to make a Gnome Background Slide Show, use this [web-based dynamic background maker](https://mkdbg.hjr265.me/). -->

After jumping through some hoops, I was able to find [a C file under the gnome-desktop](https://gitlab.gnome.org/GNOME/gnome-desktop/-/blob/89d70faa26612d35808b060a437ea06d325cbc6d/libgnome-desktop/gnome-bg/gnome-bg-slide-show.c) project that dealt with dynamic backgrounds.

A simple two-image dynamic background can be defined as follows:

``` xml
<background>
  <starttime>
    <year>2023</year>
    <month>9</month>
    <day>7</day>
    <hour>6</hour>
    <minute>0</minute>
    <second>0</second>
  </starttime>

  <static>
    <file>/home/hjr265/Pictures/Backgrounds/patterns-day.jpg</file>
    <duration>36000.0</duration>
  </static>

  <transition type="overlay">
    <duration>7200.0</duration>
    <from>/home/hjr265/Pictures/Backgrounds/patterns-day.jpg</from>
    <to>/home/hjr265/Pictures/Backgrounds/patterns-night.jpg</to>
  </transition>

  <static>
    <file>/home/hjr265/Pictures/Backgrounds/patterns-night.jpg</file>
    <duration>36000.0</duration>
  </static>

  <transition type="overlay">
    <duration>7200.0</duration>
    <from>/home/hjr265/Pictures/Backgrounds/patterns-night.jpg</from>
    <to>/home/hjr265/Pictures/Backgrounds/patterns-day.jpg</to>
  </transition>
</background>
```

This XML sets the background to `day.jpg` for 10 hours and then gradually transitions from `day.jpg` to `night.jpg` over 2 hours. The `night.jpg` stays for 10 hours and then gradually transitions back to `day.jpg` over 2 hours. After that, the slide show starts all over again.

You can, of course, add more images to the set.

## Using a Gnome Background Slide Show

To use a Gnome Background Slide Show, you have two options:

Either, use the [Gnome Tweaks](https://wiki.gnome.org/action/show/Apps/Tweaks?action=show&redirect=Apps%2FGnomeTweakTool) and select the XML as your background.

Or, add an entry to the `~/.local/share/gnome-background-properties` directory so that you can pick this dynamic background from Gnome Settings.

``` xml
<wallpapers>
  <wallpaper deleted="false">
    <name>Patterns</name>
    <filename>/home/hjr265/Pictures/Backgrounds/patterns-dynamic.xml</filename>
    <options>zoom</options>
  </wallpaper>
</wallpapers>
```

The `filename` element should contain the path to the XML defining the dynamic background.

## Side Note

The open-source community is thriving. Undoubtedly, that is all thanks to the people contributing to this community.

Perhaps emphasis on writing more documentation could be the next logical step.
