---
title: Serving JSON in Go with http.ServeContent
date: 2023-11-15T10:00:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
toc: yes
---

I know many will start with something like [Gin](https://gin-gonic.com/) whenever they are working on a JSON/HTTP-based backend in Go.

I, not entirely sure if the minority, try to stick to Go's built-in `net/http` package and, at most, use [Gorilla Mux](https://pkg.go.dev/github.com/gorilla/mux) in most of my Go projects.

And so serving something simple like JSON is no different from the package's point of view as any other content type: whatever it is, write it out to the `w`, the `http.ResponseWriter`.

But that means there are a few things worth remembering when serving JSON responses in Go over HTTP. And I am going to go over a few in this blog post.

## The Basic

``` go
// Error handling omitted for brevity.

type Message struct {
  Text string    `json:"text"`
  TS   time.Time `json:"ts"`
}

var msg = Message{
  Text: "Hello world",
  TS:   time.Date(2023, 11, 14, 10, 0, 0, 0, time.UTC),
}

func serveMessage(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(msg)
}
```

The `serveMessage` handler function uses the `json.NewEncoder` function to create a new encoder around `w` and encode the message.

You can expect the response body to look like this (minus indentations, line breaks and spaces):

``` json {linenos=false}
{
  "text": "Hello world",
  "ts": "2023-11-14T10:00:00Z"
}
```

But is that all?

## Content-Type Header

If you are serving JSON, you should set the [`Content-Type` header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type) to "application/json":

``` go
func serveMessage(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(msg)
}
```

What if we want the browser to cache the response?

## Last-Modified Header

It is where the [`http.ServeContent`](https://pkg.go.dev/net/http#ServeContent) function becomes useful.

``` go
func serveMessage(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "private, max-age=0, must-revalidate")
  b, _ := json.Marshal(msg)
  http.ServeContent(w, r, "", msg.TS, bytes.NewReader(b))
}
```

The `http.ServeContent` function will automatically add the [`Last-Modified` header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Last-Modified) to the response.

It will also check for the `If-Modified-Since` header in the request body. If present, and if it is not older than the `time.Time` passed as the `modtime` argument, the response may not include the content (and have the status `304 Not Modified` instead).

The exact headers you want to use depend on how you want the responses cached by the web browser. The example above will cause the browser to revalidate the request every time.

If the timestamp in the `TS` field doesn't change, the web browser will not have to redownload the response body.

## ETag Header

The `http.ServeContent` function also handles `If-Match`, `If-None-Match`, or `If-Range` headers. These work if you set the [`ETag` header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag) in your response:

``` go
func serveMessage(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "private, max-age=0, must-revalidate")
  b, _ := json.Marshal(msg)
  w.Header().Set("ETag", fmt.Sprintf(`"%d-%x"`, len(b), sha1.Sum(b)))
  http.ServeContent(w, r, "", msg.TS, bytes.NewReader(b))
}
```

We are generating a weak entity tag here by concatenating the length and hash of the JSON response payload.

## What Else

What else should one keep in mind when serving JSON responses in Go? Let me know if I missed something.
