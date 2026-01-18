package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gsousadev/doolar-golang/internal"
	sharedPresentation "github.com/gsousadev/doolar-golang/internal/shared/presentation"
	taskPresentation "github.com/gsousadev/doolar-golang/internal/tasks/presentation/http"
	"github.com/gsousadev/doolar-golang/tools"
	"github.com/rs/cors"
)

func StartServer() {

	compositionRoot := internal.NewCompositionRoot()

	taskHandler := taskPresentation.NewTaskManagerHandler(compositionRoot.TaskListService)
	healthHandler := sharedPresentation.NewHealthHandler()

	router := setupRouter(taskHandler, healthHandler)

	// 9. Configuração do servidor
	port := tools.GetEnv("PORT", "8080")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadTimeout:       30 * time.Second,  // Aumentado de 15s
		WriteTimeout:      300 * time.Second, // 5 minutos para streaming
		IdleTimeout:       120 * time.Second, // Aumentado de 60s
		ReadHeaderTimeout: 10 * time.Second,  // Timeout para ler headers
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

// SetupRouter configura as rotas HTTP
func setupRouter(taskHandler *taskPresentation.TaskManagerHandler, healthHandler *sharedPresentation.HealthHandler) http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler.HealthCheckHandler)
	mux.HandleFunc("POST /task-lists", taskHandler.CreateTaskList)

	return mux
}
