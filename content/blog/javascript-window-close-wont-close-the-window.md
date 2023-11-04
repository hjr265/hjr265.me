---
title: "JavaScript window.close Won't Close the Window"
date: 2023-11-04T19:00:00+06:00
tags:
  - JavaScript
  - 100DaysToOffload
---

I learned something new today.

In JavaScript within the web browser, `window.close` will not close the window if it was not opened using `window.open` or is a top-level window (or tab) with at one history entry. That is what the [documentation of `window.close`]((https://developer.mozilla.org/en-US/docs/Web/API/Window/close)) says.

It got in the way.

I was adding a page endpoint that I would link to. The link would open the page in a new tab with `target="_blank"`. And it would call `window.print` and `window.close` back to back.

The intent was so that the page loads, shows the print dialog, and closes after the user confirms or cancels the print dialog.

The solution was to use the `onclick` attribute with `window.open` in the link:

``` html {linenos=false}
<a href="javascript:;" onclick="window.open('/_/contests/{{.Contest.ID}}/participants/credential_slips#intent=print')">Print Slips</a>
```

Instead of `href`:

``` html {linenos=false}
<a href="/_/contests/{{.Contest.ID}}/participants/credential_slips#intent=print">Print Slips</a>
```

Well, at least there was a solution to this. And, `window.open` reminded me of the time when Internal Explorer was a thing and had the largest share in the web browser market. Eek!
