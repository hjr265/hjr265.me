---
title: "Firefox's Screenshot Command"
subtitle: "You can use :screenshot helper function in the web console to take screenshots in Firefox, with a bunch of useful options."
tags:
  - Firefox
  - WebBrowser
---

I needed to capture a HiDPI screenshot in Firefox without switching to a device with a HiDPI display. I was wondering if that was even possible. While searching for a how-to, I came across this `:screenshot` web console helper function.

How is it that I didn't know about it for so many years? 🤯

This function has been in [existence for about 7-8 years](https://meyerweb.com/eric/thoughts/2018/08/24/firefoxs-screenshot-command-2018/). And it has a bunch of handy options.

``` txt {linenos=false}
:screenshot --help
```

``` txt {linenos=false}
Save an image of the page
Options

     --clipboard
          type: boolean
          description: Copy screenshot to clipboard? (true/false)
          manual: True if you want to copy the screenshot instead of saving it to a file.

     --delay
          type: number
          description: Delay (seconds)
          manual: The time to wait (in seconds) before the screenshot is taken

     --dpr
          type: number
          description: Device pixel ratio
          manual: The device pixel ratio to use when taking the screenshot

     --fullpage
          type: boolean
          description: Entire webpage? (true/false)
          manual: True if the screenshot should also include parts of the webpage which are outside the current scrolled bounds.

     --selector
          type: string
          description: CSS selector
          manual: A CSS selector for use with document.querySelector which identifies a single element

     --file
          type: boolean
          description: Save to file? (true/false)
          manual: True if the screenshot should save the file even when other options are enabled (eg. clipboard).

     --filename
          type: string
          description: Destination filename
          manual: The name of the file (should have a ‘.png’ extension) to which we write the screenshot.
```

Long story short, if you want to take a HiDPI screenshot, all you need to do is:

1. Press F12 to open Developer Tools.
2. Switch to the Console tab.
3. Type `:screenshot --dpr 4`
