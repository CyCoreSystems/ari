

all:
	CGO_ENABLED=0 GOOS=linux \
		go build -a -installsuffix cgo \
			-ldflags '-w -extld ld -extldflags -static' \
			-a -o natsgw/app.static natsgw/main.go
	CGO_ENABLED=0 GOOS=linux \
		go build -a -installsuffix cgo \
			-ldflags '-w -extld ld -extldflags -static' \
			-a -o stasisStart-nats/app.static stasisStart-nats/main.go

