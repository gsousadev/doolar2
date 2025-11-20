FROM golang:1.25 AS builder

WORKDIR /app

# Copiar arquivos do módulo Go
COPY app/go.mod app/go.sum ./
RUN go mod download

# Copiar código fonte
COPY app/ ./

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o doolar-api ./cmd/http

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário da etapa de build
COPY --from=builder /app/doolar-api .

EXPOSE 8080

CMD ["./doolar-api"]