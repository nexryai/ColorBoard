package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexryai/ColorBoard/internal/service"
)

type CreateGalleryReq struct {
	Name     string `form:"name" json:"name"`
	IsPublic bool   `form:"isPublic" json:"isPublic"`
}

func handleGalleryCreateAPI(ctx *gin.Context, galleryService service.IGalleryService) {
	userId := ctx.MustGet("userId").(string)

	var req CreateGalleryReq
	err := ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdId, err := galleryService.CreateGallery(
		&service.GalleryCreateParam{
			Name: req.Name,
			IsPublic: req.IsPublic,
			UserId: userId,
		},
	)

	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to create gallery",
		})
		return
	}

	ctx.JSON(200, createdId)
}

func handleGalleryGetAPI(ctx *gin.Context, galleryService service.IGalleryService) {
    userId := ctx.MustGet("userId").(string)
	id := ctx.Param("id")

    gallery, err := galleryService.GetGallery(userId, id)
    if err != nil {
        ctx.JSON(500, gin.H{
            "error": "Failed to get gallery",
        })
        return
    }

    ctx.JSON(200, gallery)
}

func handleListGalleriesAPI(ctx *gin.Context, galleryService service.IGalleryService) {
    userId := ctx.MustGet("userId").(string)

    galleries, err := galleryService.GetGalleriesByUserId(userId)
    if err != nil {
        ctx.JSON(500, gin.H{
            "error": "Failed to list galleries",
        })
        return
    }

    ctx.JSON(200, galleries)
}

func handleGalleryUploadAPI(ctx *gin.Context, galleryService service.IGalleryService) {
    userId := ctx.MustGet("userId").(string)
	galleryId := ctx.Param("id")
	if galleryId == "" {
		ctx.String(http.StatusBadRequest, "invalid request")
        return
	}
	
	// blurhashフィールドの取得
    blurhash := ctx.PostForm("blurhash")
    fmt.Println("Blurhash:", blurhash)

    // lossless_dataファイルの処理
    losslessFile, _, err := ctx.Request.FormFile("lossless_data")
    if err != nil {
        ctx.String(http.StatusBadRequest, fmt.Sprintf("Error reading lossless data: %s", err.Error()))
        return
    }
    defer losslessFile.Close()

    // thumbnail_dataファイルの処理
    thumbnailFile, _, err := ctx.Request.FormFile("thumbnail_data")
    if err != nil {
        ctx.String(http.StatusBadRequest, fmt.Sprintf("Error reading thumbnail data: %s", err.Error()))
        return
    }
    defer thumbnailFile.Close()

	// Add image to gallery
	res, err := galleryService.AddImage(losslessFile, thumbnailFile, userId, galleryId, blurhash)
	if err != nil {
        ctx.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err.Error()))
        return
    }

    ctx.String(http.StatusOK, res)
}

func ConfigGalleryAPIRouter(router *gin.Engine, galleryService service.IGalleryService) {
	router.POST("/api/gallery/create", func(ctx *gin.Context) {
		handleGalleryCreateAPI(ctx, galleryService)
	})

	router.GET("/api/gallery/:id", func(ctx *gin.Context) {
        handleGalleryGetAPI(ctx, galleryService)
    })

    router.GET("/api/gallery/list", func(ctx *gin.Context) {
        handleListGalleriesAPI(ctx, galleryService)
    })

	router.POST("/api/gallery/:id/upload", func(ctx *gin.Context) {
		handleGalleryUploadAPI(ctx, galleryService)
	})
}
