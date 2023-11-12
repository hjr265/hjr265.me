---
title: "Multi-threaded Downloads in Go"
date: 2023-11-09T14:00:00+06:00
tags:
  - Go
  - 100DaysToOffload
---

The word multi-threaded here is an artifact of how download managers in the past worked.

The idea is to download a large file in parts, in parallel, over multiple TCP streams at once. In certain circumstances this can speed up the download significantly.

Let's start with a naive way of downloading a file in Go:

``` go
// Error handling omitted for brevity.

// Perform a GET request.
resp, _ := http.Get(url)
defer resp.Body.Close()

// Create the output file.
f, _ := os.Create("output.ext")

// Copy from the response body to the file.
io.Copy(f, resp.Body)
f.Close()
```

The code above is downloading the entire file in a single stream.

However, you need more code to download a file in multiple streams parallely.

Let us start by defining a type and a few consts:

``` go
// A chunk represents the part of the file and holds the relevant HTTP response.
type chunk struct {
  resp  *http.Response
  start int64
  end   int64
}

const (
  nthreads = 3           // Number of download threads
  bufsize  = 5*1024*1024
  readsize = 1024*1024
)
```

First, we will perform a GET request like before:

``` go
resp, err := http.Get(url)
catch(err)
defer resp.Body.Close()
```

Next, we will check if the server really supports multi-threaded downloads. For this we need two things:

- We need to know the total size of the download. We need the `Content-Length` header in the response.
- The `Accept-Ranges` header with the value "bytes".

If the server or the response do not meet these two conditions, we continue the download in a single stream.

Otherwise, we continue and plan the chunks, the parts of the file that we will download in separate streams:

``` go
// Parse the size and determine the chunk size.
sz, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
chsz := (sz + int64(nthreads-1)) / int64(nthreads)

chunks := []chunk{
  // Use the response from the first request in the first chunk.
  {
    resp: resp,
    end:  chsz,
  },
}

// Plan the remaining threads.
for i := 1; i < nthreads; i++ {
  // Prepare a request.
  req, _ := http.NewRequest("GET", url, nil)

  // Request download from an offset. Set the Range header accordingly.
  start := chsz * int64(i)
  end := min(start+chsz, sz)
  req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

  resp, _ := http.DefaultClient.Do(req)
  defer resp.Body.Close()

  chunks = append(chunks, chunk{
    resp:  resp,
    start: start,
    end:   end,
  })
}
```

Next, create the output file and spawn a goroutine for each chunk:

``` go
// Create the output file.
f, _ := os.Create("output.ext")
m := sync.Mutex{} // Mutex to synchronize writes to the file.

wg := sync.WaitGroup{}
wg.Add(len(chunks))
for _, ch := range chunks {
  go func(ch chunk) {
    defer wg.Done()
    defer ch.resp.Body.Close()

    // Prepare a buffered limited reader.
    r := io.LimitReader(ch.resp.Body, ch.end-ch.start)
    br := bufio.NewReaderSize(r, 5*1024*1024)

    buf := make([]byte, 1024*1024)
    pos := ch.start
    stop := false
    for !stop {
      // Read from buffer reader until EOF.
      n, err := br.Read(buf)
      if err == io.EOF {
        stop = true
        err = nil
      }

      m.Lock()
      f.Seek(pos, 0) // Seek the file to the appropriate position before writing data.
      f.Write(buf[:n])
      m.Unlock()
      pos += int64(n)
    }
  }(ch)
}
wg.Wait()

f.Close()
```

As multiple goroutines are writing to the same file and at different positions within the file, we use a `sync.Mutex` to synchronize access to it.

Once all the goroutines end, the file is closed, and the download is complete.

And that's the core idea to multi-threaded downloads in Go. But you have to implement more than this to reliability download a file, including proper retries and error handling.

You can also take this further and better utilize the bandwidth by dynamically splitting in-progress chunks as the goroutines finish downloading.
