

all:
	go build ./
	go build ./client/native
	go build ./audio
	go build ./prompt
	go build ./record

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
