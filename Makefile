build:
	go build -o mockserver

run: build
	./mockserver

generate:
	go generate ./...

.PHONY: test
test:
	go test -v -race -timeout=10s $$(go list ./... | grep -v /component_tests/)

comp-test:
	go test -v -timeout=10s ./component_tests/...

docker-build:
	docker build -t mockserver .

docker-run: docker-build
	docker run -t mockserver
