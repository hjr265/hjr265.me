---
title: 'Go Tidbit: How Much Was Read From io.Reader?'
date: 2023-02-05T12:50:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

You have an `io.Reader`, and you are about to pass it to some utility that will read from it. But the utility won't tell you how much it has read from the `io.Reader`.

How do you figure out how many bytes of data were in the `io.Reader`?

Like this snippet of Go+AWS code:

``` go
func uploadToS3(r io.Reader) (int64, error) {
	_, err := s3manager.NewUploader(awsSess).Upload(&s3manager.UploadInput{
		Bucket:             aws.String("some-bucket"),
		Key:                aws.String(b.Path),
		Body:               r,
	})
	// How many bytes were in r?
	return 0, err
}
```

You could implement a type that wraps the `io.Reader` and proxies all calls to `Read` while keeping track of the number of bytes read.

Or, you could use an `io.LimitedReader`:

``` go
func uploadToS3(r io.Reader) (int64, error) {
	lr := io.LimitedReader{R: r, N: math.MaxInt64} // Set the limit to the maximum possible value of N (int64).
	_, err := s3manager.NewUploader(awsSess).Upload(&s3manager.UploadInput{
		Bucket:             aws.String("some-bucket"),
		Key:                aws.String(b.Path),
		Body:               &lr,
	})
	n := math.MaxInt64 - lr.N // Subtract to see how many bytes have been read from r.
	return n, err
}
```

This method works while the number of bytes in the `io.Reader` is not greater than the maximum possible value of a 64-bit integer. You probably aren't hitting that limit. If you are, you have other challenges to deal with.
