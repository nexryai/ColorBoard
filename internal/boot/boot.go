package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/nexryai/ColorBoard/internal/server"
)

func Boot() {
	// Boot the server
	router := gin.Default()
	server.ServeClient(router)

	router.Run(":9000")
}
