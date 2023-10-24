---
title: Serving hjr265.me From an S3-like Bucket Using a Caddy Module
date: 2023-10-24T15:00:00+06:00
tags:
  - Caddy
  - S3
  - 100DaysToOffload
---

I serve hjr265.me from an S3-like bucket hosted on Linode Object Storage. I have a Caddy instance that serves some of my Hugo-built websites, including this one.

I use Hugo's deployment function and `s3cmd` to deploy these websites.

## Why Both?

Hugo's deployment function uses the [blob package from the Go Cloud Development Kit](https://gocloud.dev/howto/blob/). This package comes with a limitation by design. It cannot set S3 access control lists (ACLs) for the uploaded objects. So neither can Hugo.

I have been using `s3cmd` to apply the "public-read" ACL for all the objects uploaded by Hugo.

I am using Caddy to serve the website using Linode Object Storage's website endpoint.

The catch? Unless the objects had a "public-read" ACL set, Caddy would serve 404.

And that meant, between the time `hugo deploy` was run and `s3cmd setacl` was run, updated parts of my site would return 404.

## Why Not Just `s3cmd`?

It doesn't work well for this use case unless you hold its hands.

Hugo can set different `content-type` and `cache-control` headers based on filename patterns. S3cmd cannot unless you upload the files in separate batches and add these details to the command line for each set of files.

Hugo can gzip the files before they are uploaded. S3cmd cannot.

## Solution: Use the `caddy.fs.s3` Module

If I can configure `caddy` to retrieve files from the S3-like bucket using an access key ID and secret pair instead of using Linode Object Storage's website endpoint, I can avoid some of these workarounds.

First, to add a module to `caddy`, I build it with `xcaddy`.

This Caddy module by Márk Sági-Kazár is what I needed: [github.com/sagikazarmark/caddy-fs-s3](https://github.com/sagikazarmark/caddy-fs-s3)

Here's a Dockerfile to build Caddy with the `caddy.fs.s3` module:

``` dockerfile
FROM caddy:2.7.4-builder-alpine AS builder

RUN xcaddy build \
    --with github.com/sagikazarmark/caddy-fs-s3

FROM caddy:2.7.4-alpine

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY caddy/Caddyfile /etc/caddy/Caddyfile
```

I now need a Caddyfile that could use this module with a `file_server`.

``` caddyfile
hjr265.me {
	encode gzip zstd

	file_server {
		fs s3 {
			bucket BUCKET_NAME
			region LINODE_REGION

			endpoint LINODE_REGION.linodeobjects.com
			profile PROFILE_NAME
		}
		index index.html
	}
}
```

Here:

- I am gzipping the responses. I had to stop making Hugo upload gzipped objects since Caddy doesn't set the `content-encoding` header based on the S3 response.
- I use the `file_server` directive with the `s3` module.
- The module is configured to use a specific credential profile from the `~/.aws/credentials` file.
- The `file_server` directive uses "index.html" for the index file.

## Serving a 404 Page

The above Caddyfile almost works the way I wanted. Except it wasn't serving the site's 404 page (/404.html) when accessing a non-existent URL.

To solve that, I need to add a `handle_errors` directive.

``` caddyfile
(s3-file-server) {
	file_server {
		fs s3 {
			bucket BUCKET_NAME
			region LINODE_REGION

			endpoint LINODE_REGION.linodeobjects.com
			profile BUCKET_NAME
		}
		index index.html
	}
}

hjr265.me {
	encode gzip zstd

	import s3-file-server

	handle_errors {
		rewrite * /{err.status_code}.html
		import s3-file-server
	}
}
```

Here:

- I moved the `file_server` directive to a snippet of its own. I can use it again inside `handle_errors` without duplicating the `file_server` directive.
- The `handle_errors` directive is set to serve the "/404.html" page for 404 responses (or the corresponding error page for any non-200 response).

## Wrap Up

Now, I can say goodbye to `s3cmd` from my Hugo-built site deployment pipelines and not worry about setting a public-read ACL on uploaded objects.

The pipelines now have to do half as much work as before.
