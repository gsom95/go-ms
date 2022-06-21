package server

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"time"
)

type CreateTaskRequest struct {
	Text string    `json:"text"`
	Due  time.Time `json:"due"`
}

type CreateTaskResponse struct {
	ID int `json:"id"`
}

func (ts *TaskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("createTaskHandler at path: %s", r.URL.Path)

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var rt CreateTaskRequest
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := ts.store.CreateTask(rt.Text, rt.Due)
	if err != nil {
		http.Error(w, fmt.Sprintf("store CreateTask error: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(CreateTaskResponse{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}
