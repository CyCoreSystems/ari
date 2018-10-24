SHELL = /usr/bin/env bash

EVENT_SPEC_FILE = internal/eventgen/json/events-2.0.0.json

all: dep api clients contributors extensions

contributors:
	write_mailmap > CONTRIBUTORS

protobuf: ari.proto
	protoc -I. -I./vendor --gogofast_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,plugins=grpc:. ari.proto

dep:
	dep ensure

api:
	go build ./
	go build ./stdbus
	go build ./rid

test:
	go test `go list ./... | grep -v /vendor/`

check: all
	golangci-lint run
	#gometalinter --disable=gotype client/native ext/...

clients:
	go build ./client/native
	go build ./client/arimocks

extensions:
	go build ./ext/audiouri
	go build ./ext/bridgemon
	go build ./ext/keyfilter
	go build ./ext/play
	go build ./ext/record

events:
	go build -o bin/eventgen ./internal/eventgen/...
	@./bin/eventgen internal/eventgen/template.tmpl ${EVENT_SPEC_FILE} |goimports > events_gen.go
	
mock:
	go get -u github.com/vektra/mockery/.../
	rm -Rf vendor/ client/arimocks
	mockery -name . -outpkg arimocks -output client/arimocks
	dep ensure

ci: check
