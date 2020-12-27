FROM golang:1.15 AS builder

RUN apt-get update && apt-get install -y \
    libjpeg-dev

COPY . /app
WORKDIR /app

RUN go build -ldflags="-extldflags=-static" -o miru-search ./cmd/miru-search/
RUN go build -ldflags="-extldflags=-static" -o miru-insert ./cmd/miru-insert/
RUN go build -ldflags="-extldflags=-static" -o miru-plot ./cmd/miru-plot/

FROM scratch
COPY  --from=builder /app/miru-search /usr/bin/miru-search
COPY  --from=builder /app/miru-insert /usr/bin/miru-insert
COPY  --from=builder /app/miru-plot /usr/bin/miru-plot