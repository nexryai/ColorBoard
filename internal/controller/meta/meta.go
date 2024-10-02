package meta

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	FirebaseApiKey            string `json:"firebaseApiKey"`
	FirebaseAuthDomain        string `json:"firebaseAuthDomain"`
	FirebaseProjectId         string `json:"firebaseProjectId"`
	FirebaseStorageBucket     string `json:"firebaseStorageBucket"`
	FirebaseMessagingSenderId string `json:"firebaseMessagingSenderId"`
	FirebaseAppId             string `json:"firebaseAppId"`
}

func ConfigMetaRouter(router *gin.Engine) {
	meta := Meta{
		FirebaseApiKey:            os.Getenv("FIREBASE_API_KEY"),
		FirebaseAuthDomain:        os.Getenv("FIREBASE_AUTH_DOMAIN"),
		FirebaseProjectId:         os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseStorageBucket:     os.Getenv("FIREBASE_STORAGE_BUCKET"),
		FirebaseMessagingSenderId: os.Getenv("FIREBASE_MESSAGING_SENDER_ID"),
		FirebaseAppId:             os.Getenv("FIREBASE_APP_ID"),
	}

	router.GET("/meta", func(ctx *gin.Context) {
		// Firebaseの公開トークンを返す
		ctx.JSON(http.StatusOK, meta)
	})
}
