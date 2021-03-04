package goprofiler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initAPI() {
	router := gin.Default()

	router.GET("/api/realtime", func(ctx *gin.Context) {
		profiler := GetProfilerImpl()
		apis, _ := profiler.GetAllApis()

		jsonRet := gin.H{}
		for _, api := range apis {
			realtimeStats, _ := profiler.GetRealtimeStats(api)
			jsonRet[api] = realtimeStats
		}

		ctx.JSON(http.StatusOK, jsonRet)
	})

	router.GET("/api/history/second", func(ctx *gin.Context) {
		profiler := GetProfilerImpl()
		apis, _ := profiler.GetAllApis()

		jsonRet := gin.H{}
		for _, api := range apis {
			secondStats, _ := profiler.GetHistorySecondStats(api)
			jsonRet[api] = secondStats
		}

		ctx.JSON(http.StatusOK, jsonRet)
	})

	router.GET("/api/history/minute", func(ctx *gin.Context) {
		profiler := GetProfilerImpl()
		apis, _ := profiler.GetAllApis()

		jsonRet := gin.H{}
		for _, api := range apis {
			minuteStats, _ := profiler.GetHistoryMinuteStats(api)
			jsonRet[api] = minuteStats
		}

		ctx.JSON(http.StatusOK, jsonRet)
	})

	router.GET("/api/history/hour", func(ctx *gin.Context) {
		profiler := GetProfilerImpl()
		apis, _ := profiler.GetAllApis()

		jsonRet := gin.H{}
		for _, api := range apis {
			hourStats, _ := profiler.GetHistoryHourStats(api)
			jsonRet[api] = hourStats
		}

		ctx.JSON(http.StatusOK, jsonRet)
	})

	go router.Run(":39001")
}
