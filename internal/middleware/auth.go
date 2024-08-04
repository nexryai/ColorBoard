package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/nexryai/ColorBoard/internal/logger"
)

var (
	log = logger.GetLogger("AuthMiddleware")
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := gothic.Store.Get(ctx.Request, gothic.SessionName)
		if err != nil {
			log.ErrorWithDetail("Failed to get session", err)
			ctx.AbortWithStatusJSON(500, gin.H{
				"error": "Failed to get session",
			})
		}

		for k, v := range session.Values {
			log.Info(fmt.Sprintf("Key: %s, Value: %v", k, v))
		}

		ctx.Next()
	}
}
