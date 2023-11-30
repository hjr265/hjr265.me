---
title: 'Scanning a Website for Broken Links in Go'
date: 2023-11-20T23:15:00+06:00
tags:
  - Go
  - 100DaysToOffload
---

Yes, I know there are paid and free tools for doing this. And yes, I know there are tools for this that I can run locally.

But this exercise allowed me to try out the well-designed Go package [github.com/gocolly/colly](https://pkg.go.dev/github.com/gocolly/colly).

Colly is a web scraping framework for Go.

Here is how I used it to quickly scan my website (the one you are on right now) for broken links.

First I defined a type for links to check and the URL of the page they appear on:

``` go
type link struct {
  href    string
  pageURL string
}
```

I also wrote a rudimentary function to check if a link is okay:

``` go
func isLinkOkay(l link) bool {
  r, err := http.NewRequest("GET", l.href, nil)
  r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0 Config/93.2.8101.2")
  if err != nil {
    return false
  }
  resp, err := http.DefaultClient.Do(r)
  if err != nil {
    return false
  }
  defer resp.Body.Close()
  io.Copy(io.Discard, resp.Body)
  return resp.StatusCode >= 200 && resp.StatusCode < 400
}
```

Make sure to set a timeout on the default HTTP client.

``` go
http.DefaultClient.Timeout = 10 * time.Second
```

Next, let us define a worker function to check links as they are scanned from the website:

``` go
func checkWorker(checkCh chan link) {
  cache := map[string]bool{}
  for link := range checkCh {
    okay, hit := cache[link.href]
    if !hit {
      okay = isLinkOkay(link)
      cache[link.href] = okay // Cache the status of the check to avoid duplicating effort.
    }
    if !okay {
      // For now, we will only print each broken link and the URL of the page it is on.
      fmt.Println(link.href)
      fmt.Println(":", link.pageURL)
    }
  }
}
```

Finally, the function to crawl the website:

``` go
func crawl(checkCh chan link, domain string) {
  // Create a new Colly collector.
  c := colly.NewCollector(
    colly.AllowedDomains(domain), // Crawl URLs that belong to the site being scanned only.
    colly.Async(true),
  )
  c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 7}) // Allow some parallelism.

  c.OnHTML("a[href]", func(e *colly.HTMLElement) {
    href := e.Request.AbsoluteURL(e.Attr("href"))
    if href != "" && !strings.HasPrefix(href, "https://"+domain) {
      // Send each external URL over the check channel. You could send
      // internal URLs over as well if that is something you want to check as
      // well.
      checkCh <- link{href, e.Request.URL.String()}
    }
    e.Request.Visit(href) // Enqueue all URLs discovered on the pages.
  })

  c.Visit("https://" + domain) // Begin the crawl from the homepage.

  c.Wait()
}
```

And finally, the `main` function to wrap it all together.

``` go
func main() {
  http.DefaultClient.Timeout = 10 * time.Second

  checkCh := make(chan link, 500)
  go checkWorker(checkCh)

  crawl(checkCh, "hjr265.me")
}
```

And that's it. I can run this program to identify any broken links on my website.

Colly, an easy-to-use scraping framework, makes it possible to do more than just detect broken links. I can use it to perform other routine audits, like checking to see if my images are missing `alt` attributes, if my pages have the correct meta tags, and more.
