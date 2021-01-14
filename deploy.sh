#!/bin/bash

set -e

hugo deploy -v
s3cmd setacl -rP s3://hjr265.me/
