---
title: "Fixing Creality CR-10 Smart Touch Screen Orientation Issue"
date: 2023-10-13T10:00:00+06:00
tags:
  - 3DPrinting
  - CR10Smart
  - 100DaysToOffload
  - Creality
---

The Creality CR-10 Smart is a very _sensitive_ 3D printer, especially when updating firmware.

Are you using an SD card that is too small? The printer won't update. Too large? Again, the printer won't update.

You formatted it to FAT32 but used an allocation size that isn't exactly 4096 KB. Tough luck: the printer won't update.

You successfully updated the firmware to `CR-10 Smart Marlin2.0.6SWV1.0.14HWCRC2405V1.2.zip`? Well, now your touch screen is going to behave weirdly.

When I updated my printer to that firmware (i.e. v1.0.14), it mirrored the touch screen across both axes.

The home tab, which appears in the top-left corner, is still in the top-left corner but is activated when you touch the screen in the bottom-right corner.

The settings tab, still in its bottom-left corner, is activated when you touch the screen in the top-right corner.

Not fun.

Fortunately, I came across this excellent video on how to fix this: [youtube.com/watch?v=W-xeqMR47NA](https://www.youtube.com/watch?v=W-xeqMR47NA)

I had to patch a file (`DWIN_SET/T5LCFG_CR10Smart.CFG`) in the touch screen firmware using "DWIN DGUS" to fix this issue.

You can find a link to the tool in the description area of the YouTube video. Or, you can search for `dgus tool v7.618` on the Internet.

Unfortunately, the tool is a Windows program. I was able to use it from within a virtualized installation of Windows.

Once you have downloaded DGUS, extract the ZIP, and start `DGUS Tool V7.618.exe`.

{{< box >}}

If the tool starts in a non-English language, and you want to change the language to English, navigate to the last tab of the ribbon menu and change the language from the dropdown.

{{< image src="language.png" alt="Screenshot of the last tab on DWIN DGUS ribbon menu" caption="Last tab on DWIN DGUS ribbon menu" >}}

{{< /box >}}

Next, click the "Config Generator" link within the "DGUS config tool" box.

{{< image src="start.png" alt="Screenshot of DWIN DGUS start page" caption="DWIN DGUS start page" >}}

A new window will open where you can modify `DWIN_SET/T5LCFG_CR10Smart.CFG`.

{{< image src="config.png" alt="Screenshot of DWIN DGUS Config Generator" caption="DWIN DGUS Config Generator" >}}

Click on the "Open CFG" button and open the `T5LCFG_CR10Smart.CFG` from inside the `DWIN_SET` directory of the firmware.

It will show the values stored in the `.CFG` file.

I changed the values for `X-axis Data Opt` and `Y-axis Data Opt`. The firmware came with these two dropdowns set to "Xmax-0" and "Ymax-0".

I changed these to "0-Xmax" and "0-Ymax" and clicked on the "Save CFG" button to save the changes back to the same file.

{{< image src="fix.png" alt="Screenshot of DWIN DGUS Config Generator after changes" caption="DWIN DGUS Config Generator after changes" >}}

With this modified configuration file, I loaded it onto a micro SD card and flashed the touch screen.

And now, the touch screen is oriented with the content and the tap areas are no longer mirrored across both axes.
