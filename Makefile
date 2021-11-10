build:
	go build -o mockserver

run: build
	./mockserver

generate:
	go generate ./...

test:
	go test -v -race -timeout=10s ./...

integration-test:
	go test -v -timeout=10s ./test/...

docker-build:
	docker build -t mockserver .

docker-run: docker-build
	docker run -t mockserver
