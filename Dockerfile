FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o payment-gateway ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/payment-gateway .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./payment-gateway"]