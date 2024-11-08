# Etapa de build - usa uma imagem com Go para compilar o projeto, especificando a arquitetura
FROM --platform=linux/amd64 golang:1.21.3-alpine AS builder

# Cria e define o diretório de trabalho dentro do container
WORKDIR /app

# Copia o arquivo go.mod e go.sum para instalar as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código do projeto para o container
COPY . .

# Compila o projeto Go dentro da pasta cmd, gerando um binário chamado "app"
RUN go build -o app ./cmd

# Etapa final - cria uma imagem mínima para executar o binário
FROM alpine:latest

# Define o diretório de trabalho no container final
WORKDIR /root/

# Copia o binário gerado na etapa de build para o novo container
COPY --from=builder /app/app .

# Copia o arquivo .env da pasta /cmd para o diretório de trabalho do container
#COPY ./cmd/.env ./

# Expõe a porta que o servidor Go utiliza (ajuste conforme o projeto)
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./app"]
