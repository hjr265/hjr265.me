---
title: "Resurrecting a 12-Year-Old Node.js Project With Claude Code"
date: 2026-03-07T17:30:00+06:00
tags:
  - NodeJS
  - Docker
  - ClaudeCode
toc: yes
---

This weekend, I set out to get some high-resolution screenshots of zebra-algo (you may remember it as the now-defunct [CodeMarshal](/projects/#codemarshal) or algo.codemarshal.org). My motivation was to document its design and features, as zebra-algo is a competitive programming contest platform I built back in 2014.

The plan was straightforward: run it locally, open it in a browser, take the screenshots, and be done.

Simple, except the codebase is twelve years old, making things tricky.

Before this, I revived an old Go project. Even though Go is known for stability, outdated dependencies like MongoDB and Elasticsearch made the process slow and tricky.

This time, I tried something different by bringing Claude Code into the process, running it inside [Zed](https://zed.dev/). What followed was half a working day of surprisingly enjoyable archaeology.

## What zebra-algo Actually Is

The project is a contest platform with problem preparation, code submissions, automated judging, real-time standings, and clarifications. I wanted to collect screenshots of the platform to document on my website.

{{< image src="contest-dashboard.png" alt="Contest dashboard" caption="Contest dashboard" >}}

The platform's stack is [IcedCoffeeScript](https://github.com/maxtaco/coffee-script) for async logic, Express `4.0.0-rc3`, Mongoose `3.8.8` with MongoDB 2.6, Kue `0.6.2` for Redis queues, ZeroRPC using [ZeroMQ](https://zeromq.org/), Knox as S3 client, and Socket.io `0.9.16`. Client-side dependencies are managed with Bower.

The version I have isn’t even the latest release. It’s from when I was personally paying for the project, so it’s a bit older and rougher around the edges.

## The Compatibility Gauntlet

Getting this running in 2026 meant working through a cascade of problems, each one a small lesson in how much the [Node.js](https://nodejs.org/) ecosystem has shifted over a decade.

**Picking the right Node.js version was already a puzzle.** The build scripts for `zeromq` use ES6 const, which Node 4 rejects outside strict mode in certain transitive dependencies. But `bignum@0.6.2` (pulled in by `mongoose-shortid`) uses `nan@1.x`, which doesn’t compile against Node 4’s V8 API. The answer was Node 6: new enough for `zeromq`, workable for `bignum` if you swap in a current `bignum` with `nan@2.x` and build it manually, then overwrite the nested copy.

**Debian Jessie, the base OS for the `node:6` Docker image, has been archived.** Its apt `Valid-Until` dates have expired, and its GPG keys are no longer trusted by default. You have to redirect `sources.list` to `archive.debian.org` and pass `-o Acquire::Check-Valid-Until=false` to get any packages installed at all.

**Every installable `node-gyp` is too new.** Any version you can `npm install -g` today uses async/await internally (requires Node 7.6+) or has transitive deps using spread syntax. The workaround: use the `node-gyp` bundled inside `npm@3` itself, which is already compatible with Node 6 by definition.

**`npm@3`’s flat dependency tree broke old hardcoded paths.** `mongoose@3.8.8` requires `mongodb/node_modules/bson` (a path that only exists under `npm@2`’s deeply nested layout). `npm@3` flattened everything, so that path no longer exists. The fix is to explicitly recreate the nested structure in the Dockerfile after installation.

**`kue@0.6.2` ignores Redis configuration.** The `createQueue({redis: {host, port}})` API wasn’t added until `kue@0.7`. Passing options in `0.6.2` does nothing; it always connects to `127.0.0.1:6379`. The fix: monkey-patch `kue.redis.createClient` before any queue is created.

**Knox constructs S3 URLs for virtual hosts by default.** That means `bucket.endpoint`, which is fine for AWS but breaks against a local MinIO instance. Had to add `S3_PATH_STYLE`, `S3_ENDPOINT`, and `S3_PUBLIC_ENDPOINT` environment variables and thread them through the S3 abstraction layer.

There were a handful of others, too: `bson` shipping a prebuilt binary for the wrong platform, `connect-redis@1.x` not supporting a `url` option, `zeromq` moving from a nested path to a flat one in `npm@3`’s layout.

## The Dockerfile That Makes It All Work

After working through all this with Claude Code, the Dockerfile almost reads like a log of the problems I faced. Each RUN step is a workaround for something that stopped working during the project reviving process.

``` dockerfile
FROM node:6

# node:6 (Jessie) already has gcc 4.9, build-essential, and python.
# Just add the native library deps.
RUN echo "deb http://archive.debian.org/debian jessie main contrib" > /etc/apt/sources.list && \
    apt-get -o Acquire::Check-Valid-Until=false update && \
    apt-get install -y --no-install-recommends --force-yes \
        libgmp-dev \
        libzmq3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY package.json bower.json ./

# Skip native build scripts:
#   - zeromq@4.x prebuild-install tries to download a prebuilt and falls back to
#     preinstall.js + node-gyp; we drive that manually below.
#   - bignum@0.6.2 uses nan@1.x which does not compile against Node 6's V8.
RUN npm install --ignore-scripts && \
    npm install -g bower && \
    bower install --allow-root

# Use npm's own bundled node-gyp rather than a globally-installed one.
# Any node-gyp fetchable from the registry today has transitive deps that use
# async/await which requires Node 7.6+. npm@3 ships its own compatible version.
ENV NODE_GYP="node /usr/local/lib/node_modules/npm/node_modules/node-gyp/bin/node-gyp.js"

# bignum@0.6.2 (inside mongoose-shortid/node_modules) uses nan@1.x, incompatible
# with Node 6's V8. Install current bignum (nan@2.x) with --ignore-scripts so npm
# doesn't also trigger the nested bignum@0.6.2 build, then build it manually, and
# drop it in place of the old version.
RUN npm install bignum --ignore-scripts && \
    cd /app/node_modules/bignum && $NODE_GYP rebuild && \
    rm -rf /app/node_modules/mongoose-shortid/node_modules/bignum && \
    cp -r /app/node_modules/bignum /app/node_modules/mongoose-shortid/node_modules/bignum

# mongoose@3.8.8 hard-codes require('mongodb/node_modules/bson'); npm@3 flattened
# bson to the top level. The npm package includes a prebuilt bson.node for a
# different platform; rebuild it from source, then copy to where mongoose expects it.
RUN cd /app/node_modules/bson && ($NODE_GYP rebuild || true) && \
    mkdir -p /app/node_modules/mongodb/node_modules && \
    cp -r /app/node_modules/bson /app/node_modules/mongodb/node_modules/bson

# Build zeromq via its install script.
# npm@3 flattens deps so zeromq is at the top level, not nested under zerorpc.
RUN cd /app/node_modules/zeromq && \
    node scripts/preinstall.js && $NODE_GYP rebuild

RUN cd /app/node_modules/hiredis && \
    $NODE_GYP rebuild || true

COPY . .

EXPOSE 5000

CMD ["./node_modules/.bin/iced", "web.iced"]
```

The docker-compose.yml wires up era-appropriate versions of each service: MongoDB 2.6, Redis 2.8, and MinIO standing in for S3, along with the environment variables needed to make everything talk to each other.

``` yml
version: "2"

services:
  web:
    build: .
    ports:
      - "5000:5000"
    environment:
      PORT: 5000
      MONGO_URL: mongodb://mongo/zebra-algo
      REDIS_URL: redis://redis:6379
      S3_KEY: minioadmin
      S3_SECRET: minioadmin
      S3_BUCKET: zebra-algo
      S3_ENDPOINT: http://minio:9000
      S3_PUBLIC_ENDPOINT: http://localhost:9000
      S3_PATH_STYLE: "true"
    depends_on:
      - mongo
      - redis
      - minio

  mongo:
    image: mongo:2.6

  redis:
    image: redis:2.8

  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
```

`S3_PATH_STYLE`, `S3_ENDPOINT`, and `S3_PUBLIC_ENDPOINT` are all new — they didn’t exist in the original codebase. Claude Code added them to the S3 abstraction layer so Knox would route requests correctly to MinIO rather than trying to resolve `zebra-algo.minio` as a hostname.

## What Made Today Different

Here’s the thing: the problems above aren’t surprising alone. Anyone who has done this kind of archaeology knows these failures. Old apt repos expire. Native modules break across Node versions. Flat versus nested dependency trees cause path assumptions to crumble. You work through them one by one.

What was different today was the process.

With the Go project a few weeks ago, I was alone with the error messages. Each failure meant reading, thinking, searching, and trying. That’s the job. But it’s quiet, solitary work, and the cognitive load of holding the whole broken system in your head while working through it is real.

With Claude Code running in Zed, it felt more like pair programming. I would paste an error or describe what I was seeing and get back not just a fix but an explanation of why it was failing. It didn’t always get things right on the first try. But the back-and-forth had a different texture. The failures felt more like puzzles we worked through together than walls I climbed alone.

The project is genuinely hard to revive. IcedCoffeeScript is a language most people today have never encountered. The combination of legacy native modules, an archived OS base image, `npm@2`-era path assumptions, and undocumented behavior changes across library versions, it’s a lot. And yet the session moved forward steadily, with Claude Code identifying version-compatibility issues, proposing a monkey patch for kue, and working out the MinIO/Knox path-style problem without me having to spell out every detail.

## The Part That Surprised Me

After all that infrastructure work, the application ran correctly. The business logic in the IcedCoffeeScript file, the contest management, submission pipeline, and scoring, is intact and readable. The architecture for judging, where ZeroRPC calls a separate Saber service over ZeroMQ, is a sound design that holds up.

The code also has a certain historical charm. It predates async/await in JavaScript by three years, yet it reads cleanly because IcedCoffeeScript’s await/defer constructs solved the same problem in their own way. Looking at a working, running system that made those choices, choices that were genuinely forward-thinking at the time, is a bit like reading old correspondence from someone who turned out to be right.

I got my screenshots. The platform looks exactly as I remembered it.

{{< imagegrid layout="three" >}}
{{< image src="problem-view.png" alt="Problem view" caption="Problem view" >}}
{{< image src="problem-edit-outline.png" alt="Problem outline edit form" caption="Outline edit form" >}}
{{< image src="problem-edit-parameters.png" alt="Problem parameters edit form" caption="Parameters edit form" >}}
{{< /imagegrid >}}

## What’s Next

I have a few more old projects sitting in a similar state, things I built years ago, and didn’t revisit them for one reason or another for quite a while. I have been meaning to document them properly with screenshots and write-ups before they become genuinely unrecoverable.

After today, I think I will let Claude Code do the heavy lifting on those, too.
