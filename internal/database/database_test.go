package database

import (
	"context"
	"go-todo/internal/models"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTodoCRUD(t *testing.T) {
    srv := New()

    // Create
    todo := &models.Todo{
        Title:       "Test Todo",
        Description: "Test Description",
        Completed:   false,
    }
    if err := srv.CreateTodo(todo); err != nil {
        t.Fatalf("CreateTodo failed: %v", err)
    }

    // List
    todos, err := srv.GetTodos()
    if err != nil {
        t.Fatalf("GetTodos failed: %v", err)
    }
    if len(todos) == 0 {
        t.Fatalf("expected at least one todo, got 0")
    }

    // Get (by ID)
    created := todos[len(todos)-1]
    got, err := srv.GetTodo(created.ID)
    if err != nil {
        t.Fatalf("GetTodo failed: %v", err)
    }
    if got.Title != todo.Title {
        t.Errorf("expected title %q, got %q", todo.Title, got.Title)
    }

    // Update
    created.Title = "Updated Title"
    created.Completed = true
    if err := srv.UpdateTodo(&created); err != nil {
        t.Fatalf("UpdateTodo failed: %v", err)
    }
    updated, err := srv.GetTodo(created.ID)
    if err != nil {
        t.Fatalf("GetTodo after update failed: %v", err)
    }
    if updated.Title != "Updated Title" || !updated.Completed {
        t.Errorf("update did not persist changes")
    }

    // Delete
    if err := srv.DeleteTodo(created.ID); err != nil {
        t.Fatalf("DeleteTodo failed: %v", err)
    }
    _, err = srv.GetTodo(created.ID)
    if err == nil {
        t.Fatalf("expected error after deleting todo, got nil")
    }
}

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	database = dbName
	password = dbPwd
	username = dbUser

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
