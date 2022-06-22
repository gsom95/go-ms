package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	taskstore "github.com/gsom95/go-ms/task_tracker/store"
)

func (ts *TaskServer) getTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("getTaskHandler at path: %s with id %d", r.URL.Path, id)
	task, err := ts.store.GetTask(id)
	if err == nil {
		renderJSON(w, task)
		return
	}

	if errors.Is(err, taskstore.ErrTaskNotFound) {
		http.Error(w, fmt.Sprintf("not found task %d", id), http.StatusBadRequest)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
