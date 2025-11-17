package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig contém as configurações do MongoDB
type MongoConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// NewMongoConnection cria uma nova conexão com MongoDB
func NewMongoConnection(cfg MongoConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Testar a conexão
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("MongoDB connection established successfully")

	return client, nil
}

// DefaultMongoConfig retorna a configuração padrão para desenvolvimento
func DefaultMongoConfig() MongoConfig {
	return MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "doolar",
		Timeout:  10 * time.Second,
	}
}

// GetDatabase retorna a instância do banco de dados
func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}
