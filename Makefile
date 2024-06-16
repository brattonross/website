build:
	go generate
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/server

run: build
	./bin/server
