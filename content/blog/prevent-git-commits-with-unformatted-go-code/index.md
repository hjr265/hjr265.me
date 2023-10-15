---
title: 'Prevent Git Commits with Unformatted Go Code'
date: 2023-10-07T10:10:00+06:00
tags:
  - Go
  - Tidbit
  - 100DaysToOffload
  - Git
---

Git has this great feature that I think is well-known but under-used. I am talking about [Git hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).

With Git hooks, you can run scripts during different Git actions.

Like this one:

``` sh
#!/bin/sh

GOFILES=`git diff --name-only --cached | grep -e '.go$' | grep -ve 'vendor/'`
UNFMTFILES=()
for f in $GOFILES; do
  if [ -n "`gofmt -l -s ./"$f"`" ]; then
    UNFMTFILES+=("$f")
  fi
done

if [ ${#UNFMTFILES[@]} -gt 0 ]; then
  echo You have staged unformatted Go files. Please run \`go fmt\` first.
  for f in ${UNFMTFILES[@]}; do
    echo " $f"
  done
  exit 1
fi
```

This script will take a list of all the staged Go files. It will then run `gofmt` to determine if these Go files are not formatted.

If you use it as a part of the pre-commit hook of your Git repository in your Go project, it will prevent commits with unformatted Go files.

To use it as the pre-commit hook, save it as `pre-commit` inside the `.git/hooks/` directory.
