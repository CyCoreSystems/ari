SHELL = /usr/bin/env bash

EVENT_SPEC_FILE = internal/eventgen/json/events-2.0.0.json

all: api clients extensions

protobuf: pb/ari.proto
	protoc -I. --go_out=plugins=grpc:. pb/ari.proto

api:
	go build ./
	go build ./stdbus

test:
	go test `go list ./... | grep -v /vendor/`

check: all
	gometalinter --disable=gotype client/native ext/...

clients:
	go build ./client/native
	go build ./client/arimocks

extensions:
	go build ./ext/play
	go build ./ext/record

events:
	go build -o bin/eventgen ./internal/eventgen/...
	@./bin/eventgen internal/eventgen/template.tmpl ${EVENT_SPEC_FILE} |goimports > events_gen.go
	
mock:
	go get -u github.com/vektra/mockery/.../
	rm -Rf vendor/ client/arimocks
	mockery -all -outpkg arimocks -output client/arimocks
	glide i

ci: check
