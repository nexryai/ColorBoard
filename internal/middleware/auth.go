package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/nexryai/ColorBoard/internal/logger"
	"strings"
)

var (
	log = logger.GetLogger("AuthMiddleware")
)

func AuthMiddleware(sessionStore *sessions.CookieStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		needAuth := false

		if strings.HasPrefix(path, "/api/") {
			if strings.HasPrefix(path, "/api/oauth/") && strings.HasSuffix(path, "/callback") {
				needAuth = false
			} else {
				needAuth = true
			}
		}

		if needAuth {
			session, err := sessionStore.Get(ctx.Request, "app_session")
			if err != nil || len(session.Values) == 0 {
				log.Info("Unauthorized access")
				ctx.AbortWithStatusJSON(401, gin.H{
					"error": "Unauthorized",
				})

				return
			}

			ctx.Set("userId", session.Values["user_id"].(string))
			ctx.Set("authUid", session.Values["auth_uid"].(string))
		}

		ctx.Next()
	}
}
