

all:
	go build ./
	go build ./client/native
	go build ./client/nats
	go build ./audio
	go build ./prompt

gw:
	go build ./cmd/nats-gw

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
