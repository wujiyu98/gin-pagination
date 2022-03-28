package main

import (
	"github.com/gin-gonic/gin"
	pagination "github.com/wujiyu98/gin-pagination"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("example/gin-base/template/*")

	r.GET("/search", func(ctx *gin.Context) {

		p := pagination.GinInit(ctx, 10,5,1000)

		ctx.HTML(200, "index.html",gin.H{"bsPage": p.BsPage()})

	})


	r.Run()
	
}
