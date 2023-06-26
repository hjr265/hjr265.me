---
title: 'Makeshift CDN for Linode Object Storage with NGINX for Discourse'
date: 2023-06-26T9:15:00+06:00
tags:
  - NGINX
  - Discourse
  - 100DaysToOffload
  - Linode
---

If you use S3-like object storage for Discourse uploads, you will quickly realize that Discourse strongly recommends a CDN.

> Not using a CDN (or entering the bucket URL as the CDN URL) is likely to cause problems and is not supported. [\[...\]](https://meta.discourse.org/t/configure-an-s3-compatible-object-storage-provider-for-uploads/148916)

If you were to use Amazon S3, their Cloudfront service would be the easy solution. But not all cloud providers (e.g. [Linode](https://www.linode.com/lp/refer/?r=8d4f388136825d3d04a90d3f7b0ce6b29732a835) _referral_) providing S3-like services has CDN offerings. For such cases [NGINX with reverse proxy caching](https://www.nginx.com/resources/wiki/start/topics/examples/reverseproxycachingexample/) can work well as a makeshift CDN.

Since objects uploaded by Discourse to S3 buckets are marked to be publicly available, you can use NGINX to reverse proxy the bucket.

In your NGINX configuration file, [define a proxy cache](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_cache_path) followed by a server using it and reverse proxying to the bucket:

``` nginx
# Define a location for reverse proxy cache, and allocate 16 MB for keys and 8 GB for content. NGINX will remove anything that is older than ten weeks from the cache.
proxy_cache_path /tmp/nginx-discourse_uploads levels=1:2 keys_zone=discourse_uploads:16m max_size=8g inactive=10w use_temp_path=off;

server {
  listen 443 ssl http2;
  listen [::]:443 ssl http2;
  server_name uploads.discourse.example.com;

  # SSL and other configurations go here.

  location / {
    proxy_pass https://example-discourse-uploads.ap-south-1.linodeobjects.com$1; # Bucket URL goes here. The "$1" at the end is necessary.
    proxy_cache discourse_uploads;
    
    limit_except GET HEAD { # Deny all non-GET or HEAD requests.
      deny all;
    }

    proxy_hide_header x-amz-meta-s3cmd-attrs; # Remove response headers specific to S3-like storages.
    proxy_hide_header x-amz-request-id;
    proxy_hide_header x-amz-storage-class;
    proxy_hide_header x-rgw-object-type;

    proxy_cache_methods GET HEAD; # Cache GET and HEAD requests only.
    proxy_cache_valid 200 28d;    # Cache 200 responses for 28 days.
    proxy_cache_valid 403 404 1h; # Cache 403 and 404 responses for a shorter duration.

    expires max;
  }
}
```

<br>

_This post is 24th of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
