package app

import (
	"context"
	"time"

	"github.com/gsousadev/doolar2/internal/tasks/application"
	"github.com/gsousadev/doolar2/internal/tasks/infrastructure/database"
	"github.com/gsousadev/doolar2/internal/tasks/presentation"
	"github.com/gsousadev/doolar2/tools"
)

func main() {

	databaseConfig()
}

func databaseConfig() database.MongoConfig,  {
	mongoConfig := database.MongoConfig{
		URI:      tools.GetEnv("MONGO_URI", "mongodb://root:root@db:27017"),
		Database: tools.GetEnv("DB_NAME", "doolar"),
		Timeout:  10 * time.Second,
	}

	mongoClient, err := database.NewMongoConnection(mongoConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func dependencyInjection() {
	// Configuração do repositório
	taskListRepository := database.NewTaskListMongoRepository(mongoClient, mongoConfig.Database)

	// Configuração do serviço
	taskManagerService := application.NewTaskManagerService(taskListRepository)

	// Configuração do handler
	taskManagerHandler := presentation.NewTaskManagerHandler(taskManagerService)
	
}
