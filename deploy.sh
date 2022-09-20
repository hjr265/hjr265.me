#!/bin/bash

set -e

hugo
hugo deploy
s3cmd --access_key=$AWS_ACCESS_KEY_ID --secret_key=$AWS_SECRET_ACCESS_KEY setacl -rP s3://hjr265.me/
