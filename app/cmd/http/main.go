package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gsousadev/doolar2/internal/application"
	"github.com/gsousadev/doolar2/internal/infrastructure/database"
	"github.com/gsousadev/doolar2/internal/presentation"
	"github.com/rs/cors"
)

func main() {

	mongoConfig := database.MongoConfig{
		URI:      getEnv("MONGO_URI", "mongodb://root:root@db:27017"),
		Database: getEnv("DB_NAME", "doolar"),
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

	taskListRepository := database.NewTaskListMongoRepository(mongoClient, mongoConfig.Database)

	taskManagerService := application.NewTaskManagerService(taskListRepository)

	taskManagerHandler := presentation.NewTaskManagerHandler(taskManagerService)

	router := SetupRouter(taskManagerHandler)

	// 9. Configuração do servidor
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // aberto
		AllowedHeaders:   []string{"*"}, // qualquer header
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})
	server.Handler = c.Handler(router)

	// 10. Iniciar o servidor
	log.Printf("Servidor iniciado na porta %s\n", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	// 11. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Erro ao desligar servidor: %v", err)
	}

}

// getEnv retorna uma variável de ambiente ou valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupRouter configura as rotas HTTP
func SetupRouter(handler *presentation.TaskManagerHandler) http.Handler {
	mux := http.NewServeMux()

	// Task Lists
	mux.HandleFunc("/task-lists", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTaskList(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			presentation.UploadAudio(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/task-lists/", func(w http.ResponseWriter, r *http.Request) {
		// Verifica se é uma rota específica
		if r.URL.Path == "/task-lists/" {
			http.Error(w, "Task list ID required", http.StatusBadRequest)
			return
		}

		// /task-lists/{id}
		if r.Method == http.MethodGet && !hasSubpath(r.URL.Path, "/task-lists/", "tasks", "statistics") {
			handler.GetTaskList(w, r)
			return
		}

		// /task-lists/{id}
		if r.Method == http.MethodDelete && !hasSubpath(r.URL.Path, "/task-lists/", "tasks", "statistics") {
			handler.DeleteTaskList(w, r)
			return
		}

		// /task-lists/{id}/tasks
		if r.Method == http.MethodPost && containsPath(r.URL.Path, "/tasks") && !containsPath(r.URL.Path, "/pending") && !containsPath(r.URL.Path, "/status") {
			handler.AddTaskToList(w, r)
			return
		}

		// /task-lists/{id}/tasks/pending
		if r.Method == http.MethodGet && containsPath(r.URL.Path, "/tasks/pending") {
			handler.GetPendingTasks(w, r)
			return
		}

		// /task-lists/{id}/tasks/{taskId}/status
		if r.Method == http.MethodPatch && containsPath(r.URL.Path, "/tasks/") && containsPath(r.URL.Path, "/status") {
			handler.UpdateTaskStatus(w, r)
			return
		}

		// /task-lists/{id}/statistics
		if r.Method == http.MethodGet && containsPath(r.URL.Path, "/statistics") {
			handler.GetStatistics(w, r)
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	})

	return mux
}

// Helper para verificar se o path contém um subpath específico
func hasSubpath(path string, prefix string, subpaths ...string) bool {
	for _, sub := range subpaths {
		if containsPath(path, "/"+sub) {
			return true
		}
	}
	return false
}

func containsPath(path, substring string) bool {
	for i := 0; i < len(path)-len(substring)+1; i++ {
		if path[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}
