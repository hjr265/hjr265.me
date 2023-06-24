---
title: Generate SSH Known Hosts in Terraform
date: 2023-03-24T10:30:00+06:00
tags:
  - Terraform
  - 100DaysToOffload
---

I manage [Toph's](https://toph.co) infrastructure with Terraform. The setup is as multi-cloud as it gets. Beyond Terraform, I have several Ansible Playbooks for configuring and upkeeping the infrastructure. 

Something that I wanted to do with Terraform was to generate an SSH known_hosts file ahead of time as any changes are applied to the infrastructure. There isn't a straightforward way to do that in Terraform. But the workaround is simple too.

``` tf
locals {
  ssh_hosts = merge(
    [for i, v in vultr_instance.servers : { addr = v.main_ip, port = 22 }],
    [for i, v in linode_instance.servers : { addr = v.ip_address, port = 22 }],
    # ...
  )
}

resource "null_resource" "known_hosts" {
  provisioner "local-exec" {
    command = <<EOT
  rm -f known_hosts;
%{ for i, host in local.ssh_hosts }
  ssh-keyscan -p ${host.port} ${host.addr} >> known_hosts;
%{ endfor ~}
EOT
    interpreter = ["/bin/bash", "-c"]
  }
}
```

All you need to do is put all your SSH hostnames/addresses and ports in a local variable. And, then use a `local-exec` provisioner in a `null_resource` to run `ssh-keyscan` for each host, appending the output to a "known_hosts" file.

Note that our local `ssh_hosts` variable is a list of objects. Each object has the key `addr` and `port`. We use these values as a part of our `ssh-keyscan` command:

``` txt {linenos=false}
ssh-keyscan -p ${host.port} ${host.addr} >> known_hosts;
```

We loop over all the hosts and run this command once for each host.

Terraform uses `/bin/sh` as the interpreter on Linux. But in this case, we need to use Bash or equivalent to be able to run this script.

<br>

_This post is 21st of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
