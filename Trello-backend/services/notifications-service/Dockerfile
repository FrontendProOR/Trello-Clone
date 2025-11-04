# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder  

WORKDIR /app

# Kopiraj Go module i fajlove
COPY go.mod go.sum ./

# Preuzmi Go module zavisnosti
RUN go mod download

# Kopiraj sve fajlove iz direktorijuma direktno (ne koristi src)
COPY . .

# Izgradnja aplikacije
RUN GOOS=linux GOARCH=amd64 go build -o notifications-service main.go

# Stage 2: Kreiraj manju sliku za runtime
FROM alpine:latest  

WORKDIR /root/

# Instaliraj MongoDB alatke
RUN apk update && apk add --no-cache mongodb-tools

# Kopiraj izgrađeni binarni fajl
COPY --from=builder /app/notifications-service .

# Dodeli permisije
RUN chmod +x notifications-service

# Komanda koja se pokreće
CMD ["./notifications-service"]
