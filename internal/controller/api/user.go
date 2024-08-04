package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nexryai/ColorBoard/db"
	"github.com/nexryai/ColorBoard/internal/service"
)

func handleAccountInfoAPI(ctx *gin.Context, userService service.IUserService) {
	authUid := ctx.MustGet("authUid").(string)

	user, err := userService.GetUser(db.User.AuthUID.Equals(authUid))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	ctx.JSON(200, user)
	return
}

func ConfigAccountAPIRouter(router *gin.Engine, userService service.IUserService) {
	router.GET("/api/account/profile", func(ctx *gin.Context) {
		handleAccountInfoAPI(ctx, userService)
	})
}
