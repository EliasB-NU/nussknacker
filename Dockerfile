# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY src .

RUN go build -o main .

# Final Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]