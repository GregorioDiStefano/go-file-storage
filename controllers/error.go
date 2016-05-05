package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func sendError(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, gin.H{
		"err": message,
	})
}
