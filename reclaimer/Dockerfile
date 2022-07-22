FROM golang:1.18 AS builder
WORKDIR /app
COPY go.* *.go ./
RUN go get ./...
RUN go build

FROM ubuntu:22.10
COPY --from=builder /app/reclaimer /
RUN apt update && apt -y install sqlite3
ENTRYPOINT ["/reclaimer"]
