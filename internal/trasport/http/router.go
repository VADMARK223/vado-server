package http

import (
	"html/template"
	"strconv"
	"time"
	"vado_server/internal/app"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/task"
	"vado_server/internal/domain/user"
	"vado_server/internal/infra/persistence/gorm"
	"vado_server/internal/infra/token"
	"vado_server/internal/trasport/http/handler"
	"vado_server/internal/trasport/http/middleware"
	"vado_server/internal/trasport/ws"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctx *app.Context) *gin.Engine {
	// WS
	hub := ws.NewHub(ctx.Log)
	go hub.Run()
	// Сервисы
	taskSvc := task.NewService(gorm.NewTaskRepo(ctx.DB))
	tokenTTL, _ := strconv.Atoi(ctx.Cfg.TokenTTL())
	refreshTTL, _ := strconv.Atoi(ctx.Cfg.RefreshTTL)
	tokenProvider := token.NewJWTProvider(ctx.Cfg.JwtSecret, time.Duration(tokenTTL)*time.Second, time.Duration(refreshTTL)*time.Second)
	userSvc := user.NewService(gorm.NewUserRepo(ctx), tokenProvider)
	localCache := app.NewLocalCache()
	// Хендлеры
	authH := handler.NewAuthHandler(userSvc, ctx.Cfg.JwtSecret, ctx.Cfg.TokenTTL(), ctx.Cfg.RefreshTTL, ctx.Log)

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

	// Middleware
	r.Use(middleware.SessionMiddleware())
	r.Use(middleware.CheckJWT(tokenProvider, ctx.Cfg.RefreshTTL))
	r.Use(middleware.LoadUserContext(userSvc, localCache))
	r.Use(middleware.NoCache)
	r.Use(middleware.TemplateContext)

	// Публичные маршруты
	r.GET(route.Index, handler.ShowIndex(tokenProvider))
	r.GET(route.Book, handler.ShowBook)
	r.GET(route.Login, handler.ShowLogin)
	r.POST(route.Login, authH.Login)
	r.GET(route.Register, handler.ShowSignup)
	r.POST(route.Register, handler.PerformRegister(userSvc))
	r.POST(route.Logout, handler.Logout)
	r.GET("/ws", handler.ServeSW(hub, ctx.Log, tokenProvider))
	r.GET("/chat", handler.ShowChat())

	// Защищенные маршруты
	auth := r.Group("/")
	auth.Use(middleware.CheckAuthAndRedirect())
	{
		auth.GET(route.Tasks, handler.Tasks(taskSvc))
		auth.POST(route.Tasks, handler.AddTask(ctx))
		auth.DELETE("/tasks/:id", handler.DeleteTask(ctx))
		auth.PUT("/tasks/:id", handler.UpdateTask(ctx, taskSvc))
		auth.GET(route.Users, handler.ShowUsers(userSvc))
		auth.DELETE("/users/:id", handler.DeleteUser(userSvc))
		auth.GET("/grpc-test", handler.Grpc)
	}

	return r
}
