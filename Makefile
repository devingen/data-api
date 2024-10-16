REPO_NAME  := devingen/veri-api
IMAGE_TAG  := 0.1.0

.PHONY: build-docker
build-docker:
	@echo "Building Docker image"
	export GO111MODULE=on
	docker buildx build --platform linux/amd64 -t $(REPO_NAME):$(IMAGE_TAG) .

.PHONY: push-docker
push-docker:
	@echo "Pushing Docker image"
	docker push $(REPO_NAME):$(IMAGE_TAG)

.PHONY: release-docker
release-docker: build-docker push-docker
