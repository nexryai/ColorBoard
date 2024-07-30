package boot

import (
	"github.com/gin-gonic/gin"
	authController "github.com/nexryai/ColorBoard/internal/controller/auth"
	"github.com/nexryai/ColorBoard/internal/server"
)

func Boot() {
	// Boot the server
	router := gin.Default()
	server.ServeClient(router)

	authController.ConfigOAuthRouter(router)

	router.Run(":8080")
}
