FROM lastfm/base:latest

WORKDIR /app

COPY static /static

ARG tls_cert_rel
ENV TLS_CERT_PATH /cert/${tls_cert_rel}

RUN go build -o /lastfm-srv ./cmd/lastfm-srv

WORKDIR /

CMD /lastfm-srv
