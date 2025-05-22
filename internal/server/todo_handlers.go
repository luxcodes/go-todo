package server

import (
	"encoding/json"
	"fmt"
	"go-todo/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := s.db.GetTodos()
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) getTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := s.db.GetTodo(id)
	if err != nil {
		http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

type newTodo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

func (s *Server) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo newTodo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if newTodo.Title == "" || newTodo.Description == "" || newTodo.Completed == nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	todo := models.Todo{
		Title:       newTodo.Title,
		Description: newTodo.Description,
		Completed:   *newTodo.Completed,
	}

	if err := s.db.CreateTodo(&todo); err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

type updateTodo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

func (s *Server) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the todo exists
	todo, err := s.db.GetTodo(id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	var updateTodo updateTodo
	if err := json.NewDecoder(r.Body).Decode(&updateTodo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateTodo.Title == "" || updateTodo.Description == "" || updateTodo.Completed == nil {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

	todo.Title = updateTodo.Title
	todo.Description = updateTodo.Description
	todo.Completed = *updateTodo.Completed

	if err := s.db.UpdateTodo(&todo); err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	todo, err := s.db.GetTodo(id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if err := s.db.DeleteTodo(todo.ID); err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDFromPath(r *http.Request) (int, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(parts) == 0 {
        return 0, fmt.Errorf("missing todo ID")
    }
    idStr := parts[len(parts)-1]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, fmt.Errorf("invalid todo ID")
    }
    return id, nil
}
