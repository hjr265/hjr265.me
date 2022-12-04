---
title: "Single Node MongoDB Replica Set with Docker Compose"
date: 2022-12-04T15:00:00+06:00
tags:
  - mongodb
  - docker
  - 100DaysToOffload
  - dockercompose
toc: yes
---

Toph needs a few MongoDB features that only work in replica set mode. Take [Change Streams](https://www.mongodb.com/docs/manual/changeStreams/), for example. It is unavailable in standalone mode.

To ease development, all of Toph's application codebases come with Docker Compose files. With just a single `docker-compose up`, I can have any of Toph's applications running in a development environment. But having MongoDB start in replica set mode has been a bit of work.

You see, simply starting three nodes and then going in to configure them to be in a replica set wasn't something I wanted to do. I wanted only a single node to run, and I wanted it to start in replica mode, even if I was running `docker-compose up` on a freshly installed machine. 

So here is what I had to do:

## 1. Add MongoDB to `docker-compose.yml`

``` yml
services:
  # ...
  mongo:
    build:
      context: devel/mongodb
      dockerfile: Dockerfile
    hostname: mongo
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=secret
      - MONGO_INITDB_USERNAME=cat
      - MONGO_INITDB_PASSWORD=keyboard
      - MONGO_INITDB_DATABASE=databez
    volumes:
      - mongo:/data/db
    command: --keyFile /etc/keyfile --replSet rs0 --bind_ip_all

  # Sets up a single node replica set.
  mongoinit:
    image: mongo:6.0
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=secret
    volumes:
      - ./devel/mongodb/init.sh:/init.sh
    entrypoint: ["bash", "/init.sh"]
    depends_on:
      - mongo
```

## 2. Add `Dockerfile` to `devel/mongodb/`

``` dockerfile
FROM mongo:6.0

COPY --chown=999:999 initdb.d/10-createuser.sh /docker-entrypoint-initdb.d/10-createuser.sh
COPY --chown=999:999 keyfile /etc/keyfile

CMD ["--keyFile", "/etc/keyfile", "--replSet", "rs0", "--bind_ip_all"]
```

## 3. Add `10-createuser.sh` to `devel/mongodb/initdb.d/`

``` sh
#!/bin/bash

set -Eeuo pipefail
 
mongo -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase admin "$MONGO_INITDB_DATABASE" <<EOF
    db.createUser({
        user: '$MONGO_INITDB_USERNAME',
        pwd: '$MONGO_INITDB_PASSWORD',
        roles: [
            {
                role: 'readWrite',
                db: '$MONGO_INITDB_DATABASE'
            }
        ]
    })
EOF
```

## 4. Generate a `keyfile` for MongoDB Replica Set:

``` sh
openssl rand -base64 756 > ./devel/mongodb/keyfile
chmod 400 ./devel/mongodb/keyfile
```

## And, We Are Done!

You can now start the containers with Docker Compose:

``` sh
docker-compose up
```

<br>

_This post is 4th of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
