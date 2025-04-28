FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /app

COPY go.mod .
# COPY go.sum .
# RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.21
WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 11555

CMD ["./server"]
