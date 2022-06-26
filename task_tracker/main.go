package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gsom95/go-ms/task_tracker/server"
)

func main() {
	mux := http.NewServeMux()
	s := server.NewTaskServer()

	mux.HandleFunc("/task/", server.BasicAuth(s.TaskHandler))

	port := os.Getenv("SERVERPORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Printf("err: %v\n", err)
	}
}
