# Doolar - Domain-Driven Design em Go

Uma aplicação Go implementando os princípios de Domain-Driven Design (DDD) e Clean Architecture para gerenciamento de tarefas com múltiplas implementações de persistência.

## 🚀 Funcionalidades

- **Gerenciamento de Tarefas**: Crie e gerencie listas de tarefas com status e estatísticas
- **API REST**: Endpoints HTTP para todas as operações de tarefas
- **Múltiplas Persistências**: Suporte para GORM (PostgreSQL) e MongoDB
- **Unit of Work Pattern**: Transações atômicas com operações enfileiradas
- **Clean Architecture**: Separação completa entre domínio, aplicação, infraestrutura e apresentação
- **Dependency Inversion**: Camadas dependem de abstrações, não de implementações concretas

## 📁 Estrutura do Projeto

```
app/
├── cmd/
│   └── http/                        # Entry point HTTP
│       ├── main.go                  # Composition Root - orquestra dependências
│       └── router.go                # Configuração de rotas
├── internal/
│   ├── domain/                      # Camada de Domínio (regras de negócio)
│   │   ├── entity/
│   │   │   ├── entity.go           # Entidade base com UUID v6
│   │   │   └── task_list/
│   │   │       ├── task_entity.go           # Tarefa com status
│   │   │       └── task_list_entity.go      # Aggregate Root
│   │   ├── repository/
│   │   │   └── task_list_repository.go      # Interface do repositório
│   │   └── valueObject/
│   ├── application/                 # Camada de Aplicação (casos de uso)
│   │   ├── task_manager_interface.go        # Interface TaskManager
│   │   └── task_manager_service.go          # Implementação dos use cases
│   ├── infrastructure/              # Camada de Infraestrutura (detalhes técnicos)
│   │   └── database/
│   │       ├── connection.go                # GORM connection
│   │       ├── mongo_connection.go          # MongoDB connection
│   │       ├── task_list_gorm_repository.go # Repositório GORM + Unit of Work
│   │       └── task_list_mongo_repository.go # Repositório MongoDB + Unit of Work
│   └── presentation/                # Camada de Apresentação (HTTP handlers)
│       ├── task_manager_handler.go          # Handlers HTTP
│       └── task_presenter.go                # DTOs de resposta
└── docker/
    ├── Dockerfile                   # Multi-stage build
    └── docker-compose.yaml          # MongoDB + Replica Set
```

## 🛠️ Stack Tecnológico

- **Go 1.24+**
- **GORM** - ORM para PostgreSQL
- **MongoDB Driver** - Driver oficial do MongoDB
- **UUID v6** - Identificadores únicos para entidades
- **Testify** - Framework de testes
- **Docker** - Containerização

## 📋 Pré-requisitos

- Go 1.24 ou superior
- Docker e Docker Compose (para MongoDB)
- PostgreSQL (opcional, para usar GORM)

## ⚙️ Instalação e Execução

### 1. Clone o repositório:
```bash
git clone https://github.com/gsousadev/doolar-golang.git
cd doolar-golang/app
```

### 2. Instale as dependências:
```bash
go mod download
```

### 3. Configure o MongoDB com Docker:
```bash
cd ../docker
docker-compose up -d
```

### 4. Execute o servidor HTTP:
```bash
cd ../app
go run cmd/http/main.go
```

O servidor iniciará em `http://localhost:8080` e mostrará logs coloridos:

```
2025/11/18 10:30:15 ✓ Connected to MongoDB
2025/11/18 10:30:15 🚀 Server starting on http://localhost:8080

2025/11/18 10:30:20 → POST /task-lists from 127.0.0.1:54321
2025/11/18 10:30:20 ← POST /task-lists [201] 15ms (342 bytes)
```

### 5. Teste a API:

**Opção A - Postman (Recomendado):**

Importe a coleção pronta no Postman:
```bash
# Arquivos para importar:
Doolar_API.postman_collection.json       # Coleção com todas as rotas
Doolar_Local.postman_environment.json    # Environment local
```

Veja o guia completo: [POSTMAN_GUIDE.md](POSTMAN_GUIDE.md)

**Opção B - Script de Teste:**

Use o script bash incluído:
```bash
./test-api.sh
```

**Opção C - cURL Manual:**

```bash
# Criar lista de tarefas
curl -X POST http://localhost:8080/task-lists \
  -H "Content-Type: application/json" \
  -d '{"name": "Tarefas de Casa", "description": "Lista de tarefas domésticas"}'
```

Ou teste manualmente:
```bash
curl -X POST http://localhost:8080/task-lists \
  -H "Content-Type: application/json" \
  -d '{"title": "My Shopping List"}'
```

## 🌐 API REST

### Endpoints Disponíveis

```bash
# Criar lista de tarefas
POST /task-lists
Content-Type: application/json
{
  "title": "Minha Lista"
}

# Buscar lista por ID
GET /task-lists/{id}

# Adicionar tarefa à lista
POST /task-lists/{id}/tasks
Content-Type: application/json
{
  "title": "Estudar Go",
  "description": "Aprender sobre interfaces"
}

# Listar tarefas pendentes
GET /task-lists/{id}/tasks/pending

# Atualizar status de uma tarefa
PATCH /task-lists/{id}/tasks/{taskId}/status
Content-Type: application/json
{
  "status": "in_progress"
}
# Status: pending, in_progress, completed, cancelled

# Obter estatísticas da lista
GET /task-lists/{id}/statistics

# Deletar lista
DELETE /task-lists/{id}
```

### Exemplo de Resposta

```json
{
  "message": "Task list created successfully",
  "data": {
    "id": "01JCXXX...",
    "title": "Minha Lista",
    "tasks": [
      {
        "id": "01JCYYY...",
        "title": "Estudar Go",
        "description": "Aprender sobre interfaces",
        "status": "pending"
      }
    ],
    "stats": {
      "total": 1,
      "pending": 1,
      "in_progress": 0,
      "completed": 0,
      "cancelled": 0
    }
  }
}
```

## 🏗️ Arquitetura

### Composition Root (cmd/http/main.go)

O `main.go` funciona como **Composition Root**, orquestrando todas as dependências:

```go
func main() {
    // 1. Conecta ao MongoDB
    mongoClient, err := database.NewMongoConnection(mongoConfig)
    
    // 2. Cria repositório (implementação)
    taskListRepository := database.NewTaskListMongoRepository(mongoClient, dbName)
    
    // 3. Cria serviço de aplicação (retorna interface)
    taskManagerService := application.NewTaskManagerService(taskListRepository)
    
    // 4. Cria handler HTTP (depende da interface)
    taskManagerHandler := presentation.NewTaskManagerHandler(taskManagerService)
    
    // 5. Configura rotas e inicia servidor
    router := SetupRouter(taskManagerHandler)
    server.ListenAndServe()
}
```

### Dependency Inversion Principle

```
Presentation → TaskManager (interface)
                    ↑
                    │ implementa
                    │
Application → TaskManagerService
                    ↓
                depende de
                    ↓
Domain → TaskListRepository (interface)
                    ↑
                    │ implementam
                    │
Infrastructure → GormRepository | MongoRepository
```

### Camadas

1. **Domain**: Entidades e interfaces de repositório (regras de negócio puras)
2. **Application**: Use cases e orquestração (retorna entidades diretamente)
3. **Infrastructure**: Implementações de persistência (GORM, MongoDB)
4. **Presentation**: Handlers HTTP e DTOs (transforma entidades em respostas)

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Testes com cobertura
go test -cover ./...

# Testes de um pacote específico
go test ./internal/infrastructure/database/...

# Testes com output verbose
go test -v ./...
```

## 🐳 Docker

### Build da aplicação:
```bash
cd docker
docker-compose build
```

### Executar com Docker:
```bash
docker-compose up
```

## 🔧 Variáveis de Ambiente

```bash
# MongoDB
MONGO_URI=mongodb://localhost:27017
DB_NAME=doolar

# Servidor HTTP
PORT=8080
```

## 📊 Logging

O projeto possui um **sistema de logging interno** que salva logs em arquivos locais e exibe no console simultaneamente.

### Características

- ✅ **Dual Output**: Logs salvos em arquivo E exibidos no console
- ✅ **Níveis de Log**: DEBUG, INFO, WARN, ERROR
- ✅ **Rotação Automática**: Cria novos arquivos ao atingir tamanho máximo
- ✅ **Timestamps**: Cada log com data e hora precisas
- ✅ **Cores no Console**: Visualização colorida (INFO=verde, WARN=amarelo, ERROR=vermelho)
- ✅ **Logging HTTP**: Middleware que loga automaticamente todas as requisições

### Estrutura de Logs

Os logs são salvos no diretório `app/logs/`:

```
app/logs/
├── app-2025-11-18.log           # Log do dia atual
├── app-2025-11-18-153045.log    # Log rotacionado (quando atinge 10MB)
└── app-2025-11-17.log           # Log do dia anterior
```

### Exemplo de Output

**Console com cores:**
```
[2025-11-18 15:30:45] INFO  === Iniciando aplicação Doolar ===
[2025-11-18 15:30:45] INFO  Conectando ao MongoDB: mongodb://db:27017/doolar
[2025-11-18 15:30:46] INFO  ✓ Conectado ao MongoDB com sucesso
[2025-11-18 15:30:46] INFO  🚀 Servidor HTTP iniciado em http://localhost:8080
[2025-11-18 15:30:50] INFO  → POST /task-lists from 127.0.0.1:54321
[2025-11-18 15:30:51] INFO  ← POST /task-lists [201] 150ms (342 bytes)
[2025-11-18 15:31:00] WARN  ← GET /invalid [404] 2ms (23 bytes)
```

**Arquivo (app/logs/app-2025-11-18.log):**
```
[2025-11-18 15:30:45] INFO  === Iniciando aplicação Doolar ===
[2025-11-18 15:30:45] INFO  Conectando ao MongoDB: mongodb://db:27017/doolar
[2025-11-18 15:30:46] INFO  ✓ Conectado ao MongoDB com sucesso
[2025-11-18 15:30:46] INFO  🚀 Servidor HTTP iniciado em http://localhost:8080
[2025-11-18 15:30:50] INFO  → POST /task-lists from 127.0.0.1:54321
[2025-11-18 15:30:51] INFO  ← POST /task-lists [201] 150ms (342 bytes)
[2025-11-18 15:31:00] WARN  ← GET /invalid [404] 2ms (23 bytes)
```

### Testar o Logger

Execute o script de teste para ver o logger em ação:

```bash
# Terminal 1 - Inicia o servidor
cd docker && docker compose up -d
cd ../app && go run cmd/http/main.go

# Terminal 2 - Executa testes
./test-logger.sh
```

**Visualizar logs em tempo real:**
```bash
tail -f app/logs/app-$(date +%Y-%m-%d).log
```

**Filtrar logs por nível:**
```bash
grep "ERROR" app/logs/*.log
grep "WARN\|ERROR" app/logs/*.log
```

### Documentação Completa

Para mais detalhes sobre configuração e uso avançado, veja:
- [Logger Interno - Documentação](app/internal/logger/README.md)

## 📚 Padrões de Design Implementados

- **Domain-Driven Design (DDD)**: Aggregate Root, Entities, Value Objects
- **Clean Architecture**: Separação de responsabilidades em camadas
- **Unit of Work**: Transações atômicas com operações enfileiradas
- **Repository Pattern**: Abstração de persistência
- **Dependency Inversion**: Dependências via interfaces
- **Data Mapper**: Separação entre modelo de domínio e persistência

## 🎯 Princípios SOLID

✅ **S**ingle Responsibility: Cada camada tem uma responsabilidade única  
✅ **O**pen/Closed: Fácil adicionar novos repositórios ou handlers  
✅ **L**iskov Substitution: Interfaces respeitadas pelas implementações  
✅ **I**nterface Segregation: Interfaces específicas por necessidade  
✅ **D**ependency Inversion: Depende de abstrações, não de implementações  

## 📖 Documentação Adicional

Consulte o arquivo [ARCHITECTURE.md](ARCHITECTURE.md) para uma explicação detalhada da arquitetura, fluxo de dados e decisões de design.

## 👨‍💻 Autor

**Guilherme Santos** - [@gsousadev](https://github.com/gsousadev)

## 📧 Contato

- Email: gsousadev@gmail.com
- GitHub: [@gsousadev](https://github.com/gsousadev)

---

Desenvolvido com ❤️ usando Go, DDD e Clean Architecture
