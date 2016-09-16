

all: mock
	go build ./
	go build ./client/native
	go build ./client/nc
	go build ./server/natsgw
	go build ./audio
	go build ./prompt
	go build ./record

examples:
	mkdir -p bin/
	go build -o bin/helloworld ./_examples/helloworld/
	go build -o bin/stasisStart ./_examples/stasisStart

mock:
	go generate ./client/mock

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
