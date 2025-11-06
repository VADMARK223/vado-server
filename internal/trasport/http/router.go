package http

import (
	"html/template"
	"strconv"
	"time"
	"vado_server/internal/app"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/role"
	"vado_server/internal/domain/task"
	"vado_server/internal/domain/user"
	"vado_server/internal/infra/persistence/gorm"
	"vado_server/internal/trasport/http/handler"
	"vado_server/internal/trasport/http/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter(ctx *app.Context) *gin.Engine {
	// Сервисы
	taskSvc := task.NewService(gorm.NewTaskRepo(ctx.DB))
	roleSvc := role.NewService(gorm.NewRoleRepo(ctx))
	tokenTTL, _ := strconv.Atoi(ctx.Cfg.TokenTTL)
	refreshTTL, _ := strconv.Atoi(ctx.Cfg.RefreshTTL)
	userSvc := user.NewService(gorm.NewUserRepo(ctx), time.Duration(tokenTTL)*time.Second, time.Duration(refreshTTL)*time.Second)
	// Хендлеры
	authH := handler.NewAuthHandler(userSvc, ctx.Cfg.JwtSecret, ctx.Cfg.TokenTTL)

	gin.SetMode(ctx.Cfg.GinMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// Шаблоны
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	r.SetHTMLTemplate(tmpl)
	_ = r.SetTrustedProxies(nil)
	// Статика и шаблоны
	r.Static("/static", "./web/static")
	// Favicon: отдаём напрямую, чтобы не было 404
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("web/static/favicon.ico")
	})

	// Настраиваем cookie-сессии
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("vado-session", store))
	r.Use(middleware.CheckJWT(ctx.Cfg.JwtSecret))
	r.Use(middleware.TemplateContext)

	// Публичные маршруты
	r.GET(route.Index, handler.ShowIndex(ctx.Cfg.JwtSecret))
	r.GET(route.Login, handler.ShowLogin)
	r.POST(route.Login, authH.Login)
	r.GET(route.Register, handler.ShowSignup)
	r.POST(route.Register, handler.PerformRegister(userSvc))
	r.POST(route.Logout, handler.Logout)

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(middleware.CheckAuthAndRedirect())
	{
		auth.GET(route.Tasks, handler.Tasks(taskSvc))
		auth.POST(route.Tasks, handler.AddTask(ctx))
		auth.DELETE("/tasks/:id", handler.DeleteTask(ctx))
		auth.GET(route.Users, handler.ShowUsers(userSvc))
		auth.POST(route.Users, handler.PostUser(userSvc))
		auth.DELETE("/users/:id", handler.DeleteUser(userSvc))
		auth.GET(route.Roles, handler.Roles(roleSvc))
		auth.GET("/grpc-test", handler.Grpc)
	}

	return r
}
