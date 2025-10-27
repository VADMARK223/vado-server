package http

import (
	"html/template"
	"time"
	"vado_server/internal/app/context"
	"vado_server/internal/config/route"
	"vado_server/internal/config/token"
	"vado_server/internal/domain/role"
	"vado_server/internal/domain/task"
	"vado_server/internal/domain/user"
	user2 "vado_server/internal/infra/persistence/gorm"
	"vado_server/internal/trasport/http/handler"
	"vado_server/internal/trasport/http/middleware"

	//"vado_server/internal/middleware"
	"vado_server/internal/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(ctx *context.AppContext) *gin.Engine {
	// Сервисы
	taskSvc := task.NewService(task.NewRepo(ctx.DB))
	roleSvc := role.NewService(role.NewRepo(ctx.DB))
	userSvc := user.NewService(user2.NewUserRepo(ctx), token.AccessAliveMinutes*time.Minute)
	// Хендлеры
	authH := handler.NewAuthHandler(userSvc)
	userH := handler.NewUserHandler(userSvc)

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

	r.GET(route.Index, handler.ShowIndex)
	r.GET(route.Login, handler.ShowLoginPage())
	r.POST(route.Login, authH.Login)
	r.GET(route.Register, handler.ShowRegisterPage())
	r.POST(route.Register, handler.PerformRegister(userSvc))

	r.POST(route.Logout, handler.Logout())

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(middleware.CheckAuth())
	{
		auth.GET(route.Tasks, handler.ShowTasksPage(taskSvc))
		auth.POST(route.Tasks, handler.AddTask(ctx))
		auth.DELETE("/tasks/:id", handler.DeleteTask(ctx))
		auth.GET(route.Users, userH.ShowUsers)
		auth.POST(route.Users, handler.AddUser(userSvc))
		auth.GET(route.Roles, handler.ShowRoles(roleSvc))
		auth.DELETE("/users/:id", handler.DeleteUser(ctx))
	}

	return r
}
