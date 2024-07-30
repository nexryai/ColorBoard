package auth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/google"
	"log"
	"net/http"
	"os"
)

func contextWithProviderName(ctx *gin.Context, provider string) *http.Request {
	return ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

func ConfigOAuthRouter(router *gin.Engine) {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://127.0.0.1:8080/auth/google/callback"),
		azuread.New(os.Getenv("ENTRA_ID_KEY"), os.Getenv("ENTRA_ID_SECRET"), "http://127.0.0.1:8080/auth/azuread/callback", nil),
	)

	router.GET("/auth/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		c.Request = contextWithProviderName(c, provider)

		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	router.GET("/auth/:provider/callback", func(c *gin.Context) {
		provider := c.Param("provider")
		c.Request = contextWithProviderName(c, provider)

		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			fmt.Fprintln(c.Writer, err)
			return
		}

		log.Printf("%#v", user)
	})
}
