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

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	log.Println("MongoDB connection established successfully")

	return client, nil

}
