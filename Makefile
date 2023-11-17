build:
	@go build -o bin/p2p-pocker

run: build
	@./bin/p2p-pocker

test:
	go test -v ./...