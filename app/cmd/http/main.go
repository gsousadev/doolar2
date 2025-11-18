package main

import (
	"log"
	"net/http"

	"github.com/gsousadev/doolar2/internal/presentation"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", presentation.GetTaskHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
