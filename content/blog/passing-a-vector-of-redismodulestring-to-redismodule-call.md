---
title: Passing A Vector of RedisModuleString to RedisModule_Call
date: 2022-11-26T08:00:00+06:00
tags:
  - redis
  - module
  - 100DaysToOffload
---

While building [Redis Too](https://github.com/hjr265/redis-too) (a recommendation engine module for Redis), I spent a good bit of time (and strands of hair) figuring out how to pass a variable number of arguments to a Redis command like "SUNION".

The format specified of the [`RedisModule_Call`](https://redis.io/docs/reference/modules/modules-api-ref/#RedisModule_Call) function accepts 'v' to denote "a vector of RedisModuleString". But to use it, you need to pass two arguments to `RedisModule_Call`: the array of `RedisModuleString` and the size of the array.

Assuming you have the array `args` defined as `RedisModuleString **args`, and the size of the array `n` as `size_t n`, you can call `RedisModule_Call` as follows:

``` c
RedisModule_Call(ctx, "SUNION", "v", args, n);
```

You can also mix it with other arguments in the call:

``` c
RedisModule_Call(ctx, "SUNIONSTORE", "sv", key, args, n); // `key` is a RedisModuleString.
```

---

_This post is 1st of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
