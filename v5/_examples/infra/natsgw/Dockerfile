FROM ulexus/go-minimal
COPY app.static /app
CMD ["--ari.application","example", \
	 "--ari.username","admin", \
	 "--ari.password","admin", \
	 "--ari.websocket_url","ws://asterisk:8088/ari/events", \
	 "--ari.http_url","http://asterisk:8088/ari", \
	 "--nats.url","nats://nats:4222", \
	 "-v"]
