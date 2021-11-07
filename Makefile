build:
	go build -o mockserver

run: build
	./mockserver

generate:
	go generate ./...

test:
	go test -race -timeout=10s ./...
