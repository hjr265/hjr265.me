---
title: 'Android Emulator Slow As a Snail; Reason BTRFS'
date: 2023-11-23T10:00:00+06:00
tags:
  - Android
  - AndroidEmulator
  - 100DaysToOffload
---

I have been writing software professionally for over a decade. I have been writing software for even longer than that.

This week was the first time I wrote an Android program. Some journey it was. But that is a story for another day.

Today, in this blog post, I want to share a strange issue I encountered with Android Emulator and a fix.

If you are using BTRFS, you probably have copy-on-write (COW) enabled for the files on it. If your Android Emulator boot cache is on a BTRFS partition, you will notice the emulator running as slow as a snail.

The fix?

``` sh {linenos=false}
echo "QuickbootFileBacked = off" >> ~/.android/advancedFeatures.ini
```

Recent versions of Android Emulator will use a pre-allocated file as the backing storage of the guest RAM. It allows the emulator to store Quickbook snapshots during runtime.

But this does not play well with a copy-on-write filesystem.

By running the command from above, you are telling Android Emulator to not use the feature.
