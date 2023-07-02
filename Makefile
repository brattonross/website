build:
	make build-node
	make build-go

build-node:
	pnpm run build

build-go:
	go build -o bin/main

fmt: fmt-go fmt-html

fmt-go:
	@go fmt ./...

fmt-html:
	@pnpm run format
