---
title: 'Go Tidbit: Update Checker with GitHub Releases'
date: 2023-06-24T12:11:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
---

After building [Printd](https://github.com/FurqanSoftware/toph-printd), Toph's print daemon, it became necessary to ensure that contest organizers were using the latest version of the software. Since Printd is open-source we host both the code and the release artifacts on GitHub.

The following function uses the Go client library for GitHub to check the latest release and compare the tag with the current version.

``` go
package main

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v53/github"
	"golang.org/x/mod/semver"
)

var (
	buildTag = "v0.3.0"

	repoOwner = "FurqanSoftware"
	repoName  = "toph-printd"
)

func checkUpdate(ctx context.Context) error {
	// Give your program at most 5 seconds to check for updates.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check the latest release on GitHub.
	client := github.NewClient(nil)
	rel, _, err := client.Repositories.GetLatestRelease(ctx, repoOwner, repoName)
	if err != nil || rel == nil || rel.TagName == nil {
		return err
	}

	// Check if the latest release is newer than the current version.
	if semver.Compare(*rel.TagName, buildTag) > 0 {
		log.Printf("Update available (%s)", *rel.TagName)
	}

	return nil
}
```

Since the API is accessible publicly you do not need to authenticate the GitHub API requests.

The current version is stored in the `buildTag` variable. You can easily set this variable at build time using `ldflags` as shown in [Go Tidbit: Setting Variables in Go During Build](/blog/go-tidbit-setting-variables-in-go-during-build/).

<br>

_This post is 22nd of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
