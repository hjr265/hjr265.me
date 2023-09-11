---
title: 'Schedule MongoDB Backups with GitLab CI/CD'
date: 2023-08-30T15:00:00+06:00
tags:
  - MongoDB
  - GitLab
  - 100DaysToOffload
  - Backups
---

GitLab CI/CD let's you schedule pipelines. And, in a way, I find it such a convenient way to manage my Internet crons.

One use of GitLab CI/CD schedules that I make is to backup MongoDB data. The pipeline, when run, SSHs into a server holding the MongoDB instance, runs `mongodump` and pipes it straight to `s3cmd` that then stores the dumped archived to an S3-esque bucket.

Here's what the pipeline looks like in `.gitlab-ci.yaml`:

``` yaml
snapshot:mongodb:
  image: DOCKER_IMAGE
  only:
    - schedules
    - web
  script:
    - TIMESTAMP=$(date +%Y%m%d%H%M%S)
    - mkdir -p ~/.ssh && echo "$PRODUCTION_HOST_KEYS" > ~/.ssh/known_hosts && echo "$PRODUCTION_SSH_PRIVATE_KEY" > ~/.ssh/id_rsa && chmod 0600 ~/.ssh/id_rsa
    - echo "$PRODUCTION_S3CFG" > ~/.s3cfg
    - ssh -p $SSH_PORT root@$SSH_HOST "mongodump --uri=mongodb://${MONGODB_USER}:${MONGODB_PASS}@127.0.0.1:${MONGODB_PORT}/$MONGODB_NAME --gzip --archive $MONGODUMP_ARGS" | s3cmd -v put - s3://${BUCKET_NAME}/${ARCHIVE_NAME_PREFIX}_${TIMESTAMP}.archive.gz
```

Yes, there are just four lines in that script.

The first line `TIMESTAMP=$(date +%Y%m%d%H%M%S)` sets the `TIMESTAMP` variable. We use this to name our dumped archive.

The second line prepares the CI/CD environment for the SSH connection. It expects the SSH host keys and private key in the `PRODUCTION_HOST_KEYS` and `PRODUCTION_SSH_PRIVATE_KEY` environment variables.

The third line configures `s3cmd`. You need to keep a valid `s3cmd` configuration file in the `PRODUCTION_S3CFG` environment variable. It should allow you to upload objects to the bucket where you want to keep your MongoDB dump archives.

The fourth line is a mouthful. But it's easy if you look at it in parts.

- `ssh -p $SSH_PORT root@$SSH_HOST`: This SSHs into the remote server and runs the command that follows.
- `"mongodump --uri=mongodb://..."`: This is the command that is run on the remote server. It dumps MongoDB into an archive. The archive is streamed to stdout, which in turn is streamed through the SSH connection and piped into the next command.
- `s3cmd -v put - s3://...`: This takes the dumped archive, piped to its stdin, and uploads it to an S3-esque object storage bucket.

The Docker image you need for this should have the following tools/packages:

- ca-certificates
- gnupg2
- openssh-client
- s3cmd

Something like this will do:

``` Dockerfile
FROM ubuntu:focal

RUN apt-get update && \
  apt-get install -y --no-install-recommends ca-certificates curl gnupg2 openssh-client s3cmd wget
```

And that's it. You now just need to set up a schedule for this pipeline from the GitLab web interface.

{{< image src="screen.png" alt="Screenshot of Edit Pipeline Schedule page from GitLab" caption="Screenshot of Edit Pipeline Schedule page from GitLab" >}}
