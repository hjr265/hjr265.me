---
title: Showing GitHub Stars With Static Site Generator Hugo
date: 2022-11-27T20:13:00+06:00
tags:
  - hugo
  - shortcode
  - 100DaysToOffload
---

Static site generators are one of my favourite things about the Internet. I remember when almost every website built around me was based on Joomla or WordPress. I dread that time.

My website, which you are on right now, is built with Hugo. I have a page on this website listing some of [my open-source projects](/open-source/). And I wanted an easy way to show the number of GitHub stars on my Hugo-based website for my open-source projects.

Something like this (e.g. for my [Go implementation of Redlock](https://github.com/go-redsync/redsync)):

{{< githubstars "go-redsync/redsync" >}}

The number is not hard coded. It is fetched from a GitHub API every time I regenerate my website. And for that, I can either use a shortcode (if it is a part of the content):

``` txt
{{</* githubstars "go-redsync/redsync" */>}}
```

Or, a partial (if it is a part of the layout):

``` txt
{{partial "githubstars" "go-redsync/redsync"}}
```

Implementing them both is pretty straightforward.

First, implement a partial template like so:

``` html
{{with getJSON "https://api.github.com/repos/" .}}
  {{$starsvg := resources.GetMatch "star.svg" | fingerprint}}
  <img src="{{$starsvg.RelPermalink}}" style="height: 1em; vertical-align: text-top;"> {{.stargazers_count}}
{{end}}
```

This partial template fetches the repository data from GitHub API. The star count is in the `stargazers_count` field of the response JSON.

The template adds "star.svg" from your "assets" directory to the left of the count.

Add this template to the "layouts/partials" directory.

This provides the partial template that you can now use in your layouts. And you can also use this in your shortcodes:

``` txt
{{partial "githubstars" (.Get 0)}}
```

Add this shortcode to the "layouts/shortcodes" directory.

And that's it! You can show your (or any) GitHub repository star count on your Hugo website.

---

_This post is 2nd of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
