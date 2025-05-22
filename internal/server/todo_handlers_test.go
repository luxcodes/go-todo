package server

import (
	"encoding/json"
	"go-todo/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func (m *mockDBService) GetTodos() ([]models.Todo, error) {
	result := make([]models.Todo, 0, len(m.todos))
	for _, todo := range m.todos {
		result = append(result, todo)
	}
	return result, nil
}

func (m *mockDBService) GetTodo(id int) (models.Todo, error) {
	todo, ok := m.todos[id]
	if !ok {
		return models.Todo{}, http.ErrMissingFile
	}
	return todo, nil
}

func (m *mockDBService) CreateTodo(todo *models.Todo) error {
	todo.ID = m.nextID
	m.todos[todo.ID] = *todo
	m.nextID++
	return nil
}

func (m *mockDBService) UpdateTodo(todo *models.Todo) error {
	if _, ok := m.todos[todo.ID]; !ok {
		return http.ErrMissingFile
	}
	m.todos[todo.ID] = *todo
	return nil
}

func (m *mockDBService) DeleteTodo(id int) error {
	delete(m.todos, id)
	return nil
}

func TestGetTodosHandler(t *testing.T) {
	s := &Server{db: newMockDBService()}
	created := createTestTodo(s, models.Todo{Title: "Test", Description: "Test Desc", Completed: false})

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	s.getTodosHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var todos []models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if todos[0].Title != created.Title || todos[0].Description != created.Description || todos[0].Completed != created.Completed {
		t.Errorf("unexpected todos: %+v", todos)
	}
}

func TestGetTodoHandler(t *testing.T) {
	s := &Server{db: newMockDBService()}
	created := createTestTodo(s, models.Todo{Title: "Test", Description: "Test Desc", Completed: false})

	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	w := httptest.NewRecorder()

	s.getTodoHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var todo models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if todo.Title != created.Title || todo.Description != created.Description || todo.Completed != created.Completed {
		t.Errorf("unexpected todo: %+v", todo)
	}
}

func TestCreateTodoHandler(t *testing.T) {
	s := &Server{db: newMockDBService()}
	newTodo := models.Todo{Title: "New Todo", Description: "New Desc", Completed: false}
	createdTodo := createTestTodo(s, newTodo)

	if createdTodo.Title != newTodo.Title {
		t.Errorf("unexpected todo: %+v", createdTodo)
	}
	if createdTodo.Description != newTodo.Description {
		t.Errorf("unexpected todo: %+v", createdTodo)
	}
	if createdTodo.Completed != newTodo.Completed {
		t.Errorf("unexpected todo: %+v", createdTodo)
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	s := &Server{db: newMockDBService()}
	existingTodo := createTestTodo(s, models.Todo{Title: "Test", Description: "Test Desc", Completed: false})

	updated := models.Todo{ID: existingTodo.ID, Title: "Updated Todo", Description: "Updated Desc", Completed: true}
	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPut, "/todos/1", strings.NewReader(string(body)))
	w := httptest.NewRecorder()

	s.updateTodoHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var updatedTodo models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&updatedTodo); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if updatedTodo.Title != updated.Title {
		t.Errorf("unexpected todo: %+v", updatedTodo)
	}
	if updatedTodo.Description != updated.Description {
		t.Errorf("unexpected todo: %+v", updatedTodo)
	}
	if updatedTodo.Completed != updated.Completed {
		t.Errorf("unexpected todo: %+v", updatedTodo)
	}
	if updatedTodo.ID != existingTodo.ID {
		t.Errorf("expected todo ID %d, got %d", existingTodo.ID, updatedTodo.ID)
	}
}

func TestDeleteTodoHandler(t *testing.T) {
	s := &Server{db: newMockDBService()}
	created := createTestTodo(s, models.Todo{Title: "Test", Description: "Test Desc", Completed: false})

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	w := httptest.NewRecorder()

	s.deleteTodoHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", resp.StatusCode)
	}

	_, err := s.db.GetTodo(created.ID)
	if err == nil {
		t.Fatalf("expected error after deleting todo, got nil")
	}
}

func createTestTodo(s *Server, todo models.Todo) models.Todo {
	body, _ := json.Marshal(todo)
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(string(body)))
	w := httptest.NewRecorder()
	s.createTodoHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		panic("failed to create test todo")
	}
	var created models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		panic("failed to decode created todo")
	}
	return created
}
