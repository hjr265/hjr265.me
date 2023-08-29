---
title: Private and Secure Multi-cloud Network with WireGuard
date: 2023-01-17T13:05:00+06:00
tags:
  - Infrastructure
  - WireGuard
  - 100DaysToOffload
  - MultiCloud
  - PrivateNetwork
toc: true
---

Can you make a private and secure network between servers on multiple cloud providers?

Yes, with WireGuard, you can. The solution is fast and secure. It is easy to deploy/automate. And, best of all, it is cloud-agnostic.

For convenience, I have created a [demo repository](https://github.com/hjr265/clique) with Terraform and Ansible that you can clone and try. The demo uses [Linode](https://www.linode.com/lp/refer/?r=8d4f388136825d3d04a90d3f7b0ce6b29732a835) (referral) and [Vultr](https://www.vultr.com/?ref=8025299) (referral).

![](overview.png)

In this article, I will walk you through the inner workings of the demo and this solution.

## Provisioning The Infrastructure

In the demo, we are using Terraform to provision the infrastructure. 

We are making 3 servers on Linode:

``` terraform
resource "linode_instance" "node" {
  count = 3

  region = "ap-south"
  type   = "g6-nanode-1"
  image  = "linode/ubuntu20.04"

  booted     = true
  private_ip = true

  root_pass = var.root_password
}
```

And, 3 servers on Vultr:

``` terraform
resource "vultr_instance" "node" {
  count = 3

  plan   = "vc2-1c-1gb"
  region = "sgp"
  os_id  = 387
}
```

## Generating Ansible Inventory

While provisioning the servers, we also want to generate a hosts.yml (Ansible inventory) file. For this, we will provide Terraform with a template:

``` jinja
---
all:
  hosts:
%{ for i, host in hosts ~}
    ${host.name}:
      pipelining: true
      ansible_ssh_user: root
      ansible_ssh_pass: ${host.root_password}
      ansible_host: ${host.host}
      ansible_ssh_port: ${host.ssh_port}

%{ endfor ~}
wireguard:
  hosts:
%{ for i, host in hosts ~}
%{ if can(host.wireguard_ip) ~}
    ${host.name}:
      wireguard_ip: ${host.wireguard_ip}
      wireguard_addrs:
%{ if can(host.host_wan) ~}
        wan: ${host.host_wan}:51871
%{ endif ~}
%{ if can(host.host_lan) ~}
        lan: ${host.host_lan}:51871
%{ endif ~}
      wireguard_port: 51871
      wireguard_zone: ${host.zone}

%{ endif ~}
%{ endfor ~}
  vars:
    ansible_become_method: su

    wireguard_mask_bits: 16
    wireguard_port: 51871
```

And use the `local_file` resource to generate the file based on this template:

``` terraform
resource "local_file" "ansible_inventory" {
  content = templatefile("hosts.tpl", {
    hosts = concat(
      [for i, v in linode_instance.node : { name = v.label, root_password = var.root_password, host = v.ip_address, zone = "linode-ap-south", host_wan = v.ip_address, host_lan = v.private_ip_address, ssh_port = 22, wireguard_ip = "10.15.0.${i + 1}" }],
      [for i, v in vultr_instance.node : { name = v.label != null ? v.label : "vultr-node-${i + 1}", root_password = v.default_password, host = v.main_ip, zone = "vultr-sgp", host_wan = v.main_ip, ssh_port = 22, wireguard_ip = "10.15.1.${i + 1}" }],
    )
  })
  filename = "hosts.yml"
}
```

## Applying Terraform

After initializing Terraform, you will want to prepare three environment variables:

- `LINODE_TOKEN`: Get a Linode API Token with privileges to create new servers.
- `VULTR_API_KEY`: Get a Vultr API Key with privileges to create new servers.
- `TF_VAR_root_password`: Provide a root password for the new Linode servers. Vultr will generate random passwords automatically.

After running `terraform apply` you will notice 3 servers have been created in Linode and 3 in Vultr. At the same time, you will have a `hosts.yml` file generated in the `terraform/` directory. The file will look something like this:

``` yaml
---
all:
  hosts:
    [...]:
      pipelining: true
      ansible_ssh_user: root
      ansible_ssh_pass: [...]
      ansible_host: [...]
      ansible_ssh_port: 22

    # 5 more...

wireguard:
  hosts:
    [...]:
      wireguard_ip: [...]
      wireguard_addrs:
        wan: [...]:51871
        lan: [...]:51871
      wireguard_port: 51871
      wireguard_zone: [...]

    # 5 more...

  vars:
    ansible_become_method: su

    wireguard_mask_bits: 16
    wireguard_port: 51871
```

This hosts.yml file will contain details for the 6 servers created on the two cloud providers.

## Ansible Playbook for WireGuard

The Ansible Playbook included in the demo does the following for each server:

- Installs WireGuard using APT.
- Generates a pair of public and private keys.
- Generates pre-shared keys for each pair of servers.
- Sets up WireGaurd using Systemd Network on each server.

While setting up WireGuard, it will use the local area network of the cloud provider when possible. Otherwise, it will set up WireGuard to connect over the Internet. Regardless, with WireGuard, all communication through this WireGuard network between these 6 servers is private and secure.

## Networking Over WireGuard

In the end, each server can reach the other over `10.15.x.y` IP addresses.

The way this demo is set up, the 3 Linode servers will have WireGuard IPs `10.15.0.1`, `10.15.0.2` and `10.15.0.3`. The 3 Vultr servers will have WireGuard IPs `10.15.1.1`, `10.15.1.2` and `10.15.1.3`.

You can ping the Vultr servers from one of the Linode servers like so:

``` text
root@localhost:~# ping -c 3 10.15.1.1
PING 10.15.1.1 (10.15.1.1) 56(84) bytes of data.
64 bytes from 10.15.1.1: icmp_seq=1 ttl=64 time=2.19 ms
[...]

root@localhost:~# ping -c 3 10.15.1.2
PING 10.15.1.2 (10.15.1.2) 56(84) bytes of data.
64 bytes from 10.15.1.2: icmp_seq=1 ttl=64 time=1.47 ms
[...]

root@localhost:~# ping -c 3 10.15.1.3
PING 10.15.1.3 (10.15.1.3) 56(84) bytes of data.
64 bytes from 10.15.1.3: icmp_seq=1 ttl=64 time=2.47 ms
[...]
```

## Private and Secure

Please note that this demo and the article only outline the process of automating the setup process of WireGuard. Neither this demo nor the article is a comprehensive server and/or network hardening guide. We recommend you follow all the usual best practices and use this as a tutorial for only automating the WireGuard setup.

You would, at the very least, want to do the following in any setup that is not just for testing:

- Set up firewall
- Harden SSH access
- Configure all services to listen on and connect only over the WireGuard interface

With that said, tell me what you think and how would you improve this?
