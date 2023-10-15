---
title: 'Building a Blog With Hugo'
date: 2023-10-09T10:00:00+06:00
tags:
  - Hugo
  - 100DaysToOffload
toc: yes
---

A few weeks ago, I came across a blog post on my RSS reader: [Building a blog in Django](https://til.simonwillison.net/django/building-a-blog-in-django).

This simple tutorial by Simon Willison was an enjoyable read.

I understand the topic is not novel, but I liked how it reminds you of the simple features that enhance the implementation of a blog.

Quoting from the blog post by Simon Willison:

> Here are the features I consider to be essential for a blog in 2023 (though they haven't changed much in over a decade):
>
> - Blog posts have a title, summary, body and publication date. Optional: author information, tags
> - Posts can be live or draft
> - The blog index page shows the most recent entries
> - Older entries are available via some kind of archive mechanism
> - The blog has an Atom feed
> - Entries have social media card metadata, to enhance links to them on Mastodon and Twitter
> - Markdown is a nice-to-have for editing the posts

In this blog post, I am doing a similar exercise but using Hugo.

This blog post is not a comparison between Django and Hugo in building blogs, but it is about how to implement these essential features into a blog you build with Hugo.

## Title, Summary, Body and Publication Date

Hugo [front matter](https://gohugo.io/content-management/front-matter/) is where blog post metadata goes. 

Hugo has several predefined front matter variables, including `title`, `summary`, and publication date (as `date`).

``` md
---
title: 'Building a Blog With Hugo'
date: 2023-10-09T10:00:00+06:00
summary: '...'
---

...
```

You can access these variables directly on the `Page` type from your templates.

``` html
<ul>
  {{range .Pages}}
    <li>
      <h2><a href="{{.RelPermalink}}">{{.Title}}</a></h2>
      {{.Summary}}
    </li>
  {{end}}
</ul>
```

``` html
<h1>{{.Title}}</h1>
<div>{{.PublishDate.Format ""}}</div>
<div>
  {{.Content}}
</div>
```

## Author Information

Hugo does not have predefined front matter variables for author information. But you can add anything to the front matter and access it from the `Params` variable on the `Page` type.

``` md
---
...
author:
  name: 'Mahmud Ridwan'
  url: 'https://hjr265.me'
---

...
```

``` html
<h1>{{.Title}}</h1>
<div>{{.PublishDate.Format ""}} | <a href="{{.Params.author.url}}">{{.Params.author.name}}</a></div>
<div>
  {{.Content}}
</div>
```

## Tags

Hugo, by default, defines [two taxonomies](https://gohugo.io/content-management/taxonomies/): categories and tags.

All you need to do is add a template for the `tags` taxonomy in `layout/taxonomy/tag.html`:

``` html
<h1>Blog Posts Tagged {{.Title}}</h1>
<ul>
  {{range .Pages}}
    <li>
      <h2><a href="{{.RelPermalink}}">{{.Title}}</a></h2>
      {{.Summary}}
    </li>
  {{end}}
</ul>
```

## Draft Posts

This one is easy with Hugo. Set `draft: yes` on the front matter of the post.

``` md
---
...
draft: yes
---

...
```

While working locally on your site, you can start `hugo serve` with the `-D` flag to build the blog with the draft pages.

``` sh {linenos=false}
hugo serve -D
```

Hugo will not build pages with a future publication date unless you use the `-F` flag.

The following command will build all the pages, including draft and future ones:

``` sh {linenos=false}
hugo serve -DF
```

## Recent Posts on Homepage

To put all the pages on the homepage, you need to render them like so:

``` html
<ul>
  {{range .Pages}}
    <li>
      <h2><a href="{{.RelPermalink}}">{{.Title}}</a></h2>
      {{.Summary}}
    </li>
  {{end}}
</ul>
```

If you only want the recently published blog posts, say, from the last 30 days:

``` html
<ul>
  {{range where (.Pages) ".PublishDate" ">" (now.AddDate -30 0 0)}}
    <li>
      <h2><a href="{{.RelPermalink}}">{{.Title}}</a></h2>
      {{.Summary}}
    </li>
  {{end}}
</ul>
```

## Archival Mechanism

To make a dedicated archives page, create a directory named "archives" under the `content/` directory and then put an empty "\_index.md" file in it.

It will allow you to have a page with a separate layout under the `/archives/` URL path.

Define the layout in the `layouts/_default/archives.html` file as follows:

``` html
{{range .Site.RegularPages.GroupByDate "2006-01"}}
  <h3>{{.Key}}</h3>
  <ul>
    {{range .Pages}}
      <li>
        <div><a href="{{.RelPermalink}}">{{.Title}}</a></div>
        {{.Summary}}
      </li>
    {{end}}
  </ul>
{{end}}
```

All the blog posts will be listed on this page and will be grouped by month.

## Atom Feed

Hugo will build an RSS feed for you by default. You can [configure the RSS generated] as you see fit and even provide a custom template for it.

For Atom, you can add a theme component like [github.com/kaushalmodi/hugo-atom-feed](https://github.com/kaushalmodi/hugo-atom-feed).

1. Add the theme component to your site's configuration file.

    ``` toml
    [module]
      [[module.imports]]
        path = "github.com/kaushalmodi/hugo-atom-feed"
    ```

2. Add “ATOM” to all the `Page` kinds for which you want to create Atom feeds. In our case, since our site is the blog, we have to add it for `home`:

    ``` toml
    [outputs]
      home = ["HTML", "RSS", "ATOM"]
    ```

The Atom feed for your blog will be available directly at the '/atom.xml' URL path.

## Social Media Cards

Inside the `<head>` tag of your template HTML files, add the following internal partials:

``` html
<head>
  ...
  {{template "_internal/opengraph.html" .}}
  {{template "_internal/twitter_cards.html" .}}
  ...
</head>
```

This template will use information from the front matter for the metadata of the social cards.

It will use the first 6 URLs from the `images` variable from the front matter as the images for the social media card. If not set, it will automatically use images named "*feature*", "*cover*", or "*thumbnail*".

You can configure other site-specific metadata for these partials in your site's configuration file as detailed [here](https://gohugo.io/templates/internal/#open-graph) and [here](https://gohugo.io/templates/internal/#twitter-cards).

## Markdown

Well, that is what Hugo sites are built around.

## Wrap Up

This blog post is, again, not a comparison of Django and Hugo. Hugo is a static site generator and comes built with features that make it easy to build a blog, a personal site, or anything that fits the static site pattern.

If you are looking for a static site generator and are reading this blog post ([which is built with Hugo](https://github.com/hjr265/hjr265.me), by the way), then I hope it will encourage you to give Hugo a try.
