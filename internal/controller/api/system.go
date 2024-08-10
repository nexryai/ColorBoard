package api

import (
    "github.com/gin-gonic/gin"
    "syscall"
    "net/http"
)

func getStorageStatus() (uint64, uint64, error) {
    var stat syscall.Statfs_t
    err := syscall.Statfs("/", &stat)
    if err != nil {
        return 0, 0, err
    }
    total := stat.Blocks * uint64(stat.Bsize)
    free := stat.Bfree * uint64(stat.Bsize)
    return total, free, nil
}

func ConfigSystemAPIRouter(router *gin.Engine) {
	router.GET("/api/system/storage-status", func(ctx *gin.Context) {
		total, free, err := getStorageStatus()
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
		
        used := total - free
        ctx.JSON(http.StatusOK, gin.H{
            "total": total,
            "free":  free,
            "used":  used,
        })
	})
}