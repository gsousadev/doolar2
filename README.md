# Doolar - Domain-Driven Design em Go

Uma aplicação Go implementando os princípios de Domain-Driven Design (DDD) para automação residencial e gerenciamento de tarefas.

## 🚀 Funcionalidades

- **Gerenciamento de Tarefas**: Crie e gerencie tarefas com limites de tempo e rastreamento de status
- **Listas de Tarefas**: Organize múltiplas tarefas em listas
- **Entidades de Domínio**: Entidades Person, Family Member, Device, Room, Event e Rule
- **Varredura de Rede**: Ferramenta CLI para descoberta de rede
- **Arquitetura Limpa**: Separação de responsabilidades com camadas de domínio, aplicação e infraestrutura

## 📁 Estrutura do Projeto

```
doolar-golang/
├── cmd/                      # Comandos CLI
│   ├── root.go              # Configuração do comando raiz
│   └── networkScan.go       # Comando de varredura de rede
├── internal/
│   ├── application/         # Serviços de aplicação
│   │   └── networkScan.go
│   ├── domain/              # Camada de domínio (lógica de negócio)
│   │   ├── entity/          # Entidades de domínio
│   │   │   ├── entity.go    # Entidade base com UUID
│   │   │   ├── person.go
│   │   │   ├── family_member.go
│   │   │   ├── device.go
│   │   │   ├── room.go
│   │   │   ├── event.go
│   │   │   ├── rule.go
│   │   │   └── task_list/   # Gerenciamento de tarefas
│   │   │       ├── task_entity.go
│   │   │       ├── task_with_time_limit_entity.go
│   │   │       ├── home_task_entity.go
│   │   │       ├── home_task_with_time_limit_entity.go
│   │   │       └── task_list_entity.go
│   │   └── valueObject/     # Objetos de valor
│   │       ├── action.go
│   │       ├── condition.go
│   │       ├── geographic_point.go
│   │       └── slug_value_object.go
│   └── infrastructure/      # Camada de infraestrutura
│       └── database/
│           └── person_repository.go
├── go.mod
├── go.sum
└── main.go
```

## 🛠️ Stack Tecnológico

- **Go 1.21+**
- **Cobra** - Framework CLI
- **UUID v6** - Identificadores únicos para entidades
- **Testify** - Framework de testes

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Git

## ⚙️ Instalação

1. Clone o repositório:
```bash
git clone https://github.com/gsousadev/doolar2.git
cd doolar-golang
```

2. Instale as dependências:
```bash
go mod download
```

3. Compile a aplicação:
```bash
go build -o doolar
```

## 🎯 Uso

### Execute a aplicação:
```bash
./doolar
```

### Varredura de Rede:
```bash
./doolar networkScan
```

### Execute os testes:
```bash
go test ./...
```

### Execute os testes com cobertura:
```bash
go test -cover ./...
```

## 🏗️ Entidades de Domínio

### Gerenciamento de Tarefas

- **TaskEntity**: Tarefa base com título, descrição e status (pendente, em progresso, concluída, cancelada)
- **TimedTaskEntity**: Tarefa com datas de início e fim
- **HomeTaskEntity**: Tarefa para automação residencial
- **TaskListEntity**: Coleção de tarefas

### Automação Residencial

- **Person**: Entidade de usuário com nome e informações de contato
- **FamilyMember**: Pessoa associada a uma família
- **Device**: Dispositivos IoT para automação residencial
- **Room**: Espaços físicos em uma residência
- **Event**: Eventos e gatilhos do sistema
- **Rule**: Regras de automação baseadas em condições e ações

## 🧪 Testes

O projeto inclui testes unitários abrangentes para todas as entidades de domínio:

```bash
# Execute todos os testes
go test ./...

# Execute os testes com saída detalhada
go test -v ./...

# Execute os testes de um pacote específico
go test ./internal/domain/entity/task_list/...
```

## 🎨 Padrões de Design

- **Domain-Driven Design (DDD)**: Separação clara entre domínio, aplicação e infraestrutura
- **Padrão Entity**: Todas as entidades herdam de uma Entity base com UUID
- **Value Objects**: Objetos imutáveis para conceitos como pontos geográficos e slugs
- **Padrão Repository**: Abstração de acesso a dados
- **Padrão Command**: Comandos CLI usando Cobra

## 📝 Sobre o Projeto

Este é um projeto pessoal de estudos e portfólio, desenvolvido para demonstrar conhecimentos em Go e Domain-Driven Design.

## 👨‍💻 Autor

**Guilherme Santos** - [@gsousadev](https://github.com/gsousadev)

## 📧 Contato

- Email: gsousadev@gmail.com
- GitHub: [@gsousadev](https://github.com/gsousadev)

---

Desenvolvido com ❤️ usando Go e princípios de Domain-Driven Design
