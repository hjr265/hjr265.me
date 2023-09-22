---
title: Making an Auto Scroll Bookmarklet
date: 2023-09-22T18:05:00+06:00
tags:
  - Bookmarklet
  - JavaScript
  - 100DaysToOffload
---

We have extensions and addons that add features and customizations to our web browsers. But it wasn't always the case.

During the the early 2000s (and even later) there were bookmarklets. Bookmarklets are short JavaScript scripts stored as bookmarks in your web browser.

These were popular during the 2000s. However they started to decline in popularity when browser extensions became easier to build and more and more websites started to implement Content Security Policy (CSP).

## Example from Wikipedia

The Wikipedia entry on bookmarklets provides an [example](https://en.wikipedia.org/wiki/Bookmarklet#Example):

{{< bookmarklet label="Search on Wikipedia" >}}
javascript:(function(document) {
function se(d) {
    return d.selection ? d.selection.createRange(1).text : d.getSelection(1);
}; 
let d = se(document); 
for (i=0; i<frames.length && (d==document || d=='document'); i++) d = se(frames[i].document); 
if (d=='document') d = prompt('Enter%20search%20terms%20for%20Wikipedia',''); 
open('https://en.wikipedia.org' + (d ? '/w/index.php?title=Special:Search&search=' + encodeURIComponent(d) : '')).focus();
})(document);
{{< /bookmarklet >}}

To use this bookmarklet, drag and drop the link above on your web browser's bookmark bar. Then highlight any text on this page, and click on the bookmark that you just created.

Wikipedia provides the following source code for this example:

``` js
javascript:(function(document) {
function se(d) {
    return d.selection ? d.selection.createRange(1).text : d.getSelection(1);
}; 
let d = se(document); 
for (i=0; i<frames.length && (d==document || d=='document'); i++) d = se(frames[i].document); 
if (d=='document') d = prompt('Enter%20search%20terms%20for%20Wikipedia',''); 
open('https://en.wikipedia.org' + (d ? '/w/index.php?title=Special:Search&search=' + encodeURIComponent(d) : '')).focus();
})(document);
```

This bookmarklet will use any text that you have currently highlighted on the page and will perform a Wikipedia search with it.

## Making a Bookmarklet

While preparing the blog post I published before this one, I had to capture the screen while I scrolled down a long page.

So I wrote a JavaScript `class`:

``` js
class Scroller {
  constructor(el, scrollSpeed) {
    this.el = el
    this.scrollSpeed = scrollSpeed
    this.step = this.step.bind(this)
    this.stop = this.stop.bind(this)
  }
  
  start() {
    if (this.running) return
    this.running = true
    this.scrollYStart = window.scrollY
    this.el.requestAnimationFrame(this.step)
    this.el.addEventListener('click', this.stop)
    this.el.addEventListener('wheel', this.stop)
    this.el.addEventListener('touchmove', this.stop)
    this.el.addEventListener('mousedown', this.stop)
  }
  
  step(ts) {
    if (this.tsStart === undefined) this.tsStart = ts;
    const elapsed = ts - this.tsStart;
    
    if (!this.running || this.tsPrevious === ts) return;
    this.tsPrevious = ts

    const top = this.scrollYStart + Math.min(this.scrollSpeed * elapsed, document.body.scrollHeight);
    window.scrollTo({
      top: top,
      behavior: 'instant'
    })

    if (top !== document.body.scrollHeight) {
      window.requestAnimationFrame(this.step)
    } else {
      this.stop()
    }
  }
  
  stop() {
    this.running = false
    delete this.tsStart
    delete this.scrollYStart
    this.el.removeEventListener('click', this.stop)
    this.el.removeEventListener('wheel', this.stop)
    this.el.removeEventListener('touchmove', this.stop)
    this.el.removeEventListener('mousedown', this.stop)
  }
}
```

The `Scroller` class can be used to scroll down the page automatically at a configurable speed. But it will also stop auto-scrolling if you click on the page, try to use the mouse scroll wheel, or drag the page using touch.

You can use the `Scroller` like so:

``` js
new Scroller(window, 0.025).start()
```

To convert this into a bookmarklet, I first minified the `class` code using a JavaScript code minifier.

Then I wrapped the minified `class` code and the line to use it inside an immediately invoked function expression (IIFE) and prefixed it with `javascript:`:

``` js
javascript:(function() {
  class Scroller { /* ... */ };
  new Scroller(window, 0.025).start()
})()
```

Then I added that to my bookmarks and gave it a well-thought-out name: Auto Scroll.

## And Voila!

I can now use this bookmarklet on any page to begin auto-scrolling to the end of the page. And I can stop the auto-scrolling by clicking on the page or trying to scroll myself.

Here's the bookmarklet:

{{< bookmarklet label="Auto Scroll" >}}
javascript:(function() {
  class Scroller{constructor(t,s){this.el=t,this.scrollSpeed=s,this.step=this.step.bind(this),this.stop=this.stop.bind(this)}start(){this.running||(this.running=!0,this.scrollYStart=window.scrollY,this.el.requestAnimationFrame(this.step),this.el.addEventListener("click",this.stop),this.el.addEventListener("wheel",this.stop),this.el.addEventListener("touchmove",this.stop),this.el.addEventListener("mousedown",this.stop))}step(t){void 0===this.tsStart&&(this.tsStart=t);const s=t-this.tsStart;if(!this.running||this.tsPrevious===t)return;this.tsPrevious=t;const e=this.scrollYStart+Math.min(this.scrollSpeed*s,document.body.scrollHeight);window.scrollTo({top:e,behavior:"instant"}),e!==document.body.scrollHeight?window.requestAnimationFrame(this.step):this.stop()}stop(){this.running=!1,delete this.tsStart,delete this.scrollYStart,this.el.removeEventListener("click",this.stop),this.el.removeEventListener("wheel",this.stop),this.el.removeEventListener("touchmove",this.stop),this.el.removeEventListener("mousedown",this.stop)}};
  new Scroller(window, 0.025).start()
})()
{{< /bookmarklet >}}

And, here's the final code:

``` js
javascript:(function() {
  class Scroller{constructor(t,s){this.el=t,this.scrollSpeed=s,this.step=this.step.bind(this),this.stop=this.stop.bind(this)}start(){this.running||(this.running=!0,this.scrollYStart=window.scrollY,this.el.requestAnimationFrame(this.step),this.el.addEventListener("click",this.stop),this.el.addEventListener("wheel",this.stop),this.el.addEventListener("touchmove",this.stop),this.el.addEventListener("mousedown",this.stop))}step(t){void 0===this.tsStart&&(this.tsStart=t);const s=t-this.tsStart;if(!this.running||this.tsPrevious===t)return;this.tsPrevious=t;const e=this.scrollYStart+Math.min(this.scrollSpeed*s,document.body.scrollHeight);window.scrollTo({top:e,behavior:"instant"}),e!==document.body.scrollHeight?window.requestAnimationFrame(this.step):this.stop()}stop(){this.running=!1,delete this.tsStart,delete this.scrollYStart,this.el.removeEventListener("click",this.stop),this.el.removeEventListener("wheel",this.stop),this.el.removeEventListener("touchmove",this.stop),this.el.removeEventListener("mousedown",this.stop)}};
  new Scroller(window, 0.025).start()
})()
```
