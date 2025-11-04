FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src
COPY src/config/password_blacklist.txt ./src/config/password_blacklist.txt  

RUN GOOS=linux GOARCH=amd64 go build -o users-service ./src/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk update && apk add --no-cache mongodb-tools

COPY --from=builder /app/users-service .
COPY --from=builder /app/src/config/password_blacklist.txt ./src/config/password_blacklist.txt  
RUN chmod +x users-service

CMD ["./users-service"]
