

all:
	go build -v ./...

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
