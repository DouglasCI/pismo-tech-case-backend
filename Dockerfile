FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 simplifies setup, skipping C compiler.
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]