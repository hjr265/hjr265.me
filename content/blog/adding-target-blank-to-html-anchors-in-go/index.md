---
title: Adding target="_blank" to User-Generated HTML Anchors in Go
date: 2022-12-06T16:00:00+06:00
tags:
  - Go
  - HTML
  - 100DaysToOffload
---

Working with user-generated content is always ~~a nightmare~~ interesting.

Let's say you are building a blogging platform with Go. Your users write posts in Markdown that the platform then renders as HTML. And, you want to add `target="_blank"` and `rel="noreferrer noopener"` to all the external links. How do you do that?

## Annotated Code

The steps are simple:

- Parse the HTML with `golang.org/x/net/html`.
- Walk the tree. The annotated code below implements a simple `Walk` function.
- Update the nodes. Do this in the callback of the `Walk` function.
- Render the modified HTML.

``` go
package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	// Error handling omitted for brevity.

	const content = `<p>
	<a href="https://example.com">Not External</a>
	<a href="https://furqansoftware.com">External</a>
</p>`

	// Parse the HTML.
	doc, _ := html.Parse(strings.NewReader(content))

	// Walk the entire HTML tree.
	Walk(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.A && IsExternal(GetAttr(n, "href")) {
			// Set attributes on all anchors with external links.
			SetAttr(n, "target", "_blank")
			SetAttr(n, "rel", "noreferrer noopener")
		}
	})

	// Render the modified HTML.
	b := bytes.Buffer{}
	html.Render(&b, doc)
	fmt.Println(b.String())
	// Output:
	// <html><head></head><body><p>
	// 	<a href="https://example.com">Not External</a>
	// 	<a href="https://furqansoftware.com" target="_blank" rel="noreferrer noopener">External</a>
	// </p></body></html>
}

// Walk traverses the entire HTML tree and calls fn on each node.
func Walk(n *html.Node, fn func(*html.Node)) {
	if n == nil {
		return
	}
	fn(n)

	// Each node has a pointer to its first child and next sibling. To traverse
	// all children of a node, we need to start from its first child and then
	// traverse the next sibling until nil.
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		Walk(c, fn)
	}
}

// GetAttr returns the attribute on a node by its key.
func GetAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// SetAttr sets an attribute on a node.
func SetAttr(n *html.Node, key, val string) {
	for i := range n.Attr {
		a := &n.Attr[i]
		if a.Key == key {
			a.Val = val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{
		Key: key,
		Val: val,
	})
}

// IsExternal returns true if url doesn't have "https://example.com" as the
// prefix. You can do better than this.
func IsExternal(url string) bool {
	return !strings.HasPrefix(url, "https://example.com")
}
```

## Why Is The Output An Entire HTML Document?

Any HTML you parse is treated as an entire document. But more often than not, when dealing with user-generated content, you are probably dealing with an HTML fragment.

To force the final output to be a fragment (just like the input), you need to find the node for the `<body>`, turn it into a document node, and then call `html.Render` on it.

``` go
func main() {
	// ...

	// Find the <body>.
	var body *html.Node
	Walk(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.Body {
			body = n
		}
	})

	// Change the <body> element's type to DocumentNode.
	body.Type = html.DocumentNode
	body.DataAtom = 0
	body.Data = ""

	// Render the modified HTML.
	b := bytes.Buffer{}
	html.Render(&b, body)
	fmt.Println(b.String())
	// Output:
	// <p>
	// 	<a href="https://example.com">Not External</a>
	// 	<a href="https://furqansoftware.com" target="_blank" rel="noreferrer noopener">External</a>
	// </p>
}
```

## Why Not Just Modify/Extend The Markdown Renderer

Of course, you can always modify or add extensions to your Markdown renderer to emit HTML with the desired attributes. But with this approach, you can do much more than your Markdown renderer would allow you.

In Toph, we store both the Markdown and the rendered HTML in the database. But there are situations where the final HTML we present to the web browser must be modified slightly based on when it is being used. For example, when the content mentions a user with "@handle", the resulting anchor element is coloured based on the user's rating. This may change after the original Markdown is rendered and stored. With an approach similar to the one described above, we can modify the generated content on the fly before serving it to the web browser.
