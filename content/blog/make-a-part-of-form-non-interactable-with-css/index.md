---
title: Make a Part of Form Non-interactable With CSS
date: 2023-09-15T17:30:00+06:00
tags:
  - CSS
  - 100DaysToOffload
---

I have an HTML form. I want to make part of it non-interactable depending on certain conditions. I don't want to remove that part entirely.

There are so many reasons why you may want to do this.

{{< image src="screen.png" alt="Screenshot of Contest Branding Settings Form from Toph" captionMD="Screenshot of Contest Branding Settings Form from [Toph](https://toph.co/)" >}}

In this screenshot, the form allows advanced features to pay customers only. Making that part non-interactable, instead of hiding it, works as a teaser of what the paid tiers offer.

The CSS for that is straightforward:

``` css
.blocked {
  -webkit-mask-image: -webkit-gradient(linear, left top, left bottom, from(rgba(0,0,0,0.5)), to(rgba(0,0,0,0)));
  mask-image: gradient(linear, left top, left bottom, from(rgba(0,0,0,0.5)), to(rgba(0,0,0,0)));
  pointer-events: none;
}
```

The important property here is `pointer-events`. Setting it to `none` causes the browser to block all pointer events for the element. Any clicks or taps on the element will not trigger any event.

The `mask-image` property applies a fading gradient effect just to give the user visual cues on what's going on.

The element may still receive focus through the tab key on the keyboard. But you can prevent that by adding `tabindex="-1"` in the right places.

Remember that this CSS makes the form non-interactable on the front end only. You must add the any necessary logic in the backend to prevent your HTML forms from being abused.
