---
title: 'Certbot-Name.com DNS Challenge Hooks'
date: 2023-09-01T17:00:00+06:00
tags:
  - Certbot
  - NameDotCom
  - 100DaysToOffload
  - DNSChallenge
---

With Certbot, for SSL certificates, I tend to use its DNS challenge over the HTTP challenge method. I find it a lot easier to automate. And, I can keep the certificate renewal details and mechanisms off the balancer or application servers.

But I also have domains that are on [Name.com](https://name.com). And, for some of them, I have their DNS records in the same place.

How do you then use Certbot's DNS challenge method with Name.com domains?

The cool thing about Certbot is that you can use the `--manual` flag with `--preferred-challenges dns` and provide hook scripts. You need just two: auth and cleanup.

The auth hook, when called, is expected to set up a TXT record. Let's Encrypt will check for this record to validate the request for the new certificate.

The cleanup hook, when called, is expected to clean up the said TXT record.

The Certbot command would look something like this:

``` sh {linenos=false}
NAME_DOMAIN=example.com \
NAME_USERNAME=... \
NAME_API_TOKEN=... \
certbot -v certonly --manual --preferred-challenges dns --manual-auth-hook name-auth.sh --manual-cleanup-hook name-cleanup.sh -d $NAME_DOMAIN
```

Make sure to set the `NAME_USERNAME` and `NAME_API_TOKEN` environment variables. You can get an API token from Name.com following the ["Signing up for API access" guide](https://www.name.com/support/articles/360007597874-signing-up-for-api-access).

Here are the two shell scripts that we will use for the two hooks:

**name-auth.sh**

``` bash
#!/bin/bash

HOST="_acme-challenge.$CERTBOT_DOMAIN"
HOST=${HOST%.$NAME_DOMAIN}

curl -s -u "$NAME_USERNAME:$NAME_API_TOKEN" "https://api.name.com/v4/domains/$NAME_DOMAIN/records" -X POST --data '{"host":"'$HOST'","type":"TXT","answer":"'$CERTBOT_VALIDATION'","ttl":300}'

echo Waiting 30 seconds for DNS changes to propagate
sleep 30
```

**name-cleanup.sh**

``` bash
#!/bin/bash

IDS=`curl -s -u "$NAME_USERNAME:$NAME_API_TOKEN" "https://api.name.com/v4/domains/$NAME_DOMAIN/records" | jq '.records[] | select(.host != null) | select(.host | startswith("_acme-challenge")) | .id'`
for id in $IDS; do
  curl -s -u "$NAME_USERNAME:$NAME_API_TOKEN" "https://api.name.com/v4/domains/$NAME_DOMAIN/records/$id" -X DELETE
done
```

You will need `curl` and `jq` as dependencies for these scripts.

Save these two scripts in the directory from where you will be running `certbot`. And, that's it.

Certbot can now be used to request for SSL certificates for domains on Name.com using the DNS challenge method.
