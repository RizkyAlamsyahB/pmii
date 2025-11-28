
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
# Download dependency
RUN go mod download

COPY . .

# CGO_ENABLED=0 wajib buat Alpine biar bisa jalan tanpa library C eksternal
RUN CGO_ENABLED=0 GOOS=linux go build -o pmii-backend main.go

# Gunakan image Alpine kosong (Super Kecil)
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/pmii-backend .

EXPOSE 8080

# Jalankan aplikasi
CMD ["./pmii-backend"]