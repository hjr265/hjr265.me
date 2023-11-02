---
title: "Setting Up Prometheus DNS-SRV Discovery with Terraform and DNSimple"
date: 2023-11-02T10:00:00+06:00
tags:
  - Prometheus
  - Terraform
  - 100DaysToOffload
  - DNSimple
---

If you are using Prometheus to collect metrics from your server, and you don't have a static set of servers, then you should set up automated discovery.

There are many ways you can set up automated discovery in Prometheus. However, one of my preferred vendor-agnostic ways of doing this is to use DNS-SRV records.

How does it work?

Let's say you are using Terraform to manage your infrastructure.

You will be creating a DNS-SRV record for each of your servers. As you create or destroy servers, these records will be created and destroyed.

Then, on the Prometheus end, you will configure it to look up these DNS-SRV records and determine the set of targets on the fly.

In Terraform, using [DNSimple](https://dnsimple.com/r/a21d2b08edfce9) (referral) only as an example:

``` tf
resource "dnsimple_zone_record" "SRV_node_prometheus" {
  for_each = toset([
    # Hostnames or IP addresses...
  ])

  zone_name = "example.com"
  name      = "_prometheus_node._tcp"
  value     = "1 9100 ${each.key}"
  type      = "SRV"
  ttl       = 300
  priority  = 10
}
```

The three fields in the value are weight, port, and target (i.e. hostname or IP address).

Once you have applied the Terraform changes, you can now tell Prometheus to use DNS-SRV discovery in the `prometheus.yml` file:

``` yaml
scrape_configs:
  - job_name: node
    dns_sd_configs:
      - names:
        - _prometheus_node._tcp.example.com
```

Once you have restarted Prometheus with this configuration, it will automatically detect new servers you add to the fleet or remove from it.
