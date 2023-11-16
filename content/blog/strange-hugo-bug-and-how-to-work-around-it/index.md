---
title: Strange Hugo Bug and How to Work Around It
date: 2023-11-15T10:45:00+06:00
tags:
  - Hugo
  - 100DaysToOffload
---

I was about to deploy my site with the latest blog post this morning and found Hugo broken. Or, my Hugo `config.toml` is broken. It depends on how you want to look at it.

I ran `hugo deploy` and bam! I see an internal template error.

``` txt {linenos=false}
ERROR render of "taxonomy" failed: template: _internal/_default/rss.xml:3:9: executing "_internal/_default/rss.xml" at <site>: can't evaluate field email in type string
```

According to `pacman`, the package manager on Arch Linux, I received an update for Hugo yesterday.

The error message indicated that the built-in `rss.xml` template is trying to use a field named email.

So I looked up the [`_default/rss.xml`](https://github.com/gohugoio/hugo/blob/d4016dd5cd57a27f19a5472c6031d156066860b7/tpl/tplimpl/embedded/templates/_default/rss.xml) file on GitHub.

On line 3, Hugo is interpolating the field `email` from `params.author` in `config.toml`.

But my `params.author` is a string, my name, not an object.

I noticed a contributor to the project already added a fix for this to the Hugo repository since then. But the version I have installed on my computer does not have the fix. It has the bug.

To work around this issue, I had to update the `config.toml` file to have an object for the `params.author` field:

``` toml
[params]
github = "hjr265"

[params.author]
name  = "Mahmud Ridwan"
email = "m@hjr265.me"
```

Instead of:

``` toml
[params]
author = "Mahmud Ridwan"
github = "hjr265"
```
