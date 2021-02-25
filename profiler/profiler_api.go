package profiler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initAPI() {
	router := gin.Default()

	router.GET("/api/realtime", func(ctx *gin.Context) {
		profiler := GetProfilerImpl()
		apis, _ := profiler.GetAllApis()
		ctx.JSON(http.StatusOK, "{}")
	})

	router.GET("/api/history/second", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "{}")
	})

	router.GET("/api/history/minute", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "{}")
	})

	router.GET("/api/history/hour", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "{}")
	})

	router.Run(":39001")
}
