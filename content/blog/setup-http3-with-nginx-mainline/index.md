---
title: 'Setup HTTP3 with NGINX Mainline'
date: 2023-10-08T10:00:00+06:00
tags:
  - NGINX
  - 100DaysToOffload
---

HTTP3 is here. Well, almost.

If you are using NGINX, you can update to the mainline version and start using HTTP3 today experimentally.

## Installing NGINX Mainline

As of writing this blog post, NGINX v1.24 is the latest stable version. But, HTTP3 is available in v1.25.

On Ubuntu/Debian-esque servers, the easiest way to install the mainline NGINX version is to [use NGINX's official repository](https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-open-source/#installing-a-prebuilt-debian-package-from-the-official-nginx-repository).

``` sh
apt install curl gnupg2 ca-certificates lsb-release debian-archive-keyring

# Import the official NGINX signing key.
curl https://nginx.org/keys/nginx_signing.key | gpg --dearmor \
    | tee /usr/share/keyrings/nginx-archive-keyring.gpg >/dev/null

# Add the APT repository.
echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] \
http://nginx.org/packages/mainline/debian `lsb_release -cs` nginx" \
    | tee /etc/apt/sources.list.d/nginx.list

# Install NGINX.
apt update
apt install nginx=1.25.*

nginx -v # Verify installation.
```

Alternatively, you can also [do this in Docker](https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-docker/).

## Enable HTTP3 for an NGINX Server Block (Virtual Host)

To enable HTTP3 for a `server` block, first configure SSL and verify the server block works.

Next, add a `listen` directive with the `quic` parameter and a response header advertising the support for HTTP3:

``` nginx
server ... {
  listen 443 quic reuseport; # For HTTP3
  listen 443 ssl;

  ssl_certificate     certs/example.com.crt;
  ssl_certificate_key certs/example.com.key;

  location / {
    add_header Alt-Svc 'h3=":443"; ma=86400'; # Response header to advertise HTTP3 support
  }
}
```

It is worth noting that if you have multiple server blocks where you want to enable HTTP3, you can use the `reuseport` parameter only once per `address:port` pair.

## Open 443/UDP

If you have the server behind a firewall or with _iptables_ rules blocking 443/UDP, reconfigure and open the port. HTTP3, unlike earlier versions of HTTP, is a UDP-based protocol.

_Yes, it took me a bit to realize my firewall was blocking 443/UDP._

## Restart NGINX

Validate your NGINX configuration with `nginx -t` and restart NGINX.

For some reason, reloading NGINX with `systemctl reload` didn't quite get HTTP3 working.

## Test HTTP3

Finally, test to make sure that HTTP3 is working.

You may use something like [http3check.net](https://http3check.net/) or `curl --http3` (if built with HTTP3 capabilities).
