# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia i file del modulo di Go
COPY go.mod* ./
# Se hai dipendenze, decommentare:
# RUN go mod download

# Copia il codice sorgente
COPY . .

# Compila l'applicazione
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia il binario dalla build stage
COPY --from=builder /app/main .

# Esponi la porta
EXPOSE 8080

# Comando di avvio
CMD ["./main"]