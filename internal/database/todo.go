package database

import (
	"go-todo/internal/models"
)

func (s *service) GetTodos() ([]models.Todo, error) {
	rows, err := s.db.Query("SELECT id, title, description, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *service) GetTodo(id int) (models.Todo, error) {
	var todo models.Todo
	err := s.db.QueryRow("SELECT id, title, description, completed FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed)
	if err != nil {	
		return models.Todo{}, err
	}
	return todo, nil
}

func (s *service) CreateTodo(todo *models.Todo) error {
	_, err := s.db.Exec("INSERT INTO todos (title, description, completed) VALUES ($1, $2, $3)", todo.Title, todo.Description, todo.Completed)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateTodo(todo *models.Todo) error {
	_, err := s.db.Exec("UPDATE todos SET title = $1, description = $2, completed = $3 WHERE id = $4", todo.Title, todo.Description, todo.Completed, todo.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTodo(id int) error {
	_, err := s.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
