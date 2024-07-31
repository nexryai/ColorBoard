package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nexryai/ColorBoard/client"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	clientDistFS     http.FileSystem = http.FS(client.DistFS)
	clientDistPrefix                 = "build"
)

func serveFromClientFS(ctx *gin.Context, pathPrefix string) {
	path := ctx.Param("filepath")
	rawFile, err := clientDistFS.Open(clientDistPrefix + "/" + pathPrefix + path)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		ctx.Status(http.StatusNotFound)
		return
	}

	p := strings.Split(path, ".")
	if len(p) < 1 {
		ctx.Status(http.StatusNotFound)
		return
	}

	fileExt := p[len(p)-1]
	switch fileExt {
	case "css":
		ctx.Header("Content-Type", "text/css")
	case "js":
		ctx.Header("Content-Type", "application/javascript")
	case "json":
		ctx.Header("Content-Type", "application/json")
	case "png":
		ctx.Header("Content-Type", "image/png")
	case "jpg":
		ctx.Header("Content-Type", "image/jpeg")
	case "webp":
		ctx.Header("Content-Type", "image/webp")
	case "svg":
		ctx.Header("Content-Type", "image/svg+xml")
	case "woff":
		ctx.Header("Content-Type", "font/woff")
	case "woff2":
		ctx.Header("Content-Type", "font/woff2")
	case "ttf":
		ctx.Header("Content-Type", "font/ttf")
	case "otf":
		ctx.Header("Content-Type", "font/otf")
	case "map":
		ctx.Header("Content-Type", "application/json")
	case "wasm":
		ctx.Header("Content-Type", "application/wasm")
	default:
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Status(http.StatusOK)
	_, _ = io.Copy(ctx.Writer, rawFile)
}

func ServeClient(router *gin.Engine) {
	router.GET("/_app/*filepath", func(ctx *gin.Context) {
		serveFromClientFS(ctx, "_app")
	})

	router.GET("/wasm/*filepath", func(ctx *gin.Context) {
		serveFromClientFS(ctx, "wasm")
	})

	// Serve the index.html
	router.NoRoute(func(ctx *gin.Context) {
		rawFile, err := clientDistFS.Open(clientDistPrefix + "/200.html")
		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		ctx.Status(http.StatusOK)
		_, _ = io.Copy(ctx.Writer, rawFile)
	})
}
