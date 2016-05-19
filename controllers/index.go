package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	return
}
