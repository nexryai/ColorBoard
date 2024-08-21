package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nexryai/ColorBoard/internal/logger"
	"strings"
)

var (
	log = logger.GetLogger("AuthMiddleware")
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path

		if strings.HasPrefix(path, "/api/") {
			// ToDO: Cookieに載ってるJWTトークンが正しいか検証
		}

		ctx.Next()
	}
}
