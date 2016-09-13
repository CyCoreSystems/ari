

all: mock
	go build ./
	go build ./client/native
	go build ./client/nc
	go build ./server/natsgw
	go build ./audio
	go build ./prompt
	go build ./record

mock:
	go generate ./client/mock

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
