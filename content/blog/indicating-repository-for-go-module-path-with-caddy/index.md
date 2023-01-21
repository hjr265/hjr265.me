---
title: Using Caddy to Indicate GitLab Repository for Go Module Path with Different Domain
htmltitle: Using Caddy to Serve Canonical "go get" Import Path
date: 2023-01-21T11:05:00+06:00
tags:
  - Caddy
  - Go
  - 100DaysToOffload
  - GitLab
---

The self-hosted GitLab at my company, [Furqan Software](https://furqansoftware.com), is home to all the company Go projects. The domain where Furqan Software's GitLab is accessible is different from the domain in Go module paths. That means `go get` doesn't work out of the box for our Go projects.

For example, let's say you have a Go project with the module path `go.example.com/tools/aglet`. And the corresponding GitLab project is at `https://gitlab.example.com/tools/aglet`. If you run `go get go.example.com/tools/aglet`, you will see an error from `go get` about not finding the repository.

> When the go command downloads a module in direct mode, it starts by locating the repository that contains the module. [\[...\]](https://go.dev/ref/mod#vcs-find)

To allow `go get go.example.com/tools/aglet` to work, you need to serve a specific HTML page at `https://go.example.com/tools/aglet?go-get=1` with the right meta tag.

With Caddy, you can do this easily, although with a few limitations (more on that at the end of this post).

Caddy has built-in support for interpolating Go templates. We can leverage that to serve the HTML page with the right meta tag:

``` html
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
{{$path := .Req.URL.Query.Get "p" -}}
{{$parts := splitList "/" $path -}}
{{$repo := printf "/%s/%s" (index $parts 1) (index $parts 2) -}}
<meta name="go-import" content="go.example.com{{$repo}} git https://gitlab.example.com{{$repo}}.git">
<meta name="go-source" content="go.example.com{{$repo}} https://gitlab.example.com{{$repo}} https://gitlab.example.com{{$repo}}/-/tree/master{/dir} https://gitlab.example.com{{$repo}}/-/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://gitlab.example.com{{$repo}}">
</head>
<body>
<a href="https://gitlab.example.com{{$repo}}">Redirecting to GitLab...</a>
</body>
</html>
```

To use the template, your Caddyfile will look something like this:

``` text
go.example.com {
  root * /etc/caddy/go-get/
  file_server
  templates
  encode gzip
  try_files index.html?p={path}
}
```

With this, Caddy expects the template HTML to be present in the `/etc/caddy/go-get` directory as index.html.

And Caddy can then start serving the required HTML with the meta tag at URLs like `https://go.example.com/tools/aglet?go-get=1` for commands like `go get go.example.com/tools/aglet`.

The template here assumes that all Go modules will have GitLab URLs like `https://gitlab.example.com/{group}/{project}`. Although GitLab supports nested groups, this solution won't be able to handle it.

Alternate and better solutions exist: like [github.com/unknwon/go-import-server](https://github.com/unknwon/go-import-server). And you can always use the same domain for your self-hosted GitLab instance and Go module paths to avoid relying on this quirky workaround.

<br>

_This post is 12th of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
