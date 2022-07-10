package shop

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
    // c.JSON(http.StatusOK, gin.H{
    //     "message": "Hello www.topgoer.com121",
    // })
	c.HTML(http.StatusOK, "app/shop/index.html", gin.H{
		"title": "app/shop/index.html",
	})
}