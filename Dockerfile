# Etapa de build
FROM golang:1.22-alpine AS builder
WORKDIR /app
# Copia os arquivos de módulo e baixa as dependências
COPY go.mod go.sum ./
RUN go mod download
# Copia o restante do código-fonte
COPY . .
# Constrói o executável da sua aplicação Go
# O nome do executável será 'server' e estará em /app/server
RUN go build -o /app/server ./cmd/server/main.go 

# Etapa final: cria uma imagem menor e mais segura
FROM alpine:latest
WORKDIR /app
# Copia o executável da etapa de build para o diretório de trabalho final
COPY --from=builder /app/server .
# Copia o arquivo .env para o contêiner
COPY .env . 

# Expõe a porta que sua aplicação Go vai escutar
EXPOSE 8080 

# Define o comando que será executado quando o contêiner iniciar
# Certifique-se de que o nome aqui corresponda ao nome do executável criado acima.
CMD ["./server"]