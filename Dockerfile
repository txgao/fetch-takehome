FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build 
RUN go build -o main ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 4000

CMD ["/app/main"]
