---
title: Bump Calendar Versioning (CalVer) Bash Script for Git
date: 2022-12-05T09:00:00+06:00
tags:
  - CalVer
  - Git
  - 100DaysToOffload
  - Bash
  - Script
---

For Toph, SemVer didn't make much sense. [CalVer](https://calver.org/) did.

Toph is closed-source and available to users as a SaaS. We are not trying to maintain or communicate backwards-compatible/incompatible changes to Toph with our users. With CalVer (short for Calender Versioning), the versioning would be tied to the release date.

The pattern we wanted to use: `{YYYY}.{0M}.{SEQ}`.

- `{YYYY}`: Four-digit year.
- `{0M}`: Zero-padded month.
- `{SEQ}`: Starts from 0. Goes up by 1 for each bump. Resets after the end of the month.

For example: `2022.12.3`.

I needed a simple Bash script to bump Git tags with CalVer versions.

We use this for Toph and other related projects that use CalVer.

``` sh
#!/bin/bash

set -e

echo 'Pulling current tags.'
git pull --ff-only --tags

LATEST=`git describe --tags --abbrev=0`
echo 'Latest tag:' $LATEST

NEWPRE="v$(date +%Y).$(date +%m)"
NEWTAG=''
if [[ "$LATEST" = $NEWPRE* ]]
then
	NEWTAG="${LATEST%.*}.$((${LATEST##*.}+1))"
else
	NEWTAG="${NEWPRE}.0"
fi
echo 'Next tag:' $NEWTAG

read -p 'Create a new release? (Press "y" to confirm.) ' -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
	git tag $NEWTAG
	git push
	git push --tags
fi
```

Looks something like this when run:

``` txt
Toph/platform [master] » ./release.sh 
Pulling current tags.
Already up to date.
Latest tag: v2022.12.1
Next tag: v2022.12.2
Create a new release? (Press "y" to confirm.) 
```

On `Y`, a new Git tag is created and pushed to `origin`. The CI/CD takes care of the rest. Perhaps a story for another day.

<br>

_This post is 5th of my [#100DaysToOffload](/tags/100daystooffload/) challenge. Want to get involved? Find out more at [100daystooffload.com](https://100daystooffload.com/)._
