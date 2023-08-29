---
title: Hugo Footnote for the 100 Days to Offload Challenge
date: 2023-07-01T19:50:00+06:00
tags:
  - Hugo
  - 100DaysToOffload
---

Yesterday I posted my 25th blog post for the [#100DaysToOffload](/tags/100daystooffload/). That's 25% of the challenge.

If it wasn't clear by the post [Showing GitHub Stars With Static Site Generator Hugo](/blog/showing-github-stars-with-static-site-generator-hugo/), I use Hugo for this site.

All this time, I was manually adding a footnote to each of the blog posts:

> This post is {n}th of my #100DaysToOffload challenge. Want to get involved? Find out more at 100daystooffload.com.

Figuring out what `{n}` is for each blog post wasn't fun.

I replaced that with a simple snippet of Hugo template code:

``` go-html-template
{{if in .Params.tags "100DaysToOffload"}}
	{{$nth := 0}}
	{{$pages := where .Site.Pages "Params.tags" "intersect" (slice "100DaysToOffload")}}
	{{range $i, $_ := $pages}}
		{{if eq .Permalink $.Permalink}}
			{{$nth = sub (len $pages) $i}}
		{{end}}
	{{end}}
	<br>
	<p><em>This post is {{$nth | humanize}} of my <a href="/tags/100daystooffload/">#100DaysToOffload</a> challenge. Want to get involved? Find out more at <a href="https://100daystooffload.com/" target="_blank" rel="noreferrer noopener">100daystooffload.com</a>.</em></p>
{{end}}
```

This template, in the context of a blog post, checks to see if the blog post has the 100DaysToOffload tag. If it does, we look at all the blog posts with the same tag and determine the index of the current blog post among the list.

The template uses the `humanize` function to display the n-th ordinal.

A preview of what the end result looks like should be just below.
