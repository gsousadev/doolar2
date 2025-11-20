FROM golang:1.25 AS builder

WORKDIR /app

# Copiar arquivos do módulo Go
COPY app/go.mod app/go.sum ./
RUN go mod download

# Copiar código fonte
COPY app/ ./

EXPOSE 8080

CMD ["go", "run", "./cmd/http"]