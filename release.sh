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
