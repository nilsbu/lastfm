FROM lastfm/base:latest

WORKDIR /app

RUN go build -o /lastfm-srv ./cmd/lastfm-srv

WORKDIR /

CMD /lastfm-srv
