package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	if len(pathParts) == 1 {
		switch r.Method {
		case http.MethodPost:
			ts.createTaskHandler(w, r)
		default:
			http.Error(w, fmt.Sprintf("expect POST method, but recieved: %s", r.Method), http.StatusMethodNotAllowed)
		}
		return
	}

	if len(pathParts) != 2 {
		http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "expect numeric <id> in task handler: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ts.getTaskHandler(w, r, id)
	case http.MethodDelete:
		ts.deleteTaskHandler(w, r, id)
	default:
		http.Error(w, fmt.Sprintf("expect DELETE method, but recieved: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func renderJSON(w http.ResponseWriter, data any) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
