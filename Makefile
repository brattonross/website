run: build
	./bin/server

build: build-client build-server

build-client: build-js build-css

build-js:
	bun build main.js --outdir public

build-css:
	bun tailwindcss -i styles.css -o public/styles.css

build-server:
	go build -o ./bin/server ./main.go
