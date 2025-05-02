.PHONY: all build run clean test

all: build

build:
	go build -o go_mcp_server_youtube_search main.go

run:
	go run main.go

clean:
	rm -f go_mcp_server_youtube_search