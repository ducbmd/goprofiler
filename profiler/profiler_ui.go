package profiler

import (
	"net/http"
	"sync"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

var once sync.Once

func InitUI() {
	once.Do(func() {
		initUI()
	})
}

func initUI() {
	router := gin.Default()

	//new template engine
	router.HTMLRender = ginview.Default()

	router.GET("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "index", gin.H{"title": "Index title!"})
	})

	router.GET("/page", func(ctx *gin.Context) {
		//render only file, must full name with extension
		ctx.HTML(http.StatusOK, "page.html", gin.H{"title": "Page file title!!"})
	})

	router.Run(":39001")
}
