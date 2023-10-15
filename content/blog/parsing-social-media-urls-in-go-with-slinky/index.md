---
title: 'Parsing Social Media URLs in Go With Slinky'
date: 2023-10-14T10:20:00+06:00
tags:
  - Go
  - 100DaysToOffload
---

Toph now allows programmers to show up to 5 social media URLs on their profile pages.

Instead of showing the entire URL, I wanted to show the important bits from the URL.

{{< image src="screen.png" alt="Screenshot of a profile panel from Toph" caption="Screenshot of a profile panel from Toph" >}}

To do that, I had to parse the social media URLs and extract information like the username or profile ID (when it is a GitHub, Twitter, LinkedIn, Facebook, etc. profile URL) or the instance of Mastodon.

So, I built [Slinky](https://github.com/FurqanSoftware/slinky). It is a Go package that parses social media URLs into structured data.

Right now, it supports URLs from the following social media platforms:

``` go
// Facebook
"facebook.com"
"www.facebook.com"
"web.facebook.com"
"fb.me"

// FLOSS.social
"floss.social"

// Fostodon
"fosstodon.org"

// GitHub
"github.com"
"*.github.io"

// Instagram
"instagram.com"
"www.instagram.com"

// LinkedIn
"linkedin.com"
"www.linkedin.com"

// Telegram
"t.me"

// Twitter
"twitter.com"

// YouTube
"youtube.com"
"www.youtube.com"
```

When you pass a URL like "https://github.com/hjr265" to `slinky.Parse`, Slinky will parse the URL into this:

``` go
&URL{
  Service: slinky.GitHub,
  Type:    "User",
  ID:      "hjr265",
  Data:    map[string]string{,
    "username": "hjr265",
  },
}
```

In the case of a floss.social or a fosstodon.org URL, you will get a `URL` value like this:

``` go
&URL{
  Service: slinky.Fosstodon,
  Type:    "Profile",
  ID:      "hjr265",
  Data:    map[string]string{,
    "username": "hjr265",
    "platform": "Mastodon",
  },
}
```

If you need to parse social media URLs in your Go application, try [Slinky](https://github.com/FurqanSoftware/slinky).

And, if you want to extend Slinky to support additional social media platforms, feel free to open an issue with details or send a pull request.
