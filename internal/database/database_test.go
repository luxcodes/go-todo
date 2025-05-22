package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-todo/internal/models"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	created.Description = "Updated Description"
	created.Completed = true
	if err := srv.UpdateTodo(&created); err != nil {
		t.Fatalf("UpdateTodo failed: %v", err)
	}
	updated, err := srv.GetTodo(created.ID)
	if err != nil {
		t.Fatalf("GetTodo after update failed: %v", err)
	}
	if updated.Title != "Updated Title" || updated.Description != "Updated Description" || !updated.Completed {
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

func runMigrations(connStr string) error {
	log.Printf("Running migrations with connection string: %s", connStr)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Printf("runMigrations: failed to open DB: %v", err)
		return err
	}
	defer db.Close()

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		log.Printf("runMigrations: failed to create migrate driver: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres", driver)
	if err != nil {
		log.Printf("runMigrations: failed to create migrate instance: %v", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("runMigrations: migration failed: %v", err)
		return err
	}

	log.Printf("runMigrations: migrations applied successfully")
	return nil
}

func startPgContainer() (teardown func(context.Context, ...testcontainers.TerminateOption) error, err error) {
	database := "testdb"
	username := "testuser"
	password := "testpass"
	schema := "public"

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(database),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, err
	}
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, err
	}

	// Set environment variables for use in New()
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_DATABASE", database)
	os.Setenv("DB_USERNAME", username)
	os.Setenv("DB_PASSWORD", password)
	os.Setenv("DB_SCHEMA", schema)

	// Build connection string for migrations
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		username, password, host, port.Port(), database, schema,
	)
	if err := runMigrations(connStr); err != nil {
		return nil, err
	}

	return dbContainer.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := startPgContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	code := m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container")
	}
	os.Exit(code)
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
