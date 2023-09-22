---
title: 'HTML <img> Markup (and Hugo Shortcode) for the Modern Web'
date: 2023-09-22T10:30:00+06:00
tags:
  - HTML
  - Image
  - 100DaysToOffload
  - Hugo
toc: yes
---

I am by no means an expert in HTML and CSS. But I have been tweaking and tuning a few of my Hugo-powered websites, including this one. And, I had the opportunity to explore a few of the modern HTML and CSS features.

Starting last month I made it a point to use more images in my blog posts.

And, I am glad I did. When I wrote the [blog post on hiding files in ZIP archives](/blog/hiding-files-in-zip-archives/), I enjoyed preparing the ZIP format illustrations.

On the web, images are more than the pixels on the screens.

If you want to fine-tune your website for the best visitor experience, then you need a little more than just `<img src="...">`.

The `<img>` element is almost as old as the web itself, but the one that browsers support now has come a long way.

In this blog post, I am outlining a few things to keep in mind when implementing markup around images. You may already be familiar with some or all of these, but at least it will serve as a reminder to me for future work.

## `alt` Attribute

This isn't a modern feature. But it is an important one. If I am fine-tuning the markup of a website around images, this one should come before any modern features.

And, it is only today that I learned how the `alt` attribute wasn't a part of the `<img>` element when it was introduced.

> The `img` tag was included in the HTML 2.0 specification released in 1995 by the W3C. Some alternatives were proposed, like the `fig` tag that included the much needed `alt` attribute for users unable to see images. [\[...\]](https://thehistoryoftheweb.com/the-origin-of-the-img-tag/)

The `alt` attribute does more than just help search engines make sense of images on the web. It plays a key role in making the website more accessible.

A screen reader can read out the value of the `alt` attribute to someone who is relying on it.

It will also work as a fallback when the image fails to load.

## `<figure>` and `<figcaption>`

{{< image src="caption.png" alt="Screenshot of an image on the web with a caption" caption="Screenshot of an image on the web with a caption" >}}

Images that are a part of a page's main content should ideally have a caption. To add a caption to an image, wrap the `<img>` element inside a `<figure>` element and add a `<figcaption>`.

``` html
<figure>
  <img src="...">
  <figcaption>...</figcaption>
</figure>
```

## `loading="lazy"` Attribute

{{< video src="lazy.mp4" muted="true" >}}

Gone are the days when lazy loading images meant JavaScript code (with or without the `IntersectionObserver`)

Today, if I want to lazy load an image on a page, I can add the `loading="lazy"` attribute to my `<img>` element. All modern browsers support it.

``` html
<img src="..." loading="lazy">
```

## `<picture>` and `<source>`

There are two ways to markup responsive images:

- Use `srcset` attribute on an `<img>` element.
- Use a `<picture>` elements with multiple `<source>` elements and a fallback `<img>` element.

I am using the second option:

``` html
<picture>
  <source type="image/webp" media="(max-width: 575.98px)" srcset="A... 1x, B... 2x">
  <source type="image/webp" media="(min-width: 576px)" srcset="C... 1x, D... 2x">
  <img srcset="E... 1x, F... 2x" src="G..." alt="...">
</picture>
```

In the markup above, I am providing 7 image URLs.

There are four different images for two different screen sizes and two different pixel densities:

- `A...`: Will be shown on screens smaller than 576px and standard pixel density.
- `B...`: Will be shown on HiDPI screens smaller than 576px.
- `C...`: Will be shown on screens larger than (or equal to) 576px and standard pixel density.
- `D...`: Will be shown on HiDPI screens larger than (or equal to) 576px.

On web browsers that do not support the `<picture>` element, the `srcset` or `src` attributes on the `<img>` element will take over.

If the web browser supports the `srcset` attribute, then either `E...` or `F...` will be selected, depending on the pixel density of the screen. Otherwise, the image will fall back to the `src` attribute.

You can learn a lot more about [responsive images on MDN Web Docs](https://developer.mozilla.org/en-US/docs/Learn/HTML/Multimedia_and_embedding/Responsive_images).

## Hugo Shortcode

For images on this blog, I came up with a simple Hugo shortcode.

``` html
{{- $src := .Get "src" -}}
{{- $alt := .Get "alt" -}}
{{- $caption := .Get "caption" -}}
{{- $captionMD := .Get "captionMD" -}}
{{- $captionHTML := .Get "captionHTML" -}}

<figure>
  {{$image := .Page.Resources.Get $src}}
  <a href="{{$image.RelPermalink}}" target="_blank">
    <picture>
      {{$resized := $image}}
      {{if gt $resized.Width 750}}
        {{$resized = $image.Resize (printf "%dx" 750)}}
      {{end}}
      {{$resized2x := $image}}
      {{if gt $resized2x.Width (mul 2 750)}}
        {{$resized2x = $image.Resize (printf "%dx" (mul 2 750))}}
      {{end}}
      {{$mobile := $image}}
      {{if gt $mobile.Width 504}}
        {{$mobile = $image.Resize (printf "%dx" 504)}}
      {{end}}
      {{$mobile2x := $image}}
      {{if gt $mobile2x.Width (mul 2 504)}}
        {{$mobile2x = $image.Resize (printf "%dx" (mul 2 504))}}
      {{end}}
      {{if eq $image.MediaType.Type "image/jpeg"}}
        {{$resizedWebp := $resized.Resize (printf "%dx%d webp" $resized.Width $resized.Height)}}
        {{$resized2xWebp := $resized2x.Resize (printf "%dx%d webp" $resized2x.Width $resized2x.Height)}}
        {{$mobileWebp := $mobile.Resize (printf "%dx%d webp" $mobile.Width $mobile.Height)}}
        {{$mobile2xWebp := $mobile2x.Resize (printf "%dx%d webp" $mobile2x.Width $mobile2x.Height)}}
        <source type="image/webp" media="(max-width: 575.98px)" srcset="{{$mobileWebp.RelPermalink}} 1x, {{$mobile2xWebp.RelPermalink}} 2x">
        <source type="image/webp" media="(min-width: 576px)" srcset="{{$resizedWebp.RelPermalink}} 1x, {{$resized2xWebp.RelPermalink}} 2x">
      {{end}}
      <img srcset="{{$resized.RelPermalink}} 1x, {{$resized2x.RelPermalink}} 2x" src="{{$resized.RelPermalink}}" alt="{{$alt}}" loading="lazy">
    </picture>
  </a>
  {{if or $captionHTML $captionMD $caption}}
    <figcaption>
      {{- if $captionHTML}}
        {{$captionHTML | safeHTML}}
      {{else if $captionMD}}
        {{$captionMD | markdownify}}
      {{else}}
        {{$caption}}
      {{end -}}
    </figcaption>
  {{end}}
</figure>
```

This shortcode will take a `src`, `alt` and `caption` attribute. Instead of `caption` you can also provide `captionHTML` or `captionMD` for raw HTML to Markdown respectively.

If the original image is large enough, it will produce 3 pairs of images:

- For small screens (< 576px). The format is WebP.
- For large screens (≥ 577px). The format is WebP.
- For fallback.

Each pair contains images of two different pixel densities. One for standard displays and the other for HiDPI displays.

Finally, the image will be wrapped in an anchor linking to the original image.

The resulting markup looks like this:

``` html
<figure>
  <a href="..." target="_blank">
    <picture>
      <source type="image/webp" media="(max-width: 575.98px)" srcset="... 1x, ... 2x">
      <source type="image/webp" media="(min-width: 576px)" srcset="... 1x, ... 2x">
      <img srcset="... 1x, ... 2x" src="..." alt="..." loading="lazy">
    </picture>
  </a>
  <figcaption>...</figcaption>
</figure>
```

## Wrap Up

The web has come a long way.

It might be difficult to keep up. But it is hard not to appreciate the thoughts that are being put into forwarding the platform.
