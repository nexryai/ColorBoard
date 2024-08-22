package firebase

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nexryai/fireproxy"
)

func ConfigFirebaseProxyRouter(router *gin.Engine) {
	projectName := os.Getenv("FIREBASE_PROJECT_ID")
	fireproxy.ConfigFirebaseAuthenticationProxy(router, projectName)
}
