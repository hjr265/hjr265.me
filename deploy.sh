#!/bin/bash

set -e

rm -rf public
hugo --minify
hugo deploy --maxDeletes=5
s3cmd --access_key=$AWS_ACCESS_KEY_ID --secret_key=$AWS_SECRET_ACCESS_KEY setacl -rP s3://hjr265.me/
