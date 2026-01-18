FROM golang:1.25

WORKDIR /app

# Copiar arquivos do módulo Go
COPY app/go.mod app/go.sum ./
RUN go mod download

# Copiar código fonte
COPY app/ ./

EXPOSE 8080

CMD ["tail","-f","/dev/null"]