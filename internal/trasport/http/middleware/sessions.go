package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
	secret := "super-secret-key" // TODO: os.Getenv("SESSION_SECRET")
	if secret == "" {
		panic("Session cookie secret key is empty")
	}
	store := cookie.NewStore([]byte(secret))
	//store.Options(sessions.Options{
	//	HttpOnly: true,
	//	SameSite: http.SameSiteLaxMode,
	//	Path:     "/",
	//	MaxAge:   86400 * 7, // 7 days
	//})
	return sessions.Sessions("vado-session", store)
}
