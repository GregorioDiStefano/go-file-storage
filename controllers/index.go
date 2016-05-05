package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	return
}
