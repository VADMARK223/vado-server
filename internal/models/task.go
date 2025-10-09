package models

import "time"

type Task struct {
	ID          int        `json:"id" example:"1" format:"int64"`             // Уникальный идентификатор
	Name        string     `json:"name" example:"Купить молоко"`              // Название задачи
	Description string     `json:"description" example:"Купить 2 литра"`      // Описание задачи
	Completed   bool       `json:"completed" example:"false"`                 // Флаг выполнения
	CreatedAt   *time.Time `json:"created_at" example:"2025-10-05T12:00:00Z"` // Время создания задачи
	UpdatedAt   *time.Time `json:"updated_at" example:"2025-10-05T12:00:00Z"` // Время обновления задачи
}
