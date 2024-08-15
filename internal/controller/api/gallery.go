package api

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
    "github.com/nexryai/ColorBoard/db"
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

    // チェックサム
    expectedChecksum := ctx.PostForm("sha256")
    if expectedChecksum == "" {
        ctx.String(http.StatusBadRequest, "invalid hash value")
        return
    }
	
	// blurhashフィールドの取得
    blurhash := ctx.PostForm("blurhash")
    if blurhash == "" {
        ctx.String(http.StatusBadRequest, "invalid blurhash")
        return
    }

    width, err := strconv.Atoi(ctx.PostForm("width"))
    if err != nil {
		ctx.String(http.StatusBadRequest, "invalid width")
        return
	}

    height, err := strconv.Atoi(ctx.PostForm("height"))
    if err != nil {
		ctx.String(http.StatusBadRequest, "invalid height")
        return
	}

    // オリジナルファイルのreaderを取得
    losslessFile, _, err := ctx.Request.FormFile("lossless_data")
    if err != nil {
        ctx.String(http.StatusBadRequest, fmt.Sprintf("Error reading lossless data: %s", err.Error()))
        return
    }
    defer losslessFile.Close()

    // ファイルの整合性確認
    hasher := sha256.New()
    fileBuffer, err := io.ReadAll(losslessFile)
    if err != nil {
        ctx.String(http.StatusBadRequest, fmt.Sprintf("Error reading file data: %s", err.Error()))
        return
    }

    _, err = hasher.Write(fileBuffer)
    if err != nil {
        ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error hashing file data: %s", err.Error()))
        return
    }

    // ハッシュ確認
    sha256Hash := fmt.Sprintf("%x", hasher.Sum(nil))
    if sha256Hash != expectedChecksum {
        ctx.String(http.StatusBadRequest, fmt.Sprintf("checksum mismatch (calculated != received): %s != %s", sha256Hash, expectedChecksum))
        return
    }

    // thumbnail_dataファイルの処理
    thumbnailFile, _, err := ctx.Request.FormFile("thumbnail_data")
    if err != nil {
        ctx.String(http.StatusBadRequest, fmt.Sprintf("Error reading thumbnail data: %s", err.Error()))
        return
    }
    defer thumbnailFile.Close()

	// Add image to gallery
	res, err := galleryService.AddImage(
        bytes.NewReader(fileBuffer), 
        sha256Hash,
        thumbnailFile, 
        userId, 
        galleryId, 
        blurhash,
        width,
        height,
    )

	if err != nil {
        if _, e := db.IsErrUniqueConstraint(err); e {
            ctx.String(http.StatusBadRequest, "A same file already exists")
            return
        }

        ctx.String(http.StatusBadRequest, "bad request")
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
