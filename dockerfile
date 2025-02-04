# Etapa de construcción
FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o myapp ./cmd/server
RUN chmod +x myapp  # Asegurar permisos de ejecución

FROM debian:bookworm-slim

WORKDIR /root/

# Instalar dependencias necesarias para ejecutar el binario (glibc y otras)
RUN apt-get update && apt-get install -y libc6

COPY --from=builder /app/myapp .

EXPOSE 8080

CMD ["./myapp"]
