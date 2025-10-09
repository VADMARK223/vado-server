package handler

import (
	"encoding/json"
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/repository"
)

func RegisterTaskRoutes(appCtx *appcontext.AppContext) {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		taskRepo := repository.NewTaskRepository(appCtx.DB)

		tasks, err := taskRepo.GetAll()
		if err != nil {
			appCtx.Log.Errorw("failed to get tasks", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tasks)
	})
}
