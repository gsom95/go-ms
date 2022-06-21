package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	taskstore "github.com/gsom95/go-ms/task_tracker/store"
)

// TaskStorer contains methods for working with Task repository.
type TaskStorer interface {
	CreateTask(text string, due time.Time) (int, error)
	GetTask(id int) (taskstore.Task, error)
	DeleteTask(id int) error
}

// TaskServer contains handlers.
type TaskServer struct {
	store TaskStorer
}

func NewTaskServer() *TaskServer {
	s := taskstore.New()
	return &TaskServer{
		store: s,
	}
}

func (ts *TaskServer) TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/task/" {
		switch r.Method {
		case http.MethodPost:
			ts.createTaskHandler(w, r)
		default:
			http.Error(w, fmt.Sprintf("expect POST method, but recieved: %s", r.Method), http.StatusMethodNotAllowed)
		}
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 2 {
		http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
		return
	}
}
