---
title: "Stream Uploading Files to S3 with Object Writer"
date: 2021-06-06T09:50:00+06:00
tags:
  - Go
  - AWS
  - S3
---

The official AWS SDK provides an upload manager construct that allows you to upload to S3 from any `io.Reader`. Using it is straightforward, that is until you need to create and upload a potentially large ZIP file.

The solution: use the upload manager with a pipe.

```golang
// Open a pipe.
pr, pw := io.Pipe()

errch := make(chan error)

// Upload from pr in a separate Go routine.
go func() {
	_, err := s3manager.NewUploader(awssess).Upload(&s3manager.UploadInput{
		Bucket:             aws.String("brain-bucket"),
		Key:                aws.String("blobs/very-large.zip"),
		Body:               pr,
	})
	errch <- err
}()

// Create a ZIP writer around pw.
zw := zip.NewWriter(pw)

// Add stuff to zip.

zw.Close() // Finishes the ZIP.

pw.Close() // Closes pw, marks EOF in pr.

err := <-errch // If err == nil, success.
```

For convenience, I have wrapped this up in a Go package: [github.com/hjr265/s3ow](https://pkg.go.dev/github.com/hjr265/s3ow) [ [GitHub](https://github.com/hjr265/s3ow) ].
