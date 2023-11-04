---
title: 'Adding Anchor Links Next to Headings in Hugo'
date: 2023-11-03T10:00:00+06:00
tags:
  - Hugo
  - 100DaysToOffload
---

I like that the Markdown renderer in Hugo automatically adds an `id` attribute to the headings in the content.

This allows you to link to a specific section in a long article.

But, I wanted to make it easy for people to get that link.

Hugo doesn't do that by default, but makes it very easy to do with Markdown render hooks.

By using the following as the render hook for headings, I am able to show a small link icon next to the headings in my blog posts:

``` html
<h{{.Level}} id="{{.Anchor | safeURL}}">
	{{.Text | safeHTML}}
	{{$linksvg := resources.Get "link.svg"}}
	<a class="hlink" href="#{{.Anchor | safeURL}}"><img src="{{$linksvg.RelPermalink}}" alt="" style="width: 1rem;"></a>
</h{{.Level}}>
```

The markup must be saved as `layouts/_default/_markup/render-heading.html` for it to be used as the heading render hook.

And with a little bit of CSS, I can have the icon hidden until I actually hover the header.

``` css
article .hlink {
    display: none;
}
article h1:hover .hlink,
article h2:hover .hlink,
article h3:hover .hlink {
    display: inline-block;
}
```

You can see this in an action in one of my longer articles that has plenty of headings: [Hiding Files in ZIP Archives](/blog/hiding-files-in-zip-archives/).
