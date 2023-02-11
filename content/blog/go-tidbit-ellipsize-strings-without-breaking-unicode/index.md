---
title: 'Go Tidbit: Ellipsize Strings Without Breaking Unicode'
date: 2023-02-06T21:00:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

Strings in Go support Unicode. If you come from a programming language like C, you may think of strings as an array of (byte-sized) characters. 

In Go, you can convert a string to a byte slice, and access/manipulate each byte. But if you want to truncate or ellipsize the string to a specific length, you have to think of a string like a slice of runes.

Take this string for example: `"বাংলাদেশ"`.

If you want to ellipsize it to a length of `n` characters, you may run a code like this:

``` go
func ellipsis(s string, n int) string {
	if len(s) > n {
		return s[:n]+"…"
	}
	return s
}

func main() {
	var s = "বাংলাদেশ"
	for n := 1; n <= 8; n++ {
		fmt.Println(ellipsis(s, n))
	}
	// Output:
	// �…
	// ��…
	// ব…
	// ব�…
	// ব��…
	// বা…
	// বা�…
	// বা��…
}
```

Huh! Something needs to be fixed. What is the sane and desired output is:

``` txt
ব…
বা…
বাং…
বাংল…
বাংলা…
বাংলাদ…
বাংলাদে…
বাংলাদেশ
```

To correctly truncate a `string` with Unicode characters, convert the `string` to a slice of runes first:

``` go
func ellipsis(s string, n int) string {
	r := []rune(s) // Convert s to a slice of runes and truncate it instead.
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}

func main() {
	var s = "বাংলাদেশ"
	for n := 1; n <= 8; n++ {
		fmt.Println(ellipsis(s, n))
	}
	// Output:
	// ব…
	// বা…
	// বাং…
	// বাংল…
	// বাংলা…
	// বাংলাদ…
	// বাংলাদে…
	// বাংলাদেশ
}
```

<br>

_This post is 15th of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
