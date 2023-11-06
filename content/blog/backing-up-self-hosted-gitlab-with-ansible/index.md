---
title: 'Backing up Self-hosted GitLab With Ansible'
date: 2023-11-06T12:00:00+06:00
tags:
  - Ansible
  - GitLab
  - 100DaysToOffload
---

I have been making software for over a decade now. And one thing I have learned to love through this is automation.

After all, it is only a programmer who will spend hours automating a task that takes a few minutes to do.

There are good reasons for this.

I have a self-hosted GitLab instance for [Furqan Software](https://furqansoftware.com). And if you are self-hosting tools and services, the critical thing to do after ensuring security is automating backups.

Here is what an Ansible playbook looks like for backing up a self-hosted GitLab.

``` yaml
---
- hosts: gitlab
  gather_facts: yes
  tasks:
    - set_fact:
        backup_id: '{{ansible_date_time.iso8601_basic_short}}'

    - name: Backup GitLab
      shell: 'gitlab-backup create BACKUP={{backup_id}}'

    - name: Download backup
      fetch:
        src: '/var/opt/gitlab/backups/{{backup_id}}_gitlab_backup.tar'
        dest: './{{backup_id}}_gitlab_backup.tar'
        flat: true

    - name: Delete backup on remote
      file:
        state: absent
        path: '/var/opt/gitlab/backups/{{backup_id}}_gitlab_backup.tar'
```

Note that the `gitlab-backup create` command is sufficient to create the backup archive on the remote server. But this Ansible playbook backs up GitLab to a `.tar` file and downloads it to the local computer.
