FROM lastfm/base:latest

WORKDIR /app

COPY static /static

RUN go build -o /lastfm-srv ./cmd/lastfm-srv

WORKDIR /

CMD /lastfm-srv
