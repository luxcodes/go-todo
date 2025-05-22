package server

import (
	"go-todo/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockDBService struct {
	health map[string]string
	todos  map[int]models.Todo
	nextID int
}

func newMockDBService() *mockDBService {
	return &mockDBService{
		todos:  make(map[int]models.Todo),
		nextID: 1,
	}
}

func (m *mockDBService) Health() map[string]string { return m.health }
func (m *mockDBService) Close() error              { return nil }

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
	mockDB := &mockDBService{
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
	bodyStr := string(body)
	if !strings.Contains(bodyStr, `"status":"up"`) || !strings.Contains(bodyStr, `"message":"It's healthy"`) {
		t.Errorf("unexpected body: %s", bodyStr)
	}
}
