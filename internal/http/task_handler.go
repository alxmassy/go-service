package internalhttp

import (
	"encoding/json"
	"net/http"
	"task-service/internal/store"
)

type TaskHandler struct {
	store store.TaskStore
}

func NewTaskHandler(store store.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := h.store.List(ctx)
	if err != nil {
		http.Error(w, "Failed to list tasks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.store.Create(ctx, body.Title)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}