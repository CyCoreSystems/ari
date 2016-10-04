

all: api clients

api:
	go build ./
	go build ./stdbus

clients:
	go build ./client/native
	go build ./client/mock

extensions:
	go build ./ext/audio
	go build ./ext/prompt
	go build ./ext/record

mock:
	go generate ./client/mock

