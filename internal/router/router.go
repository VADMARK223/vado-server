package router

import (
	"html/template"
	"vado_server/internal/appcontext"
	"vado_server/internal/constants/route"
	"vado_server/internal/handler/http"
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
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	r.SetHTMLTemplate(tmpl)
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil)
	// Статика и шаблоны
	r.Static("/static", "./web/static")

	// Настраиваем cookie-сессии
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("vado-session", store))
	r.Use(middleware.CheckJWT())

	// Публичные маршруты
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	r.GET(route.Index, http.ShowIndex)
	r.GET(route.Login, http.ShowLoginPage())
	r.POST(route.Login, http.PerformLogin(cxt))
	r.GET(route.Register, http.ShowRegisterPage())
	r.POST(route.Register, http.PerformRegister(userService))

	r.POST(route.Logout, http.Logout())

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(middleware.CheckAuth())
	{
		auth.GET(route.Tasks, http.ShowTasksPage(taskService))
		auth.POST(route.Tasks, http.AddTask(cxt))
		auth.DELETE("/tasks/:id", http.DeleteTask(cxt))
		auth.GET(route.Users, http.ShowUsers(cxt))
		auth.POST(route.Users, http.AddUser(userService))
		auth.GET(route.Roles, http.ShowRoles(roleService))
		auth.DELETE("/users/:id", http.DeleteUser(cxt))
	}

	// JSON API
	r.GET("/api/hello", http.GetHello())
	//r.GET("/api/tasks", handler.GetTasksJSON(taskService))

	return r
}
