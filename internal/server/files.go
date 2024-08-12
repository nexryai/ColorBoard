package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nexryai/ColorBoard/internal/logger"
)

var (
	log = logger.GetLogger("LocalStorage")
)

func serveFromLocalDir(ctx *gin.Context, dataDir string) {
	userId := ctx.MustGet("userId").(string)
    fileId := ctx.Param("fileId")

	if !strings.HasPrefix(fileId, fmt.Sprintf("local:%s:", userId)) {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
        return
    }
    
    // 保存されているファイルのパスを生成
    filePath := filepath.Join(dataDir, fileId)

    // ファイルが存在するか確認
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Warn("File not found: ", filePath)
        ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
        return
    }

    // ファイルを開く
    file, err := os.Open(filePath)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
        return
    }
    defer file.Close()

    // Content-Typeを設定してファイルをレスポンスとして送信
    ctx.Header("Content-Type", "application/octet-stream")
    ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
    io.Copy(ctx.Writer, file)
}

func ServceLocalStorageFiles(router *gin.Engine) {
	dataDir := os.Getenv("DATA_DIR")
    if dataDir == "" {
        panic("DATA_DIR environment variable is not set")
    }

	// ほんとは/api下に配置したくないけどセッションのCookieのPathの関係でこうしている
	router.GET("/api/files/:fileId", func(ctx *gin.Context) {
		serveFromLocalDir(ctx, dataDir)
	})
}