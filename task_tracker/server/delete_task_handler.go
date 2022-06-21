package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	taskstore "github.com/gsom95/go-ms/task_tracker/store"
)

func (ts *TaskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("deleteTaskHandler at path: %s with id %d", r.URL.Path, id)
	err := ts.store.DeleteTask(id)
	switch {
	case err == nil:
		return
	case errors.Is(err, taskstore.ErrTaskNotFound):
		http.Error(w, fmt.Sprintf("not found task %d", id), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
