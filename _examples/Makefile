
stasisStart-nats: .PHONY
	CGO_ENABLED=0 GOOS=linux \
		go build -a -installsuffix cgo \
			-ldflags '-w -extld ld -extldflags -static' \
			-a -o stasisStart-nats/app.static stasisStart-nats/main.go
	docker-compose -f stasisStart-nats/docker-compose.yml build
	docker-compose -f stasisStart-nats/docker-compose.yml up

sendDtmf-nats: .PHONY
	CGO_ENABLED=0 GOOS=linux \
		go build -a -installsuffix cgo \
			-ldflags '-w -extld ld -extldflags -static' \
			-a -o sendDtmf-nats/app.static sendDtmf-nats/main.go
	docker-compose -f sendDtmf-nats/docker-compose.yml build
	docker-compose -f sendDtmf-nats/docker-compose.yml up

stasisStart: .PHONY
	CGO_ENABLED=0 GOOS=linux \
		go build -a -installsuffix cgo \
			-ldflags '-w -extld ld -extldflags -static' \
			-a -o stasisStart/app.static stasisStart/main.go
	docker-compose -f stasisStart/docker-compose.yml build
	docker-compose -f stasisStart/docker-compose.yml up

natsgw: .PHONY
	CGO_ENABLED=0 GOOS=linux \
		go build -a -installsuffix cgo \
			-ldflags '-w -extld ld -extldflags -static' \
			-a -o natsgw/app.static natsgw/main.go


.PHONY:

