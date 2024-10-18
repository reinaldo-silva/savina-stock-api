# Usa a imagem Golang
FROM golang:1.20-alpine

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos go.mod e go.sum para instalar as dependências primeiro
COPY go.mod go.sum ./

# Instala as dependências
RUN go mod tidy

# Copia todo o conteúdo da pasta raiz para o diretório de trabalho
COPY . .

# Define o diretório para compilar a aplicação a partir de cmd/api/main.go
RUN go build -o myapp ./cmd/api

# Expõe a porta 8080
EXPOSE 8080

# Define o comando para rodar a aplicação
CMD ["./myapp"]
