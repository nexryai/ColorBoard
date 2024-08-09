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

// AuthMiddleware is a middleware function that handles authentication for API routes.
// It checks if the request path starts with "/api/" and if authentication is required.
// If authentication is required, it checks if the user is authorized by checking the session values.
// If the user is not authorized, it returns a 401 Unauthorized response.
// If the user is authorized, it sets the "userId" and "authUid" values in the context and proceeds to the next middleware or handler.
// The sessionStore parameter is a pointer to a sessions.CookieStore used to retrieve session values.
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
