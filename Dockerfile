FROM golang:1.15 AS builder

RUN apt-get update && apt-get install -y \
    libjpeg-dev

COPY . /app
WORKDIR /app

RUN go build -ldflags="-extldflags=-static" -o miru-search ./cmd/miru-search/
RUN go build -ldflags="-extldflags=-static" -o miru-index ./cmd/miru-index/
RUN go build -ldflags="-extldflags=-static" -o miru-plot ./cmd/miru-plot/

FROM scratch
COPY  --from=builder /app/miru-search /usr/bin/miru-search
COPY  --from=builder /app/miru-index /usr/bin/miru-index
COPY  --from=builder /app/miru-plot /usr/bin/miru-plot