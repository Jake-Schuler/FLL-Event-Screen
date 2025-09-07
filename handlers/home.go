package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func HomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title":      "FLL Event Screen - Dashboard",
			"event_name": os.Getenv("EVENT_NAME"),
		})
	}
}
