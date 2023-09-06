---
title: 'Client-side Search in Hugo with Fuse.js'
date: 2023-09-06T19:20:00+06:00
tags:
  - Hugo
  - FuseJS
  - 100DaysToOffload
  - Search
---

[Toph Help](https://help.toph.co/) is built with Hugo - a static site generator.

As you would expect with static sites, the pages are all generated ahead of time and hosted as plain HTML. You get all the benefits of static websites, but what about search?

Client-side search is one way to work around this limitation of static websites.

You build an array of objects describing all your pages on your website. You serve it to the client as JSON. You use the JSON with client-side JavaScript to provide fast search functionality.

And, it works quite well in practice.

## Generating the Search Index

Hugo allows you to produce different types of output files, not just HTML.

You can add an `index.json` file to your `layouts/` directory to build the search index as a part of the site generation.

``` go-html-template
{{- $data := newScratch -}}
{{- $data.Add "index" slice -}}
{{- range .Site.RegularPages -}}
    {{- $contents := .Plain -}}
    {{- range .Resources.Match "*.md" -}}
        {{- $contents = printf "%s\n%s\n%s" $contents .Title .Plain -}}
    {{- end -}}
    {{- $section := .CurrentSection.Title -}}
    {{- if .CurrentSection.Parent -}}
        {{- if not .CurrentSection.Parent.IsHome -}}
            {{- $section = printf "%s / %s" .CurrentSection.Parent.Title $section -}}
        {{- end -}}
    {{- end -}}
    {{- $data.Add "index" (dict "title" .Title "section" $section "tags" .Params.tags "categories" .Params.categories "contents" $contents "permalink" .Permalink) -}}
{{- end -}}
{{- $data.Get "index" | jsonify (dict "indent" "  ") -}}
```

This template loops over all the regular pages on your site and prepares an array of objects describing the pages.

The output JSON looks something like this:

``` json
[
  {
    "categories": null,
    "contents": "To create a contest on [...] programming problems.\n",
    "permalink": "https://help.toph.co/toph/hosting-a-contest/create-a-contest/",
    "section": "Hosting a Contest",
    "tags": null,
    "title": "Create a Contest"
  },
  {
    "categories": null,
    "contents": "When you enter the [...] by the organizers.\n",
    "permalink": "https://help.toph.co/toph/contest-arena/dashboard/",
    "section": "Contest Arena",
    "tags": null,
    "title": "Dashboard"
  },
  // ...
]
```

## Search with Fuse.js

[Fuse.js](https://www.fusejs.io/) is a lightweight dependency-free fuzzy-search library.

On the front end, we use it to build the search index:

``` js
let fuse

fetch('/index.json')
  .then(resp => resp.json())
  .then(index => {
    fuse = new Fuse(index, {
      shouldSort: true,
      includeMatches: true,
      threshold: 0.0,
      tokenize: true,
      location: 0,
      distance: 100,
      maxPatternLength: 32,
      minMatchCharLength: 1,
      keys: [
        {name: 'title', weight: 0.8},
        {name: 'contents', weight: 0.5},
        {name: 'tags', weight: 0.3},
        {name: 'categories', weight: 0.3}
      ]
    })
  })
```

We can now use `fuse.search()` to find matching documents from the built index:

``` js
let results = fuse.search(query)
```

## Wrap Up

And, with that, you have client-side search functionality.

You will now need a way to render the results on the search results page or below the fancy search input field. But how you do it depends a lot on the look and feel of your static website.
