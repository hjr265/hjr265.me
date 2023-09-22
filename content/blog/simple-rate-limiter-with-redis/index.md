---
title: Simple Fixed-window Rate Limiter With Redis
date: 2023-07-01T19:50:00+06:00
tags:
  - Redis
  - Go
  - 100DaysToOffload
---

A while ago I needed a very quick rate limiter implementation. The application I was working on was already using Redis.

> **Fixed-window rate limiting:** This is a straightforward algorithm that counts the number of requests received within a fixed time window, such as one minute. Once the maximum number of requests is reached, additional requests are rejected until the next window begins. [...](https://redis.com/glossary/rate-limiting/)

With a small Redis script, I was able to implement a fixed-window rate limiter:

``` text
local current
current = redis.call("INCR", KEYS[1])
if tonumber(current) == 1 then
	redis.call("EXPIRE", KEYS[1], 60)
end
return current
```

Every time this script is run it takes a key and increments its value by 1. Whenever the key is incremented for the first time, an expiry of 60 seconds is set. It returns the current value after the increment.

The key expires 60 seconds after it is first set. Once expired, it will be set again on the next request.

Any code using this script can call it whenever a request for a rate-limitable action is received. If the value returned by the script is greater than allowed, the request is aborted due to the rate limit. If the value returned is not greater than allowed then the request is processed.

``` go
const script = `
local current
current = redis.call("INCR", KEYS[1])
if tonumber(current) == 1 then
	redis.call("EXPIRE", KEYS[1], 60)
end
return current
`

func isRateLimited(ctx context.Context, key string, limit int64) (bool, error) {	
	v, err := redisClient.Eval(ctx, script, []string{key}).Result()
	if err != nil {
		return false, err
	}
	n, _ := v.(int64)
	return n > int64(limit), nil
}
```

The function can then be used as follows:

``` go
func handleLogin(r *http.Request, w http.ResponseWriter) {
	username := r.FormValue("username")

	limited, _ := isRateLimited(context.TODO(), fmt.Sprintf("rateLimit:login:username:%s", username), 5)
	if limited {
		http.Error(w, "Too Many Attempts", http.StatusTooManyRequests)
		return
	}

	// ...
}
```

And this just works.

Note that a fixed-window rate limiter, although effective against sustained attacks, may affect the experience of legitimate users.

In the example above, we rate limit based on the username used during a login flow. This is less likely to affect legitimate users than using, for example, the remote IP address of the incoming request.
