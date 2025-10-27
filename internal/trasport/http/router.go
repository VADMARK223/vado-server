package http

import (
	"html/template"
	"vado_server/internal/app/context"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/role"
	"vado_server/internal/domain/task"
	"vado_server/internal/domain/user"
	"vado_server/internal/trasport/http/handler"

	//"vado_server/internal/middleware"
	"vado_server/internal/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cxt *context.AppContext) *gin.Engine {
	// Сервисы
	taskRepo := task.NewTaskRepository(cxt.DB)
	roleRepo := role.NewRoleRepository(cxt.DB)
	userRepo := user.NewUserRepository(cxt.DB)

	taskService := task.NewTaskService(taskRepo)
	roleService := role.NewRoleService(roleRepo)
	userService := user.NewUserService(userRepo)

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
	r.Use(checkJWT())

	// Публичные маршруты
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	r.GET(route.Index, handler.ShowIndex)
	r.GET(route.Login, handler.ShowLoginPage())
	r.POST(route.Login, handler.PerformLogin(cxt))
	r.GET(route.Register, handler.ShowRegisterPage())
	r.POST(route.Register, handler.PerformRegister(userService))

	r.POST(route.Logout, handler.Logout())

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(checkAuth())
	{
		auth.GET(route.Tasks, handler.ShowTasksPage(taskService))
		auth.POST(route.Tasks, handler.AddTask(cxt))
		auth.DELETE("/tasks/:id", handler.DeleteTask(cxt))
		auth.GET(route.Users, handler.ShowUsers(cxt))
		auth.POST(route.Users, handler.AddUser(userService))
		auth.GET(route.Roles, handler.ShowRoles(roleService))
		auth.DELETE("/users/:id", handler.DeleteUser(cxt))
	}

	// JSON API
	r.GET("/api/hello", handler.GetHello())
	//r.GET("/api/tasks", handler.GetTasksJSON(taskService))

	return r
}
