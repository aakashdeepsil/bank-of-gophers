build:
	@go build -o bin/bog

run: build
	@./bin/bog

test:
	@go test -v ./...