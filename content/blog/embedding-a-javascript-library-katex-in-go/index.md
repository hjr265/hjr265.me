---
title: Embedding a JavaScript Library (KaTeX) in Go
date: 2022-12-24T10:35:00+06:00
tags:
  - Go
  - JavaScript
  - 100DaysToOffload
  - KaTeX
---

KaTeX is the fan-favourite way of doing math and equations on the web. The library is simple to use and easy to reason about. But it is JavaScript. How do you build a Go program that renders math and equations like KaTeX?

You embed KaTeX in your Go program.

Several Go packages allow you to run JavaScript from within Go. In this post, we will use [github.com/lithdew/quickjs](https://github.com/lithdew/quickjs).

Let us first fetch the latest copy of [katex.min.js from CDNJS](https://cdnjs.com/libraries/KaTeX). KaTeX uses the MIT license; make sure to include that in your project.

Next, in your Go program, embed katex.min.js:

``` go
package katex

import (
  _ "embed"
)

//go:embed katex.min.js
var code string
```

And add a Render function that:

- Evaluates the embedded katex.min.js code
- Sets the LaTeX string as a global variable.
- Calls `katex.renderToString` on the variable with the LaTeX string.

``` go
func Render(src []byte, display bool) (string, error) {
  // Force QuickJS to run on the same thread.
  runtime.LockOSThread()
  defer runtime.UnlockOSThread()

  // Create a new QuickJS runtime. Free it after use.
  runtime := quickjs.NewRuntime()
  defer runtime.Free()

  // Create a new QuickJS context. Free it after use.
  context := runtime.NewContext()
  defer context.Free()

  globals := context.Globals()

  // Evaluate the katex.min.js code.
  result, err := context.Eval(code)
  if err != nil {
    return "", err
  }
  defer result.Free()

  // Set the LaTeX string to a global variable and call katex.renderToString on it.
  globals.Set("latexSrc", context.String(string(src)))
  if display {
    result, err = context.Eval("katex.renderToString(latexSrc, { displayMode: true })")
  } else {
    result, err = context.Eval("katex.renderToString(latexSrc)")
  }
  if err != nil {
    return "", err
  }
  defer result.Free()

  // Return the rendered equation.
  return result.String(), nil
}
```

The above code gives you the most straightforward implementation. You may want to rewrite it in a way where you evaluate katex.min.js only once and keep using the same `context` for every invocation of `Render`.

Using JavaScript libraries from within Go made packages like [goldmark-katex](https://github.com/FurqanSoftware/goldmark-katex) possible. Of course, for convenience, you are trading a bit of performance (because you are evaluating JavaScript from within Go). But sometimes, that convenience makes it worth it. For goldmark-katex, it was either this or building a KaTeX clone in Go.
