package auth

import (
	"github.com/gin-gonic/gin"
	logger "github.com/nexryai/ColorBoard/internal/logger"
	"github.com/nexryai/ColorBoard/internal/service"
)

var (
	log = logger.GetLogger("Auth")
)

func ConfigSupabaseAuthRouter(router *gin.Engine, userService service.IUserService) {
	router.GET("/auth/register-session", func(ctx *gin.Context) {
		// クライアントがSupabase Authから受け取ったセッショントークンをPOSTしてくるので、正しければCookieに載せて返す
	})
}
