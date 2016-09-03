

all:
	go build ./
	go build ./client/native
	go build ./audio
	go build ./prompt

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
