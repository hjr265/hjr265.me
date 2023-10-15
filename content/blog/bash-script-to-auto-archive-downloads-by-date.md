---
title: 'Bash Script to Auto-archive Downloads by Date'
date: 2023-10-11T10:00:00+06:00
tags:
  - Bash
  - Linux
  - 100DaysToOffload
---

Finding the files you are looking for without combing through hundreds of directories is a true time-saver and an easy productivity move.

I try to keep my files and directories in order, named and organized neatly. I don't have stale files at the base of my home directory. I have separate directories for my projects, my company stuff, and the work I do for my clients. 

But my `~/Downloads` directory is always a colossal mess.

Since I cannot find the time to go through it and clean stuff up, I wrote a script that sweeps the problem under the rug.

``` bash
#!/bin/bash

pushd ~/Downloads # Navigate to the Downloads directory in home.

for f in *; do
	if [ -d "$f" ]; then
		continue # Skip all directories.
	fi

	group=`date -r "${f@E}" '+%Y-%m'` # Get the file's modified year and month (YYYY-MM).
	if [ ! -d "archives/$group" ]; then
		mkdir -p "archives/$group" # Make an archive directory with year and month if it doesn't exist.
	fi
	mv "$f" "archives/$group/" # Move file to the archive directory.
done

popd # Return to the previous directory.
```

This script will take each file in your `~/Downloads` directory and move it to an archive directory named `YYYY-MM` where `YYYY` and `MM` are the year and month of the file's modified at timestamp (which is usually when the file was downloaded). It will create the `archive` directory and subdirectories as needed.
