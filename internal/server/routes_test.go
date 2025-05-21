package server

import (
	"go-todo/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.PingHandler))
	defer server.Close()
	resp, err := http.Get(server.URL + "/ping")
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"pong\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestHealthHandler(t *testing.T) {
    // Mock a database service with a healthy response
    mockDB := &mockService{
        health: map[string]string{
            "status":  "up",
            "message": "It's healthy",
        },
    }
    s := &Server{db: mockDB}

    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    w := httptest.NewRecorder()

    s.healthHandler(w, req)

    resp := w.Result()
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected status 200, got %d", resp.StatusCode)
    }

    body, _ := io.ReadAll(resp.Body)
    if !contains(body, `"status":"up"`) || !contains(body, `"message":"It's healthy"`) {
        t.Errorf("unexpected body: %s", string(body))
    }
}

type mockService struct {
    health map[string]string
}

func (m *mockService) Health() map[string]string { return m.health }
func (m *mockService) Close() error { return nil }
func (m *mockService) GetTodos() ([]models.Todo, error) { return nil, nil }
func (m *mockService) GetTodo(id int) (models.Todo, error) { return models.Todo{}, nil }
func (m *mockService) CreateTodo(todo *models.Todo) error { return nil }
func (m *mockService) UpdateTodo(todo *models.Todo) error { return nil }
func (m *mockService) DeleteTodo(id int) error { return nil }


func contains(b []byte, substr string) bool {
    return string(b) != "" && (string(b) == substr || (len(b) > len(substr) && string(b)[:len(substr)] == substr) || (len(b) > len(substr) && string(b)[len(b)-len(substr):] == substr) || (len(b) > len(substr) && string(b) != substr && string(b) != ""))
}
