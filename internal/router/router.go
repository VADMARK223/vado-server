package router

import (
	"fmt"
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/handlers"
	"vado_server/internal/repository"
	"vado_server/internal/services"
	"vado_server/internal/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cxt *appcontext.AppContext) *gin.Engine {
	gin.SetMode(util.GetEnv("GIN_MODE"))
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil)
	// Статика и шаблоны
	r.Static("/static", "./internal/static")
	r.LoadHTMLGlob("internal/templates/*")

	// Настраиваем cookie-сессии
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("vado-session", store))

	// Сервисы
	taskRepo := repository.NewTaskRepository(cxt.DB)
	roleRepo := repository.NewRoleRepository(cxt.DB)
	taskService := services.NewTaskService(taskRepo)
	roleService := services.NewRoleService(roleRepo)

	// Публичные маршруты
	r.GET("/", handlers.ShowIndex)
	r.GET("/login", handlers.ShowLoginPage())
	r.POST("/login", handlers.PerformLogin(cxt))

	r.POST("/logout", handlers.Logout())

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(AuthRequired())
	{
		fmt.Println("Auth")
		auth.GET("/tasks", handlers.ShowTasksPage(taskService))
		auth.POST("/tasks", handlers.AddTask(cxt))
		auth.DELETE("/tasks/:id", handlers.DeleteTask(cxt))
		auth.GET("/users", handlers.ShowUsers(cxt))
		auth.POST("/users", handlers.AddUser(cxt))
		auth.GET("/roles", handlers.ShowRoles(roleService))
		auth.DELETE("/users/:id", handlers.DeleteUser(cxt))
	}

	// JSON API
	//r.GET("/api/tasks", handlers.GetTasksJSON(taskService))

	return r
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
