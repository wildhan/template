APP_NAME=template

run:
	go run ./app/main.go

build:
	go build -o bin/$(APP_NAME) ./app/main.go

test:
	go test ./...

clean:
	rm -rf bin

lint:
	golangci-lint run ./...

deps:
	go mod tidy

full-setup: clean deps build run
