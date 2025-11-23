package http

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
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

func SetupRouter(ctx *app.Context, tokenProvider *token.JWTProvider) *gin.Engine {
	// WS
	hub := ws.NewHub(ctx.Log)
	go hub.Run()
	// Сервисы
	taskSvc := task.NewService(gorm.NewTaskRepo(ctx.DB))

	userSvc := user.NewService(gorm.NewUserRepo(ctx), tokenProvider)
	localCache := app.NewLocalCache()

	// Хендлеры
	authH := handler.NewAuthHandler(userSvc, ctx.Cfg.JwtSecret, ctx.Cfg.RefreshTTLInt(), ctx.Log)

	gin.SetMode(ctx.Cfg.GinMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// Шаблоны
	r.SetHTMLTemplate(loadTemplates())
	_ = r.SetTrustedProxies(nil)
	// Статика и шаблоны
	r.Static("/static", "./web/static")
	// Favicon: отдаём напрямую, чтобы не было 404
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("web/static/favicon.ico")
	})

	// Middleware
	r.Use(middleware.SessionMiddleware())
	r.Use(middleware.CheckJWT(tokenProvider, ctx.Cfg.RefreshTTLInt()))
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

func loadTemplates() *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"dict": dict,
	})

	err := filepath.Walk("web/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				return fmt.Errorf("parse error in %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return tmpl
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call: odd number of args")
	}
	m := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m, nil
}
