package repository

import "database/sql"

type Task struct {
	ID          int
	Name        string
	Description string
	CreatedAt   string
	Completed   bool
}

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) GetAll() ([]Task, error) {
	rows, err := r.DB.Query("SELECT id, name, description, created_at, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
