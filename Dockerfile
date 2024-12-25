FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mqtt-adaptor ./cmd

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/mqtt-adaptor .

CMD ["./mqtt-adaptor"]
