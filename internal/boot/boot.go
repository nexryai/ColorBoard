package boot

import (
	"github.com/gin-gonic/gin"
	authController "github.com/nexryai/ColorBoard/internal/controller/auth"
	"github.com/nexryai/ColorBoard/internal/middleware"
	"github.com/nexryai/ColorBoard/internal/server"
	"github.com/nexryai/ColorBoard/internal/service/account"
)

func Boot() {
	// Boot the server
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	server.ServeClient(router)

	// Resolve dependencies
	userService := account.NewUserServices()

	// Config the OAuth router
	authController.ConfigOAuthRouter(router, userService)

	router.Run(":8080")
}
