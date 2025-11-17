package main

import (
	"log"

	"github.com/gsousadev/doolar2/internal/domain/entity/task_list"
	"github.com/gsousadev/doolar2/internal/infrastructure/database"
)

func main() {
	// Conectar ao banco
	db, err := database.NewGormConnection(database.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Auto migrate
	db.AutoMigrate(&database.TaskListGormModel{})

	// Criar repositório com Unit of Work
	repo := database.NewTaskListGormRepository(db).(*database.TaskListGormRepository)

	// Criar task lists
	taskList1 := task_list.NewTaskListEntity("Compras do Supermercado")
	taskList2 := task_list.NewTaskListEntity("Tarefas de Casa")
	taskList3 := task_list.NewTaskListEntity("Projetos Pessoais")

	// Adicionar operações à pilha (não executa ainda!)
	log.Println("Adicionando operações à pilha...")
	repo.Add(taskList1)
	repo.Add(taskList2)
	repo.Add(taskList3)

	log.Printf("Operações pendentes: %d\n", repo.PendingCount())

	// Executar todas as operações em uma transação
	log.Println("Executando Flush (transação)...")
	if err := repo.Flush(); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Todas as operações foram executadas com sucesso!")

	// Buscar todas
	all, _ := repo.FindAll()
	log.Printf("Total de task lists no banco: %d\n", len(all))

	for _, tl := range all {
		log.Printf("- %s (ID: %s)\n", tl.Title, tl.ID)
	}

	// Atualizar uma task list
	taskList1.Title = "Lista de Compras Atualizada"
	repo.Update(taskList1)

	// Remover outra
	repo.Remove(taskList3.ID.String())

	log.Printf("Operações pendentes antes do segundo Flush: %d\n", repo.PendingCount())

	// Flush novamente
	if err := repo.Flush(); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Update e Delete executados!")

	// Verificar resultado final
	all, _ = repo.FindAll()
	log.Printf("Total final de task lists: %d\n", len(all))
}
