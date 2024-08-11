package api

import (
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
}
