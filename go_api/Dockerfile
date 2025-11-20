# Etapa 1: Build (Compilação)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependência
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário. 
# CGO_ENABLED=0 garante que o binário seja estático (roda em qualquer linux)
# -ldflags="-s -w" remove informações de debug para diminuir o tamanho
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server main.go

# Etapa 2: Runtime (Execução)
FROM alpine:latest

WORKDIR /root/

# Copia apenas o executável da etapa anterior
COPY --from=builder /app/server .

# Expõe a porta
EXPOSE 8080

# Comando para rodar
CMD ["./server"]