.PHONY: run build clean docs

run:
	go run cmd/server/main.go

build:
	go build -o server cmd/server/main.go

clean:
	rm -f server

docs:
	swag init -g api/index.go --dir ./
