FROM debian:8

RUN apt-get update
RUN apt-get install -y build-essential openssl libxml2-dev libncurses5-dev uuid-dev sqlite3 libsqlite3-dev pkg-config curl libjansson-dev

RUN curl -s http://downloads.asterisk.org/pub/telephony/asterisk/releases/asterisk-14.0.0-rc1.tar.gz | tar xz

WORKDIR /asterisk-14.0.0-rc1
RUN ./configure; make; make install; make samples

COPY http.conf /etc/asterisk/http.conf
COPY ari.conf /etc/asterisk/ari.conf
COPY extensions.conf /etc/asterisk/extensions.conf

CMD ["/usr/sbin/asterisk", "-f"]

