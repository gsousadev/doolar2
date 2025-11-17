package main

import (
	"context"
	"log"

	"github.com/gsousadev/doolar2/internal/domain/entity/task_list"
	"github.com/gsousadev/doolar2/internal/infrastructure/database"
)

func main() {
	// Conectar ao MongoDB
	client, err := database.NewMongoConnection(database.DefaultMongoConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Criar repositório com Unit of Work
	repo := database.NewTaskListMongoRepository(client, "doolar").(*database.TaskListMongoRepository)

	// Criar task lists
	taskList1 := task_list.NewTaskListEntity("📚 Estudar Go e MongoDB")
	taskList2 := task_list.NewTaskListEntity("🏋️ Treino da Semana")
	taskList3 := task_list.NewTaskListEntity("🛒 Lista de Compras")

	// Adicionar operações à pilha (não executa ainda!)
	log.Println("📝 Adicionando operações à pilha...")
	repo.Add(taskList1)
	repo.Add(taskList2)
	repo.Add(taskList3)

	log.Printf("⏳ Operações pendentes: %d\n", repo.PendingCount())
	log.Printf("📋 Tipos: %v\n", repo.PendingOperationTypes())

	// Executar todas as operações em uma transação MongoDB
	log.Println("\n🚀 Executando Flush (transação MongoDB)...")
	if err := repo.Flush(); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Todas as operações foram executadas com sucesso!")

	// Buscar todas
	all, _ := repo.FindAll()
	log.Printf("\n📊 Total de task lists no MongoDB: %d\n", len(all))

	for i, tl := range all {
		log.Printf("  %d. %s (ID: %s)\n", i+1, tl.Title, tl.ID)
	}

	// Atualizar uma task list
	log.Println("\n🔄 Preparando operações de Update e Delete...")
	taskList1.Title = "📚 Estudar Go, MongoDB e DDD"
	repo.Update(taskList1)

	// Remover outra
	repo.Remove(taskList3.ID.String())

	log.Printf("⏳ Operações pendentes antes do segundo Flush: %d\n", repo.PendingCount())
	log.Printf("📋 Tipos: %v\n", repo.PendingOperationTypes())

	// Flush novamente
	log.Println("\n🚀 Executando segundo Flush...")
	if err := repo.Flush(); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Update e Delete executados!")

	// Verificar resultado final
	all, _ = repo.FindAll()
	log.Printf("\n📊 Total final de task lists: %d\n", len(all))
	for i, tl := range all {
		log.Printf("  %d. %s (ID: %s)\n", i+1, tl.Title, tl.ID)
	}

	log.Println("\n🎉 Exemplo de Unit of Work com MongoDB concluído!")
}
