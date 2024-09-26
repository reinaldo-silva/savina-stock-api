# Savina Stock API

Savina Stock API é uma aplicação desenvolvida em Go que fornece uma interface para gerenciar o estoque de produtos em um e-commerce. Esta API permite realizar operações CRUD (Criar, Ler, Atualizar e Deletar) em produtos, facilitando a administração de inventário.

## Tecnologias Utilizadas

- Go (Golang)
- PostgreSQL
- Docker
- Docker Compose
- Clean Architecture
- SOLID Principles

## Funcionalidades

- Adicionar novos produtos ao estoque
- Listar produtos existentes
- Atualizar informações de produtos
- Remover produtos do estoque

## Pré-requisitos

Antes de executar o projeto, certifique-se de que você tem os seguintes itens instalados:

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Configuração do Ambiente

1. Clone o repositório:

   ```bash
   git clone <URL-do-repositório>
   cd savina-stock-api
   ```

2. Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis de ambiente:

   ```env
   DB_HOST=postgres
   DB_USER=postgres
   DB_PASSWORD=secret
   DB_NAME=ecommerce
   DB_PORT=5432
   SERVER_PORT=8080
   ```

3. Inicie o banco de dados PostgreSQL com Docker Compose:

   ```bash
   docker-compose up -d
   ```

4. Execute a aplicação Go:

   ```bash
   go run cmd/api/main.go
   ```

## Testes

Para rodar os testes, execute o seguinte comando:

```bash
go test ./...
```

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

## Licença

Este projeto está licenciado sob a Licença MIT. Consulte o arquivo [LICENSE](LICENSE) para mais detalhes.
