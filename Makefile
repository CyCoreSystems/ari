

all: api clients gateway

api:
	go build ./

clients:
	go build ./client/native
	go build ./client/nc
	go build ./client/mock

extensions:
	go build ./ext/audio
	go build ./ext/prompt
	go build ./ext/record

gateway:
	go build ./server/natsgw
	go build -o bin/ari-nats-gateway ./cmd/ari-natsgw

mock:
	go generate ./client/mock

