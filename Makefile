.PHONY: build-all-platforms
build-all-platforms: 
	docker buildx build --platform linux/amd64 . -f Dockerfile-hello-example

.PHONY: pre-commit
pre-commit: lint test
	go test ./...

.PHONY: test
test:
	go clean -testcache
	go test ./... -race -covermode=atomic -coverprofile=coverage.out

.PHONY: lint
lint:
	golangci-lint run ./...
