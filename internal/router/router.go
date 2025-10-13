package router

import (
	"vado_server/internal/appcontext"
	"vado_server/internal/constants/route"
	"vado_server/internal/handlers"
	"vado_server/internal/middleware"
	"vado_server/internal/repository"
	"vado_server/internal/services"
	"vado_server/internal/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cxt *appcontext.AppContext) *gin.Engine {
	// Сервисы
	taskRepo := repository.NewTaskRepository(cxt.DB)
	roleRepo := repository.NewRoleRepository(cxt.DB)
	userRepo := repository.NewUserRepository(cxt.DB)

	taskService := services.NewTaskService(taskRepo)
	roleService := services.NewRoleService(roleRepo)
	userService := services.NewUserService(userRepo)

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
	r.Use(middleware.CheckJWT())

	// Публичные маршруты
	r.GET(route.Index, handlers.ShowIndex)
	r.GET(route.Login, handlers.ShowLoginPage())
	r.POST(route.Login, handlers.PerformLogin(cxt))
	r.GET(route.Register, handlers.ShowRegisterPage())
	r.POST(route.Register, handlers.PerformRegister(userService))

	r.POST(route.Logout, handlers.Logout())

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(middleware.CheckAuth())
	{
		auth.GET(route.Tasks, handlers.ShowTasksPage(taskService))
		auth.POST(route.Tasks, handlers.AddTask(cxt))
		auth.DELETE("/tasks/:id", handlers.DeleteTask(cxt))
		auth.GET(route.Users, handlers.ShowUsers(cxt))
		auth.POST(route.Users, handlers.AddUser(userService))
		auth.GET(route.Roles, handlers.ShowRoles(roleService))
		auth.DELETE("/users/:id", handlers.DeleteUser(cxt))
	}

	// JSON API
	//r.GET("/api/tasks", handlers.GetTasksJSON(taskService))

	return r
}
