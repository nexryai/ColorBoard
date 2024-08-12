package boot

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	apiController "github.com/nexryai/ColorBoard/internal/controller/api"
	authController "github.com/nexryai/ColorBoard/internal/controller/auth"
	"github.com/nexryai/ColorBoard/internal/logger"
	"github.com/nexryai/ColorBoard/internal/middleware"
	"github.com/nexryai/ColorBoard/internal/server"
	"github.com/nexryai/ColorBoard/internal/service/account"
	"github.com/nexryai/ColorBoard/internal/service/gallery"
	"github.com/nexryai/ColorBoard/internal/service/storage"
)

var (
	log = logger.GetLogger("Boot")
)

func Boot() {
	// Initialize session store
	log.Info("Initializing session store...")
	storeKey := securecookie.GenerateRandomKey(32)
	storeEncryptionKey := securecookie.GenerateRandomKey(32)

	sessionStore := sessions.NewCookieStore(storeKey, storeEncryptionKey)
	sessionStore.MaxAge(36000)
	sessionStore.Options.Path = "/api"
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = true
	sessionStore.Options.SameSite = http.SameSiteDefaultMode

	// Resolve dependencies
	log.Info("Initializing services and resolving dependencies...")
	storageService, err := storage.NewLocalStorageService()
	if err != nil {
		log.FatalWithDetail("Failed to initialize service: ", err)
		os.Exit(1)
	}

	userService := account.NewUserServices()
	galleryService := gallery.NewGalleryService(storageService)

	// Boot the server
	log.Info("Configuring routes...")
	router := gin.Default()
	router.Use(middleware.AuthMiddleware(sessionStore))
	server.ServeClient(router)

	// Config the OAuth router
	authController.ConfigOAuthRouter(router, userService, sessionStore)

	// Config API routers
	apiController.ConfigAccountAPIRouter(router, userService)
	apiController.ConfigGalleryAPIRouter(router, galleryService)
	apiController.ConfigSystemAPIRouter(router)

	// Start the server
	log.Info("Starting server...")
	err = router.Run(":8080")
	if err != nil {
		log.FatalWithDetail("Failed to start server", err)
	}
}
