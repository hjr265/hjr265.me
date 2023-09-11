.PHONY: clean
clean:
	rm -rf public

.PHONY: build
build:
	hugo

.PHONY: build.production
build.production:
	hugo --minify

.PHONY: deploy
deploy:
	hugo deploy
	s3cmd --access_key=${AWS_ACCESS_KEY_ID} --secret_key=${AWS_SECRET_ACCESS_KEY} setacl --acl-public --recursive s3://hjr265.me/

.PHONY: clean-build-deploy
clean-build-deploy: clean build.production deploy
