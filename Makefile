

all:
	go build ./
	go build ./client/native

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
