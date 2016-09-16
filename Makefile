

all: api clients server extensions

api:
	go build ./

clients:
	go build ./client/native
	go build ./client/nc
	go build ./client/mock

server:
	go build ./server/natsgw

extensions:
	go build ./ext/audio
	go build ./ext/prompt
	go build ./ext/record

mock:
	go generate ./client/mock

build_test:
	docker build -t test-asterisk:13.8 ./internal/dockertest
