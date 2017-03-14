SHELL = /bin/bash

EVENT_SPEC_FILE = internal/eventgen/json/events-14.0.0-rc1.json

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

events:
	go build -o bin/eventgen ./internal/eventgen/...
	@./bin/eventgen internal/eventgen/template.tmpl ${EVENT_SPEC_FILE} |goimports > events_gen.go
	
mock:
	go generate ./client/mock

