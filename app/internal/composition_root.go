package internal

import (
	"time"

	"github.com/gsousadev/doolar2/internal/shared/infrastructure/database"
	"github.com/gsousadev/doolar2/internal/tasks/application"
	"github.com/gsousadev/doolar2/internal/tasks/application/contracts"
	"github.com/gsousadev/doolar2/internal/tasks/infrastructure/database/mongo"
	mongo_driver "go.mongodb.org/mongo-driver/mongo"
)

type CompositionRoot struct {
	TaskListService contracts.ITaskManagerService
}

func NewCompositionRoot() *CompositionRoot {

	taskListRepository := mongo.NewTaskListMongoRepository(newMongoConnection(), "doolar")

	return &CompositionRoot{
		TaskListService: application.NewTaskManagerService(taskListRepository),
	}
}

func newMongoConnection() *mongo_driver.Client {

	duration, _ := time.ParseDuration("1m")

	mongoClient, error := database.NewMongoConnection(database.MongoConfig{
		URI:      "mongodb://root:root@db:27017/",
		Database: "doolar",
		Timeout:  duration,
	})

	if error != nil {
		panic("Erro em conexão com o mongo db")
	}

	return mongoClient
}
