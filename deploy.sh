#!/bin/bash

set -e

hugo
hugo deploy
s3cmd setacl -rP s3://hjr265.me/
