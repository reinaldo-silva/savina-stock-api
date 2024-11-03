# Usa a imagem base do Go
FROM golang:1.20-alpine AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia o código da aplicação
COPY . .

# Gera o go.mod e go.sum
RUN go mod init github.com/reinaldo-silva/savina-stock || true
RUN go mod tidy

# Constrói a aplicação
RUN go build -o myapp ./cmd/api/main.go

# Nova etapa para a imagem final
FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /app

# Define a variável de ambiente para o ambiente de produção
ENV ENVIRONMENT=production

# Copia o executável da etapa de build
COPY --from=builder /app/myapp .

# Expõe a porta que sua aplicação irá utilizar
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./myapp"]
