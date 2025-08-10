---
title: "7 Useful io.Reader and io.Writer Wrappers in Go"
date: 2023-11-17T09:45:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
draft: yes
---

Whenever I think about Go in comparison to other programming languages, the first thing that comes to my mind is how Go simplifies the concepts around concurrency and writing concurrent programs.

After all, concurrency is one of the features touted on the homepage:

> Built-in concurrency and a robust standard library

But there is something else that Go simplifies and makes it very easy to wrap one's head around.

The type system.

For anyone who is reaching for your pitch and fork, hear me out.

Take `io.Reader` as an example. Anything that is an `io.Reader` can be read from using the `Read` method. And, all the utilities built around `io.Reader` will work with anything that complies with the specifications of an `io.Reader`: not just what has been built but what will be built in the future as well.

And that is why you can create custom wrapper types for `io.Reader` and `io.Writer` that are themselves `io.Reader` and `io.Writer` to do all sorts of useful things.

In this blog post I will share 7 useful `io.Reader` and `io.Writer` wrappers that I found useful in the Go projects that I have worked on.

## Detect Line Endings

## Word Counter
