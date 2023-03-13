---
title: 'Randomised Suffix in Terraform Resource Names'
date: 2023-03-13T20:30:00+06:00
tags:
  - Terraform
  - 100DaysToOffload
---

I didn't like how I named cloud resources for [Toph](https://toph.co). The resource names had sequential suffixes at the end.

Let's take Toph's storage servers, for example, where MongoDB replicas are running. On Linode, while they were there, I had named the three storage servers as `toph-storage-sgp-1`, `toph-storage-sgp-2` and `toph-storage-sgp-3`.

But as I moved them over to [Vultr](https://www.vultr.com/?ref=8025299) last week, I realised that eliminating these sequential numbers could remove any notion of preference/precedence among these three servers.

It was time to adopt a better naming scheme. For that, I chose the following:

``` txt {linenos=false}
{application}-{role}-{cloud}-{region}-{suffix}
```

I wanted `{suffix}` to be a 5-character random string with only numbers and lower-cased alphabets.

Terraform makes creating randomised strings easy. The "random_string" resource type is just for that. The strings, although random, stay unchanged between each `terraform apply`.

Here is how you can add randomised suffixes to resource names in Terraform:

``` tf
resource "random_string" "vultr_instance_worker" {
  special = false
  upper = false
  length = 5
}

resource "vultr_instance" "worker" {
  label = "toph-storage-vultr-sgp-${random_string.vultr_instance_worker.result}"
  // ...
}
```

With this example, after running `terraform apply`, the resource label will be something like "toph-storage-vultr-sgp-ns2oh".

_P.S: Yes, I know I should also add the environment (e.g. "production", "staging", etc.) into that naming scheme._ 🙃

<br>

_This post is 20th of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
