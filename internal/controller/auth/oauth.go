package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/google"
	"github.com/nexryai/ColorBoard/db"
	logger "github.com/nexryai/ColorBoard/internal/logger"
	"github.com/nexryai/ColorBoard/internal/service"
	"net/http"
	"os"
)

var (
	appUrl = os.Getenv("APP_URL")
	log    = logger.GetLogger("OAuth")
)

func getCallbackURL(provider string) string {
	return fmt.Sprintf("%s/auth/%s/callback", appUrl, provider)
}

func contextWithProviderName(ctx *gin.Context, provider string) *http.Request {
	return ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

func ConfigOAuthRouter(router *gin.Engine, userService service.IUserService) {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), getCallbackURL("google")),
		azuread.New(os.Getenv("ENTRA_ID_KEY"), os.Getenv("ENTRA_ID_SECRET"), getCallbackURL("azuread"), nil),
	)

	router.GET("/auth/:provider", func(ctx *gin.Context) {
		provider := ctx.Param("provider")
		ctx.Request = contextWithProviderName(ctx, provider)

		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	})

	router.GET("/auth/:provider/callback", func(ctx *gin.Context) {
		provider := ctx.Param("provider")
		ctx.Request = contextWithProviderName(ctx, provider)

		user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
		if err != nil {
			log.ErrorWithDetail("failed to complete user auth", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		authUid := fmt.Sprintf("%s:%s", provider, user.UserID)

		dbUser, err := userService.GetUser(
			db.User.AuthUID.Equals(authUid),
		)

		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				// Create user
				_, err := userService.CreateUser(&service.UserCreateParam{
					Name:      user.Name,
					AuthUID:   authUid,
					AvatarUrl: user.AvatarURL,
				})

				if err != nil {
					log.ErrorWithDetail("failed to create user", err)
					ctx.Status(http.StatusInternalServerError)
					return
				} else {
					log.Info("user created: ", authUid)
				}
			} else {
				log.ErrorWithDetail("failed to get user", err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
		} else if dbUser == nil {
			log.Error("UNEXPECTED STATUS: user is nil")
			ctx.Status(http.StatusInternalServerError)
			return
		} else {
			// Update avatar url
			err := userService.UpdateAvatarUrl(
				db.User.AuthUID.Equals(authUid),
				user.AvatarURL,
			)

			if err != nil {
				log.ErrorWithDetail("failed to update avatar url", err)
				ctx.Status(http.StatusInternalServerError)
				return
			} else {
				log.Info("user updated: ", authUid)
			}
		}

		hostname := ctx.Request.Host
		ctx.SetCookie("auth_uid", authUid, 3600, "/", hostname, false, false)

		// redirect to main page
		ctx.Redirect(http.StatusFound, "/")
	})
}
