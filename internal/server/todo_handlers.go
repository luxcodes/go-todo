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

// @Summary Get all todos
// @Description Get all todos
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Router /todos [get]
func (s *Server) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET %s from %s", r.URL.Path, r.RemoteAddr)
	todos, err := s.db.GetTodos()
	if err != nil {
		log.Printf("getTodosHandler: failed to fetch todos: %v", err)
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("getTodosHandler: failed to write response: %v", err)
	}
}

// @Summary Get a todo by ID
// @Description Get a todo by ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Router /todo/{id} [get]
func (s *Server) getTodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET %s from %s", r.URL.Path, r.RemoteAddr)
	id, err := parseIDFromPath(r)
	if err != nil {
		log.Printf("getTodoHandler: invalid ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := s.db.GetTodo(id)
	if err != nil {
		log.Printf("getTodoHandler: failed to fetch todo with id %d: %v", id, err)
		http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("getTodoHandler: failed to write response: %v", err)
	}
}

type newTodo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

// @Summary Create todo
// @Description Create a new todo
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body newTodo true "Todo"
// @Success 201 {object} models.Todo
// @Router /todo/create [post]
func (s *Server) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST %s from %s", r.URL.Path, r.RemoteAddr)
	var newTodo newTodo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		log.Printf("createTodoHandler: invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if newTodo.Title == "" || newTodo.Description == "" || newTodo.Completed == nil {
		log.Printf("createTodoHandler: missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	todo := models.Todo{
		Title:       newTodo.Title,
		Description: newTodo.Description,
		Completed:   *newTodo.Completed,
	}

	if err := s.db.CreateTodo(&todo); err != nil {
		log.Printf("createTodoHandler: failed to create todo: %v", err)
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	log.Printf("createTodoHandler: created todo with id %d", todo.ID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("createTodoHandler: failed to write response: %v", err)
	}
}

type updateTodo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

// @Summary Update todo
// @Description Update an existing todo
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body updateTodo true "Todo"
// @Success 200 {object} models.Todo
// @Router /todo/update/{id} [put]
func (s *Server) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PUT %s from %s", r.URL.Path, r.RemoteAddr)
	id, err := parseIDFromPath(r)
	if err != nil {
		log.Printf("updateTodoHandler: invalid ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := s.db.GetTodo(id)
	if err != nil {
		log.Printf("updateTodoHandler: todo not found with id %d: %v", id, err)
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	var updateTodo updateTodo
	if err := json.NewDecoder(r.Body).Decode(&updateTodo); err != nil {
		log.Printf("updateTodoHandler: invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateTodo.Title == "" || updateTodo.Description == "" || updateTodo.Completed == nil {
		log.Printf("updateTodoHandler: missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	todo.Title = updateTodo.Title
	todo.Description = updateTodo.Description
	todo.Completed = *updateTodo.Completed

	if err := s.db.UpdateTodo(&todo); err != nil {
		log.Printf("updateTodoHandler: failed to update todo with id %d: %v", id, err)
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	log.Printf("updateTodoHandler: updated todo with id %d", id)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("updateTodoHandler: failed to write response: %v", err)
	}
}

// @Summary Delete todo
// @Description Delete a todo by ID
// @Tags todos
// @Param id path int true "Todo ID"
// @Success 204
// @Router /todo/delete/{id} [delete]
func (s *Server) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE %s from %s", r.URL.Path, r.RemoteAddr)
	id, err := parseIDFromPath(r)
	if err != nil {
		log.Printf("deleteTodoHandler: invalid ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := s.db.GetTodo(id)
	if err != nil {
		log.Printf("deleteTodoHandler: todo not found with id %d: %v", id, err)
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if err := s.db.DeleteTodo(todo.ID); err != nil {
		log.Printf("deleteTodoHandler: failed to delete todo with id %d: %v", id, err)
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	log.Printf("deleteTodoHandler: deleted todo with id %d", id)
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
