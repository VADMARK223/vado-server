package http

import (
	"html/template"
	"vado_server/internal/app/context"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/role"
	"vado_server/internal/domain/task"
	"vado_server/internal/domain/user"
	user2 "vado_server/internal/infra/persistence/user"
	"vado_server/internal/trasport/http/handler"

	//"vado_server/internal/middleware"
	"vado_server/internal/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(ctx *context.AppContext) *gin.Engine {
	taskService := task.NewService(task.NewRepo(ctx.DB))
	roleService := role.NewService(role.NewRepo(ctx.DB))
	userService := user.NewService(user2.NewGormRepo(ctx))

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
	r.POST(route.Login, handler.PerformLogin(ctx))
	r.GET(route.Register, handler.ShowRegisterPage())
	r.POST(route.Register, handler.PerformRegister(userService))

	r.POST(route.Logout, handler.Logout())

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(checkAuth())
	{
		auth.GET(route.Tasks, handler.ShowTasksPage(taskService))
		auth.POST(route.Tasks, handler.AddTask(ctx))
		auth.DELETE("/tasks/:id", handler.DeleteTask(ctx))
		auth.GET(route.Users, handler.ShowUsers(ctx))
		auth.POST(route.Users, handler.AddUser(userService))
		auth.GET(route.Roles, handler.ShowRoles(roleService))
		auth.DELETE("/users/:id", handler.DeleteUser(ctx))
	}

	// JSON API
	r.GET("/api/hello", handler.GetHello())
	//r.GET("/api/tasks", handler.GetTasksJSON(taskService))

	return r
}
