FROM golang:1.18 AS builder
WORKDIR /app
COPY go.* *.go ./
RUN go get ./...
RUN go build

FROM ubuntu:22.10
COPY --from=builder /app/archiver /
RUN apt update && apt -y install sqlite3
ENTRYPOINT ["/archiver"]
