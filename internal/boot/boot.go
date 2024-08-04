package boot

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	apiController "github.com/nexryai/ColorBoard/internal/controller/api"
	authController "github.com/nexryai/ColorBoard/internal/controller/auth"
	"github.com/nexryai/ColorBoard/internal/middleware"
	"github.com/nexryai/ColorBoard/internal/server"
	"github.com/nexryai/ColorBoard/internal/service/account"
	"net/http"
)

func Boot() {
	// Initialize session store
	storeKey := securecookie.GenerateRandomKey(32)
	storeEncryptionKey := securecookie.GenerateRandomKey(32)

	sessionStore := sessions.NewCookieStore(storeKey, storeEncryptionKey)
	sessionStore.MaxAge(36000)
	sessionStore.Options.Path = "/api"
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = true
	sessionStore.Options.SameSite = http.SameSiteDefaultMode

	// Boot the server
	router := gin.Default()
	router.Use(middleware.AuthMiddleware(sessionStore))
	server.ServeClient(router)

	// Resolve dependencies
	userService := account.NewUserServices()

	// Config the OAuth router
	authController.ConfigOAuthRouter(router, userService, sessionStore)

	// Config API routers
	apiController.ConfigAccountAPIRouter(router, userService)

	router.Run(":8080")
}
