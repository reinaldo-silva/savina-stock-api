# Usa a imagem do Golang
FROM golang:1.20-alpine

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia apenas os arquivos go.mod e go.sum para a imagem
COPY go.mod go.sum ./

# Instala as dependências
RUN go mod tidy

# Copia todo o conteúdo do diretório atual (incluindo o código-fonte) para a imagem
COPY . .

# Altera o diretório de trabalho para o diretório onde o main.go está localizado
WORKDIR /app/cmd/api

# Compila a aplicação
RUN go build -o myapp

# Expõe a porta que a aplicação usará
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./myapp"]
