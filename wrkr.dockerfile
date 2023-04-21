FROM lastfm/base:latest

WORKDIR /app

RUN go build -o /lastfm-wrkr ./cmd/lastfm-wrkr

WORKDIR /

CMD /lastfm-wrkr
