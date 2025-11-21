# Etapa 1: Build (Compilação)
FROM golang:alpine AS builder

WORKDIR /app

# Copia os arquivos de dependência
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# --- LINHA NOVA AQUI EM BAIXO ---
# Arruma as dependências antes de compilar
RUN go mod tidy

# Compila o binário. 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server main.go

# Etapa 2: Runtime (Execução) - IGUAL AO ANTERIOR
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]