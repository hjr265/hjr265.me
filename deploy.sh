#!/bin/bash

# npm run build && \
(cd public/ && s3cmd sync --delete-removed --acl-public --exclude sitemap.xml . s3://hjr265.me/)
