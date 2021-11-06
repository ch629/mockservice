build:
	go build -o mockserver

run: build
	./mockserver
