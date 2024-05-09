---
title: "How Hard Can It Be: Use JavaScript to Close Web Browser Window After Print"
date: 2024-05-09T12:00:00+06:00
tags:
  - WebBrowser
  - JavaScript
  - HowHardCanItBe
---

So here is a _simple_ JavaScript task I had to tackle for [Toph]:

When a user clicks the Print button, open a new tab/window and activate the print dialog. Close the window when the user confirms the print or cancels the dialog.

## Attempt 1: Call `window.close()` Immediately After `window.print()`

I linked the Print button to a separate page. I used `target="_blank"` on the link so the page opens in a new tab/window.

``` html
<a href="..." target="_blank">Print</a>
```

And on that other page I tucked away a `<script>` element right before closing the `<body>`:

``` js
<script type="text/javascript">
    window.print()
    window.close()
</script>
```

Did that work? It did on Firefox.

On Google Chrome, the Internet Explorer of the current century? It closed the window as soon as it opened.

Apparently, `window.print()` on Firefox is blocking, but not on Google Chrome and its progeny.

## Attempt 2: Use Event 'afterprint'

Surely I was not using the _standard_ way of doing things, so attempt 1 only tells you how little I know about front-end JavaScript.

I changed the JavaScript to this:

``` js
<script type="text/javascript">
    window.addEventListener("afterprint", function(event) {
        window.close()
    })
    window.print()
</script>
```

Now it works on Google Chrome.

But guess what? Not on Firefox. However, at least this time it broke in a way not as frustrating as Google Chrome.

On Firefox, the event "afterprint" is triggered at a time when `window.close()` fails to close the window.

## Attempt 3: `setTimeout(..., 0)`

I don't know why I tried this. It is like my brain is used to this idea that if something doesn't work in JavaScript, you probably need to just run it on the next tick.

``` js
<script type="text/javascript">
    window.addEventListener("afterprint", function(event) {
        setTimeout(function() { window.close() }, 0)
    })
    window.print()
</script>
```

And yes, doing this did the trick. It now works on Google Chrome and the web browser _funded_ by Google.

Any other important web browser I am missing?
