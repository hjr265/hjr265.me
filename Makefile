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

.PHONY: clean-build-deploy
clean-build-deploy: clean build.production deploy
