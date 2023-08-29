---
title: 'Go Tidbit: Peek a Snippet From io.Reader'
date: 2023-02-05T13:40:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

You have an `io.Reader`, and you want to extract a small snippet from the beginning of the `io.Reader` and then put it back.

Ideally, `io.ReadSeeker` would let you do that. But not all `io.Reader`s can seek.

And, like toothpaste from its tube, you cannot put something back once you read it from an `io.Reader`. But you can do the next best thing:

``` go
// PeekSnippet takes r and returns the first n bytes from it and another
// io.Reader with the entire data.
func PeekSnippet(r io.Reader, n int) ([]byte, io.Reader, error) {
	lr := io.LimitReader(r, n)
	b, err := io.ReadAll(lr) // Read first n bytes from r through lr.
	if err != nil {
		return nil, nil, err
	}
	r = io.MultiReader(bytes.NewReader(b), r) // Make a new reader combining the bytes just read and the remaining data in r.
	return b, r, nil
}
```

As you can see, you can use `io.LimitReader` to read the first `n` bytes. Then use `io.MultiReader` to combine those bytes with the remainder of data in `r`.
