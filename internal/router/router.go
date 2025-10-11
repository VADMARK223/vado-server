package router

import (
	"vado_server/internal/appcontext"
	"vado_server/internal/handlers"
	"vado_server/internal/repository"
	"vado_server/internal/services"
	"vado_server/internal/util"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cxt *appcontext.AppContext) *gin.Engine {
	gin.SetMode(util.GetEnv("GIN_MODE"))
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil)
	r.Static("/static", "./internal/static")
	r.LoadHTMLGlob("internal/templates/*")

	r.GET("/", handlers.ShowIndex)

	repo := repository.NewTaskRepository(cxt.DB)
	taskService := services.NewTaskService(repo)
	// HTML-страница
	r.GET("/tasks", handlers.ShowTasksPage(taskService))
	// JSON API
	r.GET("/api/tasks", handlers.GetTasksJSON(taskService))

	r.POST("/tasks", handlers.AddTask(cxt))
	r.DELETE("/tasks/:id", handlers.DeleteTask(cxt))

	r.GET("/users", handlers.ShowUsers(cxt))
	r.POST("/users", handlers.AddUser(cxt))

	r.GET("/roles", handlers.ShowRoles(cxt))

	r.DELETE("/users/:id", handlers.DeleteUser(cxt))

	return r
}
