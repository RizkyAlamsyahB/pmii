
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
# Download dependency
RUN go mod download

COPY . .

# CGO_ENABLED=0 wajib buat Alpine biar bisa jalan tanpa library C eksternal
RUN CGO_ENABLED=0 GOOS=linux go build -o pmii-backend cmd/api/main.go

# Gunakan image Alpine kosong (Super Kecil)
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/pmii-backend .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/seeds ./seeds

EXPOSE 8080

# Jalankan aplikasi
CMD ["./pmii-backend"]