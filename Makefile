.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test -coverprofile=./build/cover.out ./...

.PHONY: browse-cover
browse-cover:
	go tool cover -html=./build/cover.out

.PHONY: build
build:
	go build -o ./build ./cmd/twd
