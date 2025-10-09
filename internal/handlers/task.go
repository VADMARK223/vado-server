package handlers

import (
	"database/sql"
	"net/http"
	"time"
	"vado_server/internal/appcontext"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          int        `json:"id" example:"1" format:"int64"`             // Уникальный идентификатор
	Name        string     `json:"name" example:"Купить молоко"`              // Название задачи
	Description string     `json:"description" example:"Купить 2 литра"`      // Описание задачи
	Completed   bool       `json:"completed" example:"false"`                 // Флаг выполнения
	CreatedAt   *time.Time `json:"created_at" example:"2025-10-05T12:00:00Z"` // Время создания задачи
	UpdatedAt   *time.Time `json:"updated_at" example:"2025-10-05T12:00:00Z"` // Время обновления задачи
}

func GetTasks(ctx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := ctx.DB.Query(`SELECT id, name, description, completed, created_at, updated_at FROM tasks ORDER BY created_at`)
		if err != nil {
			ctx.Log.Errorw("Failed to query tasks", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
			return
		}
		defer func(rows *sql.Rows) {
			_ = rows.Close()
		}(rows)

		var tasks []Task
		for rows.Next() {
			var t Task
			if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
				ctx.Log.Errorw("Failed to scan task", "error", err)
				continue
			}
			tasks = append(tasks, t)
		}
		c.JSON(http.StatusOK, tasks)
	}
}
