# Etapa 1: build do binário estaticamente
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instala as dependências necessárias para CGO
RUN apk add --no-cache gcc musl-dev

COPY main.go .
COPY go.mod .
COPY go.sum .

# Compila o binário com CGO habilitado
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o malicious

# Etapa 2: imagem final com scratch
FROM scratch

# Copia apenas o binário para a imagem final
COPY --from=builder /app/malicious /malicious

# Define o binário como ponto de entrada
ENTRYPOINT ["/malicious"]